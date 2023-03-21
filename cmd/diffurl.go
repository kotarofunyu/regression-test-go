/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"image/png"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/kotarofunyu/regression-test-go/urlcomparison"
	diff "github.com/olegfedoseev/image-diff"
	"github.com/sclevine/agouti"
	"github.com/spf13/cobra"
)

// diffurlCmd represents the diffurl command
var diffurlCmd = &cobra.Command{
	Use:   "diffurl",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("diffurl called")
		p, d := setupBrowser()
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
		u.Run(compareFiles)
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

func setupBrowser() (*agouti.Page, *agouti.WebDriver) {
	driver := agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{
			"--headless",
		}),
	)
	err := driver.Start()
	if err != nil {
		log.Fatal(err)
	}
	page, _ := driver.NewPage()

	return page, driver
}

func compareFiles(before, after, path string, breakpoint int) {
	diff, percent, err := diff.CompareFiles(before, after)
	if err != nil {
		log.Fatal(err, before, after)
	}
	if percent == 0.0 {
		fmt.Println("Image is same!")
		return
	}
	t := time.Now()
	ft := t.Format("20200101123045")
	diffName := "diff-" + path + "-" + strconv.Itoa(breakpoint) + "px" + "-" + ft + ".png"
	destDir := "./results/"
	f, err := os.Create(destDir + diffName)
	if err != nil {
		log.Fatal(err)
	}
	buf := new(bytes.Buffer)
	png.Encode(buf, diff)
	f.Write(buf.Bytes())
	fmt.Println("diff has written into " + destDir + diffName)
}
