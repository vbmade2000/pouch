package main

import (
	"fmt"

	"github.com/alibaba/pouch/pkg/reference"
	"github.com/spf13/cobra"
)

var rmiDescription = "Remove one or more images by reference." +
	"When the image is being used by a container, you must specify -f to delete it." +
	"But it is strongly discouraged, because the container will be in abnormal status."

// RmiCommand use to implement 'rmi' command, it remove one or more images by reference
type RmiCommand struct {
	baseCommand
	force bool
}

// Init initialize rmi command
func (rmi *RmiCommand) Init(c *Cli) {
	rmi.cli = c
	rmi.cmd = &cobra.Command{
		Use:   "rmi image ",
		Short: "Remove one or more images by reference",
		Long:  rmiDescription,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return rmi.runRmi(args)
		},
		Example: rmiExample(),
	}
	rmi.addFlags()
}

// addFlags adds flags for specific command
func (rmi *RmiCommand) addFlags() {
	rmi.cmd.Flags().BoolVarP(&rmi.force, "force", "f", false, "if image is being used, remove image and all associated resources")
}

// runRmi is the entry of rmi command
func (rmi *RmiCommand) runRmi(args []string) error {
	apiClient := rmi.cli.Client()

	for _, name := range args {
		ref, err := reference.Parse(name)
		if err != nil {
			return fmt.Errorf("failed to remove image: %v", err)
		}

		if err := apiClient.ImageRemove(ref.String(), rmi.force); err != nil {
			return fmt.Errorf("failed to remove image: %v", err)
		}
		fmt.Printf("%s\n", ref.String())
	}

	return nil
}

// rmiExample shows examples in rmi command, and is used in auto-generated cli docs.
func rmiExample() string {
	return `$ pouch rmi registry.hub.docker.com/library/busybox:latest
registry.hub.docker.com/library/busybox:latest
$ pouch create --name test registry.hub.docker.com/library/busybox:latest
container ID: e5952417f9ee94621bbeaec532be1803ae2dedeb11a80f578a6d621e04a95afd, name: test
$ pouch rmi registry.hub.docker.com/library/busybox:latest
Error: failed to remove image: {"message":"Unable to remove the image \"registry.hub.docker.com/library/busybox:latest\" (must force) - container e5952417f9ee94621bbeaec532be1803ae2dedeb11a80f578a6d621e04a95afd is using this image"}
`
}
