package cmd

import (
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
			device.RemovePackage(config.PkgRemove[i])
		}

		// Execute installations
		for i:= range config.PkgInstall {

			err := device.InstallPackage(
				config.PkgInstall[i].Name,
				config.PkgInstall[i].Package,
				config.PkgInstall[i].Source,
				baredroid.InstallType(config.PkgInstall[i].Type),
			)

			if err != nil {
				cmd.PrintErrf("Failed to install %s: %v\n", config.PkgInstall[i].Name, err)
				os.Exit(1)
			}
		
		}
	},

}

func init() {
	apply.Flags().StringVarP(&configPath, "config", "c", "", "Path to config file")

	apply.Flags().BoolVarP(&restart, "restart", "r", false, "Restart device upon successful config application")
}

// func ExecCmd(cmd *exec.Cmd) {

// 	err := cmd.Run()
// 	if err != nil {

// 		fmt.Println("Could not run command: ", cmd, err)
// 	}

// 	// fmt.Println(string(out))
// }

// func DownloadAPKFromURL(filepath string, url string) error {

// 	out, err := os.Create(filepath)
// 	if err != nil {
// 		return err
// 	}
// 	defer out.Close()

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	// Raise for status
// 	if resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("bad status: %s", resp.Status)
// 	}

// 	_, err = io.Copy(out, resp.Body)
// 	if err != nil {
// 		return err
// 	}

// 	return nil

// }