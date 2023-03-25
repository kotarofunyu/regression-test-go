/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/kotarofunyu/regression-test-go/cmd/validator"
	"github.com/kotarofunyu/regression-test-go/comparison"
	"github.com/kotarofunyu/regression-test-go/gitcomparison"
	"github.com/spf13/cobra"
)

// diffbranchCmd represents the diffbranch command
var diffbranchCmd = &cobra.Command{
	Use:   "diffbranch",
	Short: "Comparison two websites based on git branches",
	Long: `You can easily compare two git branches by providing arguments.
It requires close attention that two websites must be almost same such as production env and development env. `,
	Args: func(cmd *cobra.Command, args []string) error {
		u, err := cmd.Flags().GetString("url")
		fmt.Println(u)
		if err != nil {
			log.Fatal(err)
		}
		err = validator.ValidateUrl(u, "url")
		if err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("diffbranch called")
		gitdir, err := cmd.Flags().GetString("gitdir")
		if err != nil {
			log.Fatal(err)
		}
		beforebranch, err := cmd.Flags().GetString("beforebranch")
		if err != nil {
			log.Fatal(err)
		}
		afterbranch, err := cmd.Flags().GetString("afterbranch")
		if err != nil {
			log.Fatal(err)
		}
		url, err := cmd.Flags().GetString("url")
		if err != nil {
			log.Fatal(err)
		}
		paths, err := cmd.Flags().GetStringSlice("paths")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(paths)
		breakpoints, err := cmd.Flags().GetIntSlice("breakpoints")
		if err != nil {
			log.Fatal(err)
		}
		p, d := comparison.SetupBrowser()
		defer d.Stop()
		gc := gitcomparison.NewGitComparison(gitdir, beforebranch, afterbranch, url, paths, breakpoints, p)
		gc.Run()
	},
}

func init() {
	rootCmd.AddCommand(diffbranchCmd)
	diffbranchCmd.Flags().StringP("gitdir", "d", "", "directory that git repository is placed")
	diffbranchCmd.Flags().StringP("beforebranch", "b", "", "before branch")
	diffbranchCmd.Flags().StringP("afterbranch", "a", "", "after branch")
	diffbranchCmd.Flags().StringP("url", "u", "", "url")
	diffbranchCmd.Flags().StringSliceP("paths", "p", []string{}, "paths")
	diffbranchCmd.Flags().IntSliceP("breakpoints", "w", []int{}, "breakpoints")
	diffbranchCmd.MarkFlagRequired("gitdir")
	diffbranchCmd.MarkFlagRequired("beforebranch")
	diffbranchCmd.MarkFlagRequired("afterbranch")
	diffbranchCmd.MarkFlagRequired("url")
	diffbranchCmd.MarkFlagRequired("paths")
	diffbranchCmd.MarkFlagRequired("breakpoints")
}
