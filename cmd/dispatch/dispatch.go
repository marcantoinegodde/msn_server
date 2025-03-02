package dispatch

import (
	"github.com/spf13/cobra"
)

func BuildDispatchCmd() *cobra.Command {
	dispatchCmd := &cobra.Command{
		Use:   "dispatch",
		Short: "Start MSN dispatch server",
		Long:  `Start the MSN dispatch server, which handles incoming connections.`,
		Run: func(cmd *cobra.Command, args []string) {
			main()
		},
	}

	return dispatchCmd
}
