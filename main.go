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

type Comparer interface {
	Run()
}

type GitComparison struct {
	repository   *git.Worktree
	beforebranch string
	afterbranch  string
	baseurl      string
	paths        []string
	initheight   int
	breakpoints  []int
	page         *agouti.Page
}

type UrlComparison struct {
	beforebaseurl string
	afterbaseurl  string
	paths         []string
	initheight    int
	breakpoints   []int
	page          *agouti.Page
}

func (gc *GitComparison) Run() {
	os.Mkdir("results/", os.ModePerm)
	os.Mkdir("captures/", os.ModePerm)
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, path := range gc.paths {
		wg.Add(1)
		go func(wg *sync.WaitGroup, path string) {
			// NOTE: goroutine間でagouti.Pageを共有するので排他制御が必要
			mu.Lock()
			defer wg.Done()
			defer mu.Unlock()
			gc.page.Navigate(gc.baseurl + path)
			for _, breakpoint := range gc.breakpoints {
				before := "./captures/before-" + path + "-" + strconv.Itoa(breakpoint) + ".png"
				after := "./captures/after-" + path + "-" + strconv.Itoa(breakpoint) + ".png"
				os.Create(before)
				os.Create(after)
				var height int
				if err := gc.page.RunScript("return document.body.scrollHeight;", nil, &height); err != nil {
					log.Fatal(err)
				}
				gc.page.Size(breakpoint, height)
				if err := checkoutGitBranch(gc.repository, gc.beforebranch); err != nil {
					log.Fatal(err)
				}
				gc.page.Refresh()
				gc.page.Screenshot(before)
				if err := checkoutGitBranch(gc.repository, gc.afterbranch); err != nil {
					log.Fatal(err)
				}
				gc.page.Refresh()
				gc.page.Screenshot(after)
				compareFiles(before, after, path, breakpoint)
			}
		}(&wg, path)
	}
	wg.Wait()
	// NOTE: ファイルへの書き込みをやめてバイナリをメモリに保持して比較する方が省エネかも
	// defer rt.cleanupCaptures(before, after)
}

func cleanupCaptures(before, after string) {
	os.Remove(before)
	os.Remove(after)
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
	r, err := git.PlainOpen(gitpath)
	if err != nil {
		log.Fatal(err)
	}
	wt, err := r.Worktree()
	if err != nil {
		log.Fatal(err)
	}
	gc := GitComparison{
		repository:   wt,
		beforebranch: beforebranch,
		afterbranch:  afterbranch,
		baseurl:      baseUrl,
		paths:        paths,
		initheight:   300,
		breakpoints:  breakpoints,
		page:         page,
	}
	gc.Run()
	fmt.Printf("Completed in: %vms\n", time.Since(now).Milliseconds())
}
