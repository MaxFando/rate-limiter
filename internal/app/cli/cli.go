package cli

import (
	"context"
	"github.com/MaxFando/rate-limiter/internal/delivery/cli"
	"github.com/spf13/cobra"
	"os"
)

type Cmd struct {
	errors chan error
}

func NewCmd() *Cmd {
	return &Cmd{
		errors: make(chan error, 1),
	}
}

func (app *Cmd) Run(ctx context.Context, interruptCh chan os.Signal) {
	rootCmd := &cobra.Command{Use: "anti-bruteforce"}

	cliCmd := &cobra.Command{Use: "cli", Run: func(cmd *cobra.Command, args []string) {
		go func() {
			c := cli.New(ctx)
			c.Run(ctx, interruptCh)
		}()
	}}

	extCmd := &cobra.Command{Use: "exit", Run: func(cmd *cobra.Command, args []string) {
		interruptCh <- os.Interrupt
	}}

	rootCmd.AddCommand(cliCmd, extCmd)

	go func() {
		if err := rootCmd.Execute(); err != nil {
			app.errors <- err
		}
	}()
}

func (app *Cmd) Notify() <-chan error {
	return app.errors
}
