package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"

	"github.com/hashicorp/go-multierror"
	"github.com/hokaccha/go-prettyjson"
	"github.com/jodydadescott/home-dns-server/server"
	"github.com/jodydadescott/home-dns-server/static"
	"github.com/jodydadescott/home-dns-server/types"
	"github.com/jodydadescott/home-dns-server/unifi"
	"github.com/jodydadescott/jody-zap-logging/logging"
	"github.com/spf13/cobra"
)

const (
	BinaryName = "unifi-dns-server"
)

type Config = types.Config

var (
	configFileArg   string
	debugEnabledArg bool

	rootCmd = &cobra.Command{
		Use: BinaryName,
		//SilenceUsage: true,
	}

	generateConfigCmd = &cobra.Command{
		Use: "generate-config",
	}

	generateJsonConfigCmd = &cobra.Command{
		Use: "json",
		RunE: func(cmd *cobra.Command, args []string) error {
			o, _ := json.Marshal(types.NewExampleConfig())
			fmt.Println(string(o))
			return nil
		},
	}

	generatePrettyJsonConfigCmd = &cobra.Command{
		Use: "pretty-json",
		RunE: func(cmd *cobra.Command, args []string) error {
			o, _ := prettyjson.Marshal(types.NewExampleConfig())
			fmt.Println(string(o))
			return nil
		},
	}

	generateYamlConfigCmd = &cobra.Command{
		Use: "yaml",
		RunE: func(cmd *cobra.Command, args []string) error {
			o, _ := yaml.Marshal(types.NewExampleConfig())
			fmt.Println(string(o))
			return nil
		},
	}

	runCmd = &cobra.Command{

		Use: "run",

		RunE: func(cmd *cobra.Command, args []string) error {

			if configFileArg == "" {
				return fmt.Errorf("configFile is required")
			}

			config, err := getConfig(configFileArg)
			if err != nil {
				return err
			}

			if debugEnabledArg {
				zap.ReplaceGlobals(logging.GetDebugZapLogger())
				zap.L().Debug("debug is enabled")
			} else {
				zap.ReplaceGlobals(logging.GetDefaultZapLogger())
			}

			serverConfig := &server.Config{
				Debug:       debugEnabledArg,
				Listen:      config.Listen,
				Nameservers: config.Nameservers,
			}

			if config.Unifi != nil && config.Unifi.Enabled {
				zap.L().Debug("Unifi is enabled")
				serverConfig.AddProvider(unifi.New(config.Unifi))
			} else {
				zap.L().Debug("Unifi is not enabled")
			}

			if config.Static != nil && config.Static.Enabled {
				zap.L().Debug("static config is enabled")
				for _, v := range static.New(config.Static) {
					serverConfig.AddProvider(v)
				}
			} else {
				zap.L().Debug("static config is not enabled")
			}

			s := server.New(serverConfig)

			ctx, cancel := context.WithCancel(cmd.Context())

			interruptChan := make(chan os.Signal, 1)
			signal.Notify(interruptChan, os.Interrupt)

			go func() {
				select {
				case <-interruptChan: // first signal, cancel context
					cancel()
				case <-ctx.Done():
				}
				<-interruptChan // second signal, hard exit
			}()

			return s.Run(ctx)
		},
	}
)

func getConfig(configFile string) (*Config, error) {

	var errs *multierror.Error

	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var config Config

	err = json.Unmarshal(content, &config)
	if err == nil {
		return &config, nil
	}

	errs = multierror.Append(errs, err)

	err = yaml.Unmarshal(content, &config)
	if err == nil {
		return &config, nil
	}

	errs = multierror.Append(errs, err)

	return nil, errs.ErrorOrNil()
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	runCmd.PersistentFlags().StringVarP(&configFileArg, "config", "c", "", "config file")
	runCmd.PersistentFlags().BoolVarP(&debugEnabledArg, "debug", "d", false, "debug to STDERR")
	generateConfigCmd.AddCommand(generateJsonConfigCmd, generatePrettyJsonConfigCmd, generateYamlConfigCmd)
	rootCmd.AddCommand(runCmd, generateConfigCmd)
}
