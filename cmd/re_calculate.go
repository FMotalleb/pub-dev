/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"os"
	"path"

	"github.com/fmotalleb/go-tools/log"
	"github.com/fmotalleb/go-tools/sysctx"
	"github.com/spf13/cobra"

	"github.com/fmotalleb/pub-dev/pub"
)

// reCalculateCmd represents the reCalculate command.
var reCalculateCmd = &cobra.Command{
	Use:   "re-calculate",
	Short: "regenerates metadata of pub packages",
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := context.Background()
		ctx = sysctx.CancelWith(ctx, os.Interrupt, os.Kill)
		var err error
		if ctx, err = log.WithNewEnvLogger(ctx); err != nil {
			return err
		}
		var storage string
		if storage, err = cmd.Flags().GetString("storage"); err != nil {
			return err
		}
		if storage == "" {
			return errors.New("storage path is mandatory")
		}
		storage = path.Join(storage, "packages")
		return pub.RecalculateMetadata(ctx, storage)
	},
}

func init() {
	rootCmd.AddCommand(reCalculateCmd)
	reCalculateCmd.Flags().StringP("storage", "s", "", "recalculate listing json for the directory")
}
