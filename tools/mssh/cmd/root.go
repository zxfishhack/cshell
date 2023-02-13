package cmd

import (
	"github.com/zxfishhack/cshell/pkg/mssh"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var tagList []string
var hostList []string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mssh [-t tag1,tag2...] [--host host1,host2...] -- COMMAND [args...]",
	Short: "execute same command in multiple host",
	Long:  `execute same command in multiple host`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(tagList) == 0 && len(hostList) == 0 {
			log.Fatal("tag list or host list must be special.")
		}
		err = mssh.Init()
		if err != nil {
			return
		}
		password, _ := cmd.Flags().GetString("password")
		return mssh.Execute(tagList, hostList, password, args)
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mssh.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringArrayVarP(&tagList, "tag", "t", nil, "tag list for host filter")
	rootCmd.Flags().StringArrayVar(&hostList, "host", nil, "host list for host filter")
	rootCmd.Flags().StringP("password", "p", "", "password for identity file or host")
}
