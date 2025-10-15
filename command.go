package main

// these are the single command options
type Options struct {
	Repository string `short:"r" long:"repository" description:"The Git repository to clone" required:"true"`
	Tag        string `short:"t" long:"tag" description:"The tag to clone" optional:"true" default:"HEAD"`
	Settings   string `short:"s" long:"settings" description:"The repository-specific settings" required:"true"`
	Auth       struct {
		Token         *string `short:"T" long:"token" description:"The personal access token for authentication" optional:"true"`
		Username      *string `short:"U" long:"username" description:"The username for authentication" optional:"true"`
		Password      *string `short:"P" long:"password" description:"The password for authentication" optional:"true"`
		SSHKey        *string `short:"K" long:"sshkey" description:"The SSH key for authentication" optional:"true"`
		DefaultSSHKey bool    `short:"D" long:"default-sshkey" description:"Use default SSH key for authentication" optional:"true"`
		SSHAgent      bool    `short:"A" long:"sshagent" description:"Use SSH agent for authentication" optional:"true"`
	} `group:"Authentication Options" description:"Authentication options"`
}
