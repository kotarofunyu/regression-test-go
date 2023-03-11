package gitcomparison

import (
	"bytes"
	"fmt"
	"image/png"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	diff "github.com/olegfedoseev/image-diff"
	"github.com/sclevine/agouti"
)

type GitComparison struct {
	repository   *git.Worktree
	beforebranch string
	afterbranch  string
	baseurl      string
	paths        []string
	breakpoints  []int
	page         *agouti.Page
}

func NewGitComparison(gitpath, beforebranch, afterbranch, baseUrl string, paths []string, breakpoints []int, page *agouti.Page) *GitComparison {
	r, err := git.PlainOpen(gitpath)
	if err != nil {
		log.Fatal(err)
	}
	wt, err := r.Worktree()
	if err != nil {
		log.Fatal(err)
	}
	return &GitComparison{
		repository:   wt,
		beforebranch: beforebranch,
		afterbranch:  afterbranch,
		baseurl:      baseUrl,
		paths:        paths,
		breakpoints:  breakpoints,
		page:         page,
	}
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
