package main

import (
	"bytes"
	"flag"
	"fmt"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kotarofunyu/regression-test-go/gitcomparison"
	"github.com/kotarofunyu/regression-test-go/urlcomparison"
	diff "github.com/olegfedoseev/image-diff"

	agouti "github.com/sclevine/agouti"
)

type Comparer interface {
	Run(comparefunc func(before, after, path string, breakpoint int))
}

type ComparisonTesting struct {
	comparer    Comparer
	breakpoints []int
	paths       []string
	page        *agouti.Page
}

func (ct *ComparisonTesting) Compare() {
	ct.comparer.Run(compareFiles)
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

var (
	b  = flag.String("base_url", "", "Testing target url")
	p  = flag.String("paths", "", "paths")
	bp = flag.String("breakpoints", "", "breakpoints")
	gp = flag.String("gitpath", "", "git repository path")
	bb = flag.String("beforebranch", "main", "the git branch which is base")
	ab = flag.String("afterbranch", "", "the git branch which some changes added")
	bu = flag.String("beforeurl", "", "")
	au = flag.String("afterurl", "", "")
)

func setupArgs() (baseUrl string, paths []string, breakpoints []int, gitpath, beforebranch, afterbranch, beforeurl, afterurl string) {
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

func hogemain() {
	now := time.Now()
	baseUrl, paths, breakpoints, gitpath, beforebranch, afterbranch, beforeurl, afterurl := setupArgs()
	page, driver := setupBrowser()
	defer driver.Stop()

	ct := ComparisonTesting{}
	if len(gitpath) > 0 {
		ct.comparer = gitcomparison.NewGitComparison(gitpath, beforebranch, afterbranch, baseUrl, paths, breakpoints, page)
	} else {
		ct.comparer = urlcomparison.NewUrlComparison(beforeurl, afterurl, paths, breakpoints, page)
	}
	ct.Compare()
	fmt.Printf("Completed in: %vms\n", time.Since(now).Milliseconds())
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
