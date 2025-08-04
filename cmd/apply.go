package cmd

import (
	"context"
	"os"
	"time"

	"github.com/spf13/cobra"
	"xochitla.dev/baredroid/baredroid"
)

var configPath string
var restart bool

var apply = &cobra.Command{
	Use:   "apply",
	Short: "Apply a configuration to the device",
	Long: `The apply command allows you to apply a specific configuration to your Android device.`,

	Run: func(cmd *cobra.Command, args []string) {
		// initialize the device
		var device baredroid.Device = *baredroid.NewDevice(10 * time.Minute)

		var ctx context.Context = context.Background()

		cmd.Println("Applying configuration...")
		
		// Load config
		configPtr, err := baredroid.NewConfig(configPath)
		if err != nil {
			cmd.Println("Error loading config:", err)
			return
		}
		config := *configPtr

		// Execute removals
		for i:= range config.PkgRemove {
			device.RemovePackage(ctx, config.PkgRemove[i])
		}

		// Execute installations
		for _, pkg := range config.PkgInstall {

			err := device.InstallPackage(
				ctx,
				&pkg,
			)

			if err != nil {
				cmd.PrintErrf("Failed to install %w: %v\n", &pkg.Name, err)
				os.Exit(1)
			}
		
		}
	},

}

func init() {
	apply.Flags().StringVarP(&configPath, "config", "c", "", "Path to config file")

	apply.Flags().BoolVarP(&restart, "restart", "r", false, "Restart device upon successful config application")
}
