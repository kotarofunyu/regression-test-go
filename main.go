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
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	diff "github.com/olegfedoseev/image-diff"
	agouti "github.com/sclevine/agouti"
)

type RegressionTest struct {
	testConfig TestConfig
	page       *agouti.Page
	repository *git.Worktree
}

type TestConfig struct {
	breakpoints []int
	baseurl     string
	paths       []string
	initheight  int
	gitconf     GitConfig
}

type GitConfig struct {
	path         string
	beforebranch string
	afterbranch  string
}

func (rt *RegressionTest) Run() {
	os.Mkdir("results/", os.ModePerm)
	os.Mkdir("captures/", os.ModePerm)
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, path := range rt.testConfig.paths {
		wg.Add(1)
		go func(wg *sync.WaitGroup, path string) {
			// NOTE: goroutine間でagouti.Pageを共有するので排他制御が必要
			mu.Lock()
			defer wg.Done()
			defer mu.Unlock()
			rt.page.Navigate(rt.testConfig.baseurl + path)
			for _, breakpoint := range rt.testConfig.breakpoints {
				before := "./captures/before-" + path + "-" + strconv.Itoa(breakpoint) + ".png"
				after := "./captures/after-" + path + "-" + strconv.Itoa(breakpoint) + ".png"
				os.Create(before)
				os.Create(after)
				var height int
				if err := rt.page.RunScript("return document.body.scrollHeight;", nil, &height); err != nil {
					log.Fatal(err)
				}
				rt.page.Size(breakpoint, height)
				if err := checkoutGitBranch(rt.repository, rt.testConfig.gitconf.beforebranch); err != nil {
					log.Fatal(err)
				}
				rt.page.Refresh()
				rt.page.Screenshot(before)
				if err := checkoutGitBranch(rt.repository, rt.testConfig.gitconf.afterbranch); err != nil {
					log.Fatal(err)
				}
				rt.page.Refresh()
				rt.page.Screenshot(after)
				rt.compareFiles(before, after, path, breakpoint)
			}
		}(&wg, path)
	}
	wg.Wait()
	// NOTE: ファイルへの書き込みをやめてバイナリをメモリに保持して比較する方が省エネかも
	// defer rt.cleanupCaptures(before, after)
}

func (rt *RegressionTest) cleanupCaptures(before, after string) {
	os.Remove(before)
	os.Remove(after)
}

func (rt *RegressionTest) compareFiles(before, after, path string, breakpoint int) {
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

func setupArgs() (baseUrl string, paths []string, breakpoints []int, gitpath, beforebranch, afterbranch string) {
	b := flag.String("base_url", "", "Testing target url")
	p := flag.String("paths", "", "paths")
	bp := flag.String("breakpoints", "", "breakpoints")
	gp := flag.String("gitpath", "", "git repository path")
	bb := flag.String("beforebranch", "main", "the git branch which is base")
	ab := flag.String("afterbranch", "", "the git branch which some changes added")
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
	return
}

func checkoutGitBranch(wt *git.Worktree, destbranch string) error {
	err := wt.Checkout(
		&git.CheckoutOptions{
			Branch: plumbing.NewBranchReferenceName(destbranch),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	now := time.Now()
	baseUrl, paths, breakpoints, gitpath, beforebranch, afterbranch := setupArgs()
	page, driver := setupBrowser()
	defer driver.Stop()
	gitconf := GitConfig{
		path:         gitpath,
		beforebranch: beforebranch,
		afterbranch:  afterbranch,
	}
	mytestconf := TestConfig{
		breakpoints: breakpoints,
		baseurl:     baseUrl,
		paths:       paths,
		initheight:  300,
		gitconf:     gitconf,
	}
	r, err := git.PlainOpen(gitconf.path)
	if err != nil {
		log.Fatal(err)
	}
	wt, err := r.Worktree()
	if err != nil {
		log.Fatal(err)
	}
	rt := RegressionTest{
		mytestconf,
		page,
		wt,
	}
	rt.Run()
	fmt.Printf("Completed in: %vms\n", time.Since(now).Milliseconds())
}
