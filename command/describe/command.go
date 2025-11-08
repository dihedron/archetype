package describe

import (
	"fmt"
	"log/slog"

	"github.com/dihedron/archetype/command/base"
	"github.com/dihedron/archetype/logging"
	"github.com/dihedron/archetype/pointer"
	"github.com/dihedron/archetype/repository"
	"github.com/dihedron/archetype/settings"
	"gopkg.in/yaml.v3"
)

// Describe is the command to describe the settings and parameters of an archetype.
type Describe struct {
	base.Command
}

// Execute is the main entry point for the describe command.
func (cmd *Describe) Execute(args []string) error {

	slog.Info("executing Describe command")

	var options []repository.Option

	// 1. check that the repository URL is specified
	if cmd.URL == "" {
		slog.Error("repository URL not specified in settings")
		return fmt.Errorf("repository URL not specified in settings")
	}

	// 2. extract authentication options
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

	// 3. create an in-memory clone of the remote archetypal repository
	repo, err := repository.New(cmd.URL, options...)
	if err != nil {
		slog.Error("failed to clone remote repository", "url", cmd.URL, "error", err)
		return fmt.Errorf("failed to clone remote repository '%s': %w", cmd.URL, err)
	}

	// 4. checkout the specified tag (or 'latest' if none specified)
	if cmd.Tag == nil {
		slog.Info("no tag specified, using 'latest' as default")
		cmd.Tag = pointer.To("latest")
	}
	commit, err := repo.Commit(*cmd.Tag)
	if err != nil {
		slog.Error("failed to get commit for input tag", "tag", *cmd.Tag, "error", err)
		return fmt.Errorf("failed to get commit for input tag '%s': %w", *cmd.Tag, err)
	}

	// 5. validate the user-provided settings against the remote archetype metadata
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
	settings := &settings.Settings{
		Version:    metadata.Version,
		Parameters: map[string]any{},
	}
	for key, value := range metadata.Parameters {
		settings.Parameters[key] = value
	}
	for key, value := range metadata.Parameters {
		settings.Parameters[key] = value.Default
	}

	fmt.Printf("%s", logging.ToYAML(settings))

	return nil

}

// // VisitFile is a callback function that is invoked for each file in the repository.
// func VisitFile(file *object.File) error {
// 	fmt.Printf("%v  %9d  %s    %s\n", file.Mode, file.Size, file.Hash.String(), file.Name)

// 	if file.Name == ".archetype/metadata.yml" {
// 		reader, err := file.Blob.Reader()
// 		if err != nil {
// 			return err
// 		}
// 		defer reader.Close()

// 		scanner := bufio.NewScanner(reader)
// 		// optionally, resize scanner's capacity for lines over 64K, see next example
// 		for scanner.Scan() {
// 			fmt.Println(scanner.Text())
// 		}

// 		if err := scanner.Err(); err != nil {
// 			slog.Error("error scanning file", "error", err)
// 			return err
// 		}
// 	}

// 	return nil
// }
