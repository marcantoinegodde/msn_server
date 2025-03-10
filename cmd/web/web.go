package web

import "github.com/spf13/cobra"

func BuildWebCmd() *cobra.Command {
	webCmd := &cobra.Command{
		Use:   "web",
		Short: "Start web server",
		Long:  `Start the web server, which serves the account management application.`,
		Run: func(cmd *cobra.Command, args []string) {
			main()
		},
	}

	return webCmd
}
