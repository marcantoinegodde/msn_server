package notification

import "github.com/spf13/cobra"

func BuildNotificationCmd() *cobra.Command {
	notificationCmd := &cobra.Command{
		Use:   "notification",
		Short: "Start MSN notification server",
		Long:  `Start the MSN notification server, which handles users sessions.`,
		Run: func(cmd *cobra.Command, args []string) {
			main()
		},
	}

	return notificationCmd
}
