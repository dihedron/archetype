package prepare

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/dihedron/archetype/command/base"
	"github.com/dihedron/archetype/pointer"
	"github.com/dihedron/archetype/repository"
)

// Prepare is the command to prepare a file by escaping all Golang-template directives.
type Prepare struct {
	base.Command
	// Directory is the path to the directory to use to store the "prepared" files.
	Directory string `short:"d" long:"directory" description:"The directory where the output files are stored" required:"true" default:".archetype/prepared"`
}

const (
	CurrentVersion              = 1
	DefaultDirectoryPermissions = 0755
	DefaultFilePermissions      = 0644
)

// Execute is the main entry point for the initialise command.
// It clones the archetype repository, validates the provided settings against the archetype's metadata,
// and then processes the files in the repository, treating them as templates and executing them with the provided settings.
// The resulting files are written to the output directory.
func (cmd *Prepare) Execute(args []string) error {
	slog.Info("executing Prepare command")

	if len(cmd.Exclude) > 0 && len(cmd.Include) > 0 {
		slog.Warn("both exclude and include patterns specified; include patterns will take precedence")
		fmt.Fprintf(os.Stderr, "Both exclude and include patterns specified; include patterns will take precedence\n")
	}
	var options []repository.Option

	// 1. check that the repository URL is specified
	if cmd.URL == "" {
		slog.Error("repository URL not specified in settings")
		return fmt.Errorf("repository URL not specified in settings")
	}

	// 2. create the output directory if it does not exist; check if it is empty
	if err := os.MkdirAll(cmd.Directory, DefaultDirectoryPermissions); err != nil {
		slog.Error("failed to create output directory", "directory", cmd.Directory, "error", err)
		return fmt.Errorf("failed to create output directory '%s': %w", cmd.Directory, err)
	}
	if files, err := os.ReadDir(cmd.Directory); err != nil {
		slog.Error("failed to read output directory", "directory", cmd.Directory, "error", err)
		return fmt.Errorf("failed to read output directory '%s': %w", cmd.Directory, err)
	} else if len(files) > 0 {
		slog.Warn("output directory is not empty", "directory", cmd.Directory)
	}

	// 3. extract authentication options
	if cmd.HasAuthOptions() {
		// extract and validate auth settings
		if auth, err := cmd.AuthenticationOpts(); err != nil {
			slog.Error("error validating authentication options", "error", err)
			return fmt.Errorf("error validating authentication options: %w", err)
		} else if auth != nil {
			options = append(options, auth)
		}
	}
	//options = append(options, repository.WithProxyFromEnv())

	// 4. create an in-memory clone of the remote archetypal repository
	repo, err := repository.New(cmd.URL, options...)
	if err != nil {
		slog.Error("failed to clone remote repository", "url", cmd.URL, "error", err)
		return fmt.Errorf("failed to clone remote repository '%s': %w", cmd.URL, err)
	}

	// 5. checkout the specified tag (or 'latest' if none specified)
	if cmd.Tag == nil {
		slog.Info("no tag specified, using 'latest' as default")
		cmd.Tag = pointer.To("latest")
	}
	commit, err := repo.Commit(*cmd.Tag)
	if err != nil {
		slog.Error("failed to get commit for input tag", "tag", *cmd.Tag, "error", err)
		return fmt.Errorf("failed to get commit for input tag '%s': %w", *cmd.Tag, err)
	}

	// // 6. loop over the files and perform some processing
	// repo.ForEachFile(commit, func(file *object.File) error {
	// 	fmt.Printf("processing file %s (mode: %v, size: %d, hash: %s)...\n", file.Name, file.Mode, file.Size, file.Hash.String())
	// 	return nil
	// })

	// 6. loop over the files and perform some processing
	repo.ForEachFile(commit, FileVisitor(cmd.Directory, cmd.Exclude, cmd.Include))

	// 7. launch the script for post processing (TODO)

	return nil
}
