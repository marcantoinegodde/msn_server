package switchboard

import "github.com/spf13/cobra"

func BuildSwitchboardCmd() *cobra.Command {
	switchboardCmd := &cobra.Command{
		Use:   "switchboard",
		Short: "Start MSN switchboard server",
		Long:  `Start the MSN switchboard server, which handles users communications.`,
		Run: func(cmd *cobra.Command, args []string) {
			main()
		},
	}

	return switchboardCmd
}
