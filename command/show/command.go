package show

import (
	"bufio"
	"fmt"
	"log/slog"

	"github.com/dihedron/archetype/command/base"
	"github.com/dihedron/archetype/logging"
	"github.com/dihedron/archetype/settings"
	"github.com/go-git/go-git/v6/plumbing/object"
)

type Show struct {
	base.Command
}

func (cmd *Show) Execute(args []string) error {

	p := settings.Settings{
		Version: 1,
		Parameters: map[string]any{
			"key1": "value 1",
			"key2": "value 2",
			"key3": map[string]any{
				"key3.1": "value 3.1",
				"key3.2": "value 3.2",
				"key3.3": "value 3.3",
			},
			"key4": true,
			"key5": 12345,
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
