package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"Checkiner/pkg/checkin"
	"Checkiner/pkg/util"

	"github.com/kydance/ziwi/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultConfigDir  = "config"
	defaultConfigName = "checkiner"

	envPrefix = "CHECKINER"
)

var cfgFile string

// run is the real main entry point.
func run() error {
	_, err := json.Marshal(viper.AllSettings())
	if err != nil {
		log.Errorw("Failed to marshal viper settings", "err", err)
		return err
	}

	// Welcome
	util.SendNotify("Checkiner", "normal", "Welcome to enjoy your time with Checkiner")

	// It's time to checkin
	for {
		who, err := checkin.CheckinRun()
		// Checkiner failed
		if err != nil {
			util.SendNotify("Checkiner", "critical", who+" Check in Failed: "+err.Error())
		}
	}
}

// NewZiwiCommand creates *cobra.Command object. Then, call Execute to run application.
func NewCheckinerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "Checkiner", // Command name
		Short: "A web checkiner for clash",
		Long: `A web checkiner for clash.

Find more checkiner information at:
	https://github.com/kydance/checkiner#readme`,
		// Commands that fail to print the usage.
		SilenceUsage: false,

		// When running cmd.Execute(), it will be called.
		RunE: func(cmd *cobra.Command, args []string) error {
			// Init log
			log.Init(&log.Options{
				DisableCaller:     viper.GetBool("log.disable-caller"),
				DisableStacktrace: viper.GetBool("log.disable-stacktrace"),
				Level:             viper.GetString("log.level"),
				Format:            viper.GetString("log.format"),
				OutputPaths:       viper.GetStringSlice("log.output-paths"),
			})

			// Sync write buffer log to file
			defer log.Sync()

			return run()
		},

		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q",
						cmd.CommandPath(), args)
				}
			}

			return nil
		},
	}

	cobra.OnInitialize(initConfig)

	// Other command flags
	// ...

	// 持久性标志(PersistentFlag)，该标志可用于它所分配的命令以及该命令下的每个子命令
	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c",
		"", "The path to the ziwi configuration file. Empty string for no configuration file.")
	// 本地标志，本地标志只能在其所绑定的命令上使用
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	return cmd
}

func initConfig() {
	if cfgFile != "" {
		log.Infof("Using config file: %s", cfgFile)

		// Read config file from cfgFile.
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(filepath.Join(home, defaultConfigDir)) // $HOME/defaultConfigDir
		viper.AddConfigPath(filepath.Join(".", defaultConfigDir))  // ./defaultConfigDir

		viper.SetConfigType("yaml")
		viper.SetConfigName(defaultConfigName)
	}

	// Read matched environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read config file.
	// If a config file is specified, use it. Otherwise, search in defaultConfigDir.
	if err := viper.ReadInConfig(); err != nil {
		log.Errorw("Failed to read viper configuration file", "err", err)
	}

	log.Debugw("Using config file", "file", viper.ConfigFileUsed())
}
