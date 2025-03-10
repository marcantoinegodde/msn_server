package cmd

import (
	"fmt"
	"msnserver/cmd/dispatch"
	"msnserver/cmd/notification"
	"msnserver/cmd/switchboard"
	"msnserver/cmd/web"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "msnserver",
	Short: "MSNServer is a MSNP server",
	Long: `An implementation written in Golang of Microsoft's MSN Messenger server,
supporting the MSNP protocol.`,
}

func init() {
	rootCmd.AddCommand(dispatch.BuildDispatchCmd())
	rootCmd.AddCommand(notification.BuildNotificationCmd())
	rootCmd.AddCommand(switchboard.BuildSwitchboardCmd())
	rootCmd.AddCommand(web.BuildWebCmd())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
