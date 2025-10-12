package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/dihedron/bootstrap/metadata"
	"github.com/dihedron/bootstrap/repository"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Repository string `short:"r" long:"repository" description:"repository to clone" required:"true"`
	//Settings   string `short:"s" long:"settings" required:"true"`
}

func (o *Options) Execute(args []string) error {
	fmt.Println("in execute")
	return nil
}

func main() {

	defer cleanup()

	if len(os.Args) == 2 && (os.Args[1] == "version" || os.Args[1] == "--version") {
		metadata.Print(os.Stdout)
		os.Exit(0)
	} else if len(os.Args) == 3 && os.Args[1] == "version" && (os.Args[2] == "--verbose" || os.Args[2] == "-v") {
		metadata.PrintFull(os.Stdout)
		os.Exit(0)
	}

	// these are the single command options
	var options struct {
		Repository string `short:"r" long:"repository" description:"The Git repository to clone" required:"true"`
		Tag        string `short:"t" long:"tag" description:"The tag to clone" optional:"true" default:"HEAD"`
		Auth       struct {
			Token         *string `short:"T" long:"token" description:"The personal access token for authentication" optional:"true"`
			Username      *string `short:"U" long:"username" description:"The username for authentication" optional:"true"`
			Password      *string `short:"P" long:"password" description:"The password for authentication" optional:"true"`
			SSHKey        *string `short:"K" long:"sshkey" description:"The SSH key for authentication" optional:"true"`
			DefaultSSHKey bool    `short:"D" long:"default-sshkey" description:"Use default SSH key for authentication" optional:"true"`
			SSHAgent      bool    `short:"A" long:"sshagent" description:"Use SSH agent for authentication" optional:"true"`
		} `group:"Authentication Options" description:"Authentication options"`
	}
	if _, err := flags.Parse(&options); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			//fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		default:
			//fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}

	var auth repository.Option
	if options.Auth.Token != nil {
		if strings.HasPrefix(options.Repository, "http") {
			slog.Info("using token for authentication")
			auth = repository.WithTokenAuth(*options.Auth.Token)
		} else {
			slog.Error("token authentication is only supported for HTTP repositories")
			os.Exit(1)
		}
	} else if options.Auth.Password != nil && options.Auth.Username != nil {
		if strings.HasPrefix(options.Repository, "http") {
			slog.Info("using username and password for authentication")
			auth = repository.WithBasicAuth(*options.Auth.Username, *options.Auth.Password)
		} else {
			slog.Error("username and password authentication is only supported for HTTP repositories")
			os.Exit(1)
		}
	} else if options.Auth.SSHKey != nil {
		if strings.HasPrefix(options.Repository, "ssh://") || strings.HasPrefix(options.Repository, "git@") {
			slog.Info("using SSH key for authentication")
			auth = repository.WithSSHKey(*options.Auth.SSHKey, nil)
		} else {
			slog.Error("SSH key authentication is only supported for SSH repositories")
			os.Exit(1)
		}
	} else if options.Auth.DefaultSSHKey {
		if strings.HasPrefix(options.Repository, "ssh://") || strings.HasPrefix(options.Repository, "git@") {
			slog.Info("using default SSH key for authentication")
			auth = repository.WithDefaultSSHKey()
		} else {
			slog.Error("SSH key authentication is only supported for SSH repositories")
			os.Exit(1)
		}
	} else if options.Auth.SSHAgent {
		if strings.HasPrefix(options.Repository, "ssh://") || strings.HasPrefix(options.Repository, "git@") {
			slog.Info("using SSH agent for authentication")
			auth = repository.WithSSHAgent()
		} else {
			slog.Error("SSH agent authentication is only supported for SSH repositories")
			os.Exit(1)
		}
	} else {
		slog.Info("using anonymous authentication")
		auth = nil
	}

	repo := repository.New(
		options.Repository,
		auth,
	)
	repo.Clone()

	var reference *plumbing.Reference
	if options.Tag == "HEAD" {
		var err error
		// ... retrieving the branch being pointed by HEAD
		reference, err = repo.Head()
		if err != nil {
			slog.Error("failed to get HEAD", "error", err)
			os.Exit(1)
		}
		fmt.Println("HEAD points to:", reference.Name())
	} else {
		var err error
		reference, err = repo.Tag(options.Tag)
		if err != nil {
			slog.Error("failed to get tag", "error", err)
			os.Exit(1)
		}
	}

	repo.ForEachFile(reference, VisitFile)
}

func VisitFile(file *object.File) error {
	fmt.Printf("%v  %9d  %s    %s\n", file.Mode, file.Size, file.Hash.String(), file.Name)

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

	return nil
}
