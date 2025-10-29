package show

import (
	"bufio"
	"fmt"
	"log/slog"

	"github.com/dihedron/archetype/command/base"
	"github.com/dihedron/archetype/logging"
	"github.com/dihedron/archetype/pointer"
	"github.com/dihedron/archetype/settings"
	"github.com/go-git/go-git/v6/plumbing/object"
)

type Show struct {
	base.Command
	URL  string         `json:"url,omitempty" yaml:"url,omitempty" short:"r" long:"repository" description:"The Git repository containing the template" required:"true"`
	Tag  *string        `json:"tag,omitempty" yaml:"tag,omitempty" short:"t" long:"tag" description:"The tag or commit to clone" optional:"true" default:"latest"`
	Auth *settings.Auth `json:"auth,omitempty" yaml:"auth,omitempty" group:"Authentication Options" description:"Authentication options"`
}

func (cmd *Show) Execute(args []string) error {

	p := settings.Settings{
		Version: 1,
		Repository: settings.Repository{
			URL: "https://github.com/go-git/go-git.git",
			Tag: pointer.To("latest"),
			Auth: &settings.Auth{
				Token:    pointer.To("my-token"),
				Username: pointer.To("my-username"),
				Password: pointer.To("my-password"),
				SSHKey:   pointer.To("my-ssh-key"),
			},
		},
		Parameters: []settings.Parameter{
			{
				Name:        "name1",
				Description: "Description 1",
				Default:     "default 1",
				Value: settings.Value{
					Type:  "bool",
					Value: "value 1",
				},
			},
		},
	}

	fmt.Printf("%s", logging.ToYAML(p))

	return nil

	// fmt.Printf("Executing Show command")

	// auth, err := cmd.AuthenticationOpts()
	// if err != nil {
	// 	slog.Error("error validating authentication options", "error", err)
	// 	return err
	// }

	// repo := repository.New(
	// 	cmd.Repository,
	// 	auth,
	// )
	// repo.Clone()

	// var reference *plumbing.Reference
	// if cmd.Tag == "HEAD" {
	// 	var err error
	// 	// ... retrieving the branch being pointed by HEAD
	// 	reference, err = repo.Head()
	// 	if err != nil {
	// 		slog.Error("failed to get HEAD", "error", err)
	// 		os.Exit(1)
	// 	}
	// 	fmt.Println("HEAD points to:", reference.Name())
	// } else {
	// 	var err error
	// 	reference, err = repo.Tag(cmd.Tag)
	// 	if err != nil {
	// 		slog.Error("failed to get tag", "error", err)
	// 		os.Exit(1)
	// 	}
	// }
	// commit, err := repo.CommitFromReference(reference)
	// if err != nil {
	// 	slog.Error("failed to get commit from reference", "reference", reference.Name().String(), "error", err)
	// 	os.Exit(1)
	// }

	// repo.ForEachFile(commit, VisitFile)

	// return nil
}

func VisitFile(file *object.File) error {
	fmt.Printf("%v  %9d  %s    %s\n", file.Mode, file.Size, file.Hash.String(), file.Name)

	if file.Name == ".archetype/metadata.yml" {
		reader, err := file.Blob.Reader()
		if err != nil {
			return err
		}
		defer reader.Close()

		scanner := bufio.NewScanner(reader)
		// optionally, resize scanner's capacity for lines over 64K, see next example
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			slog.Error("error scanning file", "error", err)
			return err
		}
	}

	return nil
}
