/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/kotarofunyu/regression-test-go/cmd/validator"
	"github.com/kotarofunyu/regression-test-go/comparison"
	"github.com/kotarofunyu/regression-test-go/urlcomparison"
	"github.com/spf13/cobra"
)

// diffurlCmd represents the diffurl command
var diffurlCmd = &cobra.Command{
	Use:   "diffurl",
	Short: "Comparison two websites based on urls",
	Long: `You can easily compare two websites by providing arguments.
It requires close attention that two websites must be almost same such as production env and development env. `,
	Args: func(cmd *cobra.Command, args []string) error {
		b, err := cmd.Flags().GetString("beforeurl")
		if err != nil {
			log.Fatal(err)
		}
		a, err := cmd.Flags().GetString("afterurl")
		if err != nil {
			log.Fatal(err)
		}
		err = validator.ValidateUrl(b, "beforeurl")
		if err != nil {
			return err
		}
		err = validator.ValidateUrl(a, "afterurl")
		if err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		p, d := comparison.SetupBrowser()
		defer d.Stop()
		beforeurl, err := cmd.Flags().GetString("beforeurl")
		if err != nil {
			log.Fatal(err)
		}
		afterurl, err := cmd.Flags().GetString("afterurl")
		if err != nil {
			log.Fatal(err)
		}
		paths, err := cmd.Flags().GetStringSlice("paths")
		if err != nil {
			log.Fatal(err)
		}
		breakpoints, err := cmd.Flags().GetIntSlice("breakpoints")
		if err != nil {
			log.Fatal(err)
		}
		u := urlcomparison.NewUrlComparison(beforeurl, afterurl, paths, breakpoints, p)
		u.Run()
	},
}

func init() {
	rootCmd.AddCommand(diffurlCmd)
	diffurlCmd.Flags().StringP("beforeurl", "b", "", "before url")
	diffurlCmd.Flags().StringP("afterurl", "a", "", "after url")
	diffurlCmd.Flags().StringSliceP("paths", "p", []string{}, "path")
	diffurlCmd.Flags().IntSliceP("breakpoints", "w", []int{}, "breakpoints")
	diffurlCmd.MarkFlagRequired("beforeurl")
	diffurlCmd.MarkFlagRequired("afterurl")
	diffurlCmd.MarkFlagRequired("paths")
	diffurlCmd.MarkFlagRequired("breakpoints")
}
