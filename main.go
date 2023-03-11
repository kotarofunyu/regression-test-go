package main

import (
	"flag"
	"fmt"
	_ "image/jpeg"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kotarofunyu/regression-test-go/gitcomparison"
	"github.com/kotarofunyu/regression-test-go/urlcomparison"

	agouti "github.com/sclevine/agouti"
)

type Comparer interface {
	Run()
}

func cleanupCaptures(before, after string) {
	os.Remove(before)
	os.Remove(after)
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

func setupArgs() (baseUrl string, paths []string, breakpoints []int, gitpath, beforebranch, afterbranch, beforeurl, afterurl string) {
	b := flag.String("base_url", "", "Testing target url")
	p := flag.String("paths", "", "paths")
	bp := flag.String("breakpoints", "", "breakpoints")
	gp := flag.String("gitpath", "", "git repository path")
	bb := flag.String("beforebranch", "main", "the git branch which is base")
	ab := flag.String("afterbranch", "", "the git branch which some changes added")
	bu := flag.String("beforeurl", "", "")
	au := flag.String("afterurl", "", "")
	flag.Parse()
	baseUrl = *b
	paths = strings.Split(*p, ",")
	for _, v := range strings.Split(*bp, ",") {
		atoi, _ := strconv.Atoi(v)
		breakpoints = append(breakpoints, atoi)
	}
	beforebranch = *bb
	afterbranch = *ab
	gitpath = *gp
	beforeurl = *bu
	afterurl = *au
	return
}

func main() {
	now := time.Now()
	baseUrl, paths, breakpoints, gitpath, beforebranch, afterbranch, beforeurl, afterurl := setupArgs()
	page, driver := setupBrowser()
	defer driver.Stop()

	var comparer Comparer
	if len(gitpath) > 0 {
		comparer = gitcomparison.NewGitComparison(gitpath, beforebranch, afterbranch, baseUrl, paths, breakpoints, page)
	} else {
		comparer = urlcomparison.NewUrlComparison(beforeurl, afterurl, paths, breakpoints, page)
	}
	comparer.Run()
	fmt.Printf("Completed in: %vms\n", time.Since(now).Milliseconds())
}
