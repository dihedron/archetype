package initialise

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/dihedron/archetype/command/base"
	"github.com/dihedron/archetype/command/commons"
	"github.com/dihedron/archetype/logging"
	"github.com/dihedron/archetype/pointer"
	"github.com/dihedron/archetype/repository"
	"github.com/dihedron/archetype/settings"
)

type Initialise struct {
	base.Command
	// Settings is the path to the settings file to use for saturating the archetype variables.
	Settings settings.Settings `short:"s" long:"settings" description:"The settings used to define the archetype and the specific parameters" required:"true"`
	// Directory is the path to the directory to use for the archetype files.
	Directory string `short:"d" long:"directory" description:"The directory where the output files are stored" required:"true" default:"tmp/archetype-output"`
}

const (
	CurrentVersion              = 1
	DefaultDirectoryPermissions = 0755
	DefaultFilePermissions      = 0644
)

// Execute runs the Initialise command.
func (cmd *Initialise) Execute(args []string) error {
	slog.Info("executing Initialise command")

	var options []repository.Option

	slog.Debug("command configuration", "settings", logging.ToJSON(cmd.Settings), "directory", cmd.Directory)

	if cmd.Settings.Version != CurrentVersion {
		slog.Error("unsupported settings version", "version", cmd.Settings.Version, "expected", CurrentVersion)
		return fmt.Errorf("unsupported settings version: %d (expected %d)", cmd.Settings.Version, CurrentVersion)
	}

	if cmd.Settings.Repository.URL == "" {
		slog.Error("repository URL not specified in settings")
		return fmt.Errorf("repository URL not specified in settings")
	}

	// 1. create the output directory if it does not exist; check if it is empty
	if err := os.MkdirAll(cmd.Directory, DefaultDirectoryPermissions); err != nil {
		slog.Error("failed to create output directory", "directory", cmd.Directory, "error", err)
		return err
	}
	if files, err := os.ReadDir(cmd.Directory); err != nil {
		slog.Error("failed to read output directory", "directory", cmd.Directory, "error", err)
		return err
	} else if len(files) > 0 {
		slog.Warn("output directory is not empty", "directory", cmd.Directory)
	}

	// 2. extract authentication options
	if cmd.Settings.Repository.Auth != nil {
		// extract and validate auth settings
		auth, err := commons.AuthenticationOpts(cmd.Settings.Repository)
		if err != nil {
			slog.Error("error validating authentication options", "error", err)
			return err
		}
		if auth != nil {
			options = append(options, auth)
		}
	}
	//options = append(options, repository.WithProxyFromEnv())

	// 3. create an in-memory clone of the remote archetypal repository
	repo := repository.New(
		cmd.Settings.Repository.URL,
		options...,
	)
	if err := repo.Clone(); err != nil {
		slog.Error("failed to clone remote repository", "url", cmd.Settings.Repository.URL, "error", err)
		return err
	}

	// 4. checkout the specified tag (or 'latest' if none specified)
	if cmd.Settings.Repository.Tag == nil {
		slog.Info("no tag specified, using 'latest' as default")
		cmd.Settings.Repository.Tag = pointer.To("latest")
	}
	commit, err := repo.Commit(*cmd.Settings.Repository.Tag)
	if err != nil {
		slog.Error("failed to get commit for input tag", "tag", *cmd.Settings.Repository.Tag, "error", err)
		return err
	}

	// 5. load the parameters
	context := map[string]any{}
	fmt.Printf("---- PARAMETERS ----\n")
	for _, parameter := range cmd.Settings.Parameters {
		if parameter.Value.Value == nil && parameter.Default != nil {
			context[parameter.Name] = parameter.Default
		} else {
			context[parameter.Name] = parameter.Value.Value
		}
		fmt.Printf("'%s' => '%v' (type: %T)\n", parameter.Name, context[parameter.Name], context[parameter.Name])
	}
	fmt.Printf("---- end of PARAMETERS ----\n")

	// 6. loop over the files and perform some processing
	repo.ForEachFile(commit, FileVisitor(cmd.Directory, context))

	// 7. launch the script for post processing (TODO)

	return nil
}
