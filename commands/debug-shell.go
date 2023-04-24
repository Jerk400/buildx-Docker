package commands

import (
	"context"
	"os"
	"runtime"

	"github.com/containerd/console"
	"github.com/docker/buildx/controller"
	"github.com/docker/buildx/controller/control"
	controllerapi "github.com/docker/buildx/controller/pb"
	"github.com/docker/buildx/monitor"
	"github.com/docker/cli/cli/command"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func debugShellCmd(dockerCli command.Cli) *cobra.Command {
	var options control.ControlOptions
	var progress string

	cmd := &cobra.Command{
		Use:   "debug-shell",
		Short: "Start a monitor",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.TODO()
			c, err := controller.NewController(ctx, options, dockerCli)
			if err != nil {
				return err
			}
			defer func() {
				if err := c.Close(); err != nil {
					logrus.Warnf("failed to close server connection %v", err)
				}
			}()
			con := console.Current()
			if err := con.SetRaw(); err != nil {
				return errors.Errorf("failed to configure terminal: %v", err)
			}
			err = monitor.RunMonitor(ctx, "", nil, controllerapi.InvokeConfig{
				Tty: true,
			}, c, progress, os.Stdin, os.Stdout, os.Stderr)
			con.Reset()
			return err
		},
	}

	flags := cmd.Flags()

	flags.StringVar(&options.Root, "root", "", "Specify root directory of server to connect [experimental]")
	flags.BoolVar(&options.Detach, "detach", runtime.GOOS == "linux", "Detach buildx server (supported only on linux) [experimental]")
	flags.StringVar(&options.ServerConfig, "server-config", "", "Specify buildx server config file (used only when launching new server) [experimental]")
	flags.StringVar(&progress, "progress", "auto", `Set type of progress output ("auto", "plain", "tty"). Use plain to show container output`)

	return cmd
}

func addDebugShellCommand(cmd *cobra.Command, dockerCli command.Cli) {
	cmd.AddCommand(
		debugShellCmd(dockerCli),
	)
}
