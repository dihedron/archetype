package show

import (
	"bufio"
	"fmt"
	"log/slog"

	"github.com/dihedron/archetype/command/base"
	"github.com/dihedron/archetype/command/initialise"
	"github.com/dihedron/archetype/logging"
	"github.com/dihedron/archetype/pointer"
	"github.com/go-git/go-git/v6/plumbing/object"
)

type Show struct {
	base.Command
}

func (cmd *Show) Execute(args []string) error {

	p := initialise.Settings{
		Version: 1,
		Repository: initialise.Repository{
			URL: "https://github.com/go-git/go-git.git",
			Tag: pointer.To("latest"),
			Auth: &initialise.Auth{
				Token:    pointer.To("my-token"),
				Username: pointer.To("my-username"),
				Password: pointer.To("my-password"),
				SSHKey:   pointer.To("my-ssh-key"),
			},
		},
		Parameters: []initialise.Parameter{
			{
				Name:        "name1",
				Type:        "bool",
				Description: "Description 1",
				Default:     "default 1",
				Value:       "value 1",
			},
		},
	}

	fmt.Printf("%s\n", logging.ToYAML(p))

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
