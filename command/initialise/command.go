package initialise

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/dihedron/archetype/command/base"
	"github.com/dihedron/archetype/logging"
	"github.com/dihedron/archetype/pointer"
	"github.com/dihedron/archetype/repository"
	"github.com/dihedron/archetype/settings"
	"gopkg.in/yaml.v3"
)

// Initialise is the command to initialise a new repository from an archetype.
type Initialise struct {
	base.Command
	// Settings is the path to the settings file to use for saturating the archetype variables.
	Settings settings.Settings `short:"s" long:"settings" description:"The settings used to transform the archetype into an actual repository" required:"true"`
	// Directory is the path to the directory to use for the archetype files.
	Directory string `short:"d" long:"directory" description:"The directory where the output files are stored" required:"true" default:".archetype/output"`
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
func (cmd *Initialise) Execute(args []string) error {
	slog.Info("executing Initialise command")

	var options []repository.Option

	slog.Debug("command configuration", "settings", logging.ToJSON(cmd.Settings), "directory", cmd.Directory)

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

	// 6. validate the user-provided settings against the remote archetype metadata
	file, err := commit.File(".archetype/metadata.yml")
	if err != nil {
		slog.Error("failed to get archetype metadata file from repository", "error", err)
		return fmt.Errorf("failed to get archetype metadata file from repository: %w", err)
	}
	var contents string
	if contents, err = file.Contents(); err != nil {
		slog.Error("failed to get contents of archetype metadata file", "error", err)
		return err
	}
	var metadata settings.Metadata
	if err := yaml.Unmarshal([]byte(contents), &metadata); err != nil {
		slog.Error("failed to unmarshal archetype metadata file", "error", err)
		return err
	}
	slog.Info("loaded archetype metadata", "version", metadata.Version, "parameters", logging.ToJSON(metadata.Parameters))
	if cmd.Settings.Version != metadata.Version {
		slog.Error("unsupported settings version", "version", cmd.Settings.Version, "expected", metadata.Version)
		return fmt.Errorf("unsupported settings version: %d (expected %d)", cmd.Settings.Version, metadata.Version)
	}
	// 7. load and validate the parameters from the settings
	context := map[string]any{}
	fmt.Printf("---- %s ----\n", yellow("PARAMETERS"))
	for key, value := range cmd.Settings.Parameters {
		meta, ok := metadata.Parameters[key]
		if !ok {
			slog.Error("unsupported parameter in settings", "parameter", key)
			return fmt.Errorf("unsupported parameter in settings: %s", key)
		}
		expected := strings.ReplaceAll(meta.Type, "interface {}", "any")
		got := strings.ReplaceAll(fmt.Sprintf("%T", value), "interface {}", "any")
		if expected != got {
			slog.Error("parameter type mismatch", "parameter", key, "expected", meta.Type, "got", fmt.Sprintf("%T", value))
			return fmt.Errorf("parameter type mismatch for '%s': expected %s, got %s", key, meta.Type, fmt.Sprintf("%T", value))
		}
		if value == nil && meta.Default != nil {
			context[key] = meta.Default
		} else {
			context[key] = value
		}
		fmt.Printf("'%s' => '%s' (type: %s)\n",
			green(key),
			fmt.Sprintf("%v", green(context[key])),
			blue(fmt.Sprintf("%T", value)),
		)
		//fmt.Sprintf("%T", blue(parameter.Value)))
	}
	fmt.Printf("---- %s ----\n", yellow("PARAMETERS"))

	// 6. loop over the files and perform some processing
	repo.ForEachFile(commit, FileVisitor(cmd.Directory, context))

	// 7. launch the script for post processing (TODO)

	return nil
}
