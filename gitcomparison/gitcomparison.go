package gitcomparison

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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

func (gc *GitComparison) Run(comparefunc func(before, after, path string, breakpoint int)) {
	createOutputDir()
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
				var height int
				if err := gc.page.RunScript("return document.body.scrollHeight;", nil, &height); err != nil {
					log.Fatal(err)
				}
				gc.page.Size(breakpoint, height)
				if err := checkoutGitBranch(gc.repository, gc.beforebranch); err != nil {
					log.Fatal(err)
				}
				gc.page.Refresh()
				before, err := saveCapture("before", path, breakpoint, gc)
				if err != nil {
					log.Fatal(err)
				}
				if err := checkoutGitBranch(gc.repository, gc.afterbranch); err != nil {
					log.Fatal(err)
				}
				gc.page.Refresh()
				after, err := saveCapture("after", path, breakpoint, gc)
				if err != nil {
					log.Fatal(err)
				}
				comparefunc(before, after, path, breakpoint)
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

func createOutputDir() error {
	err := os.Mkdir("results/", os.ModePerm)
	if err != nil {
		return err
	}
	err = os.Mkdir("captures/", os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func saveCapture(timing, path string, breakpoint int, gc *GitComparison) (string, error) {
	dest := "./captures/" + timing + "-" + path + "-" + strconv.Itoa(breakpoint) + ".png"
	_, err := os.Create(dest)
	if err != nil {
		return "", err
	}
	err = gc.page.Screenshot(dest)
	if err != nil {
		return "", err
	}
	return dest, nil
}
