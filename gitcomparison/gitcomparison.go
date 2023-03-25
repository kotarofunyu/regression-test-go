package gitcomparison

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/kotarofunyu/regression-test-go/comparison"
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
	if err := createOutputDir("results/", "captures/"); err != nil {
		log.Fatal(err)
	}
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
				height, err := getPageHeight(gc)
				if err != nil {
					log.Fatal(err)
				}
				if err := setPageSize(gc, breakpoint, height); err != nil {
					log.Fatal(err)
				}
				if err := checkoutGitBranch(gc.repository, gc.beforebranch); err != nil {
					log.Fatal(err)
				}
				gc.page.Refresh()
				beforefilename := comparison.NewFileName("before", path, breakpoint)
				bf, err := saveCapture(beforefilename, gc)
				if err != nil {
					log.Fatal(err)
				}
				if err := checkoutGitBranch(gc.repository, gc.afterbranch); err != nil {
					log.Fatal(err)
				}
				gc.page.Refresh()
				afterfilename := comparison.NewFileName("after", path, breakpoint)
				af, err := saveCapture(afterfilename, gc)
				if err != nil {
					log.Fatal(err)
				}
				comparefunc(bf.Name(), af.Name(), path, breakpoint)
			}
		}(&wg, path)
	}
	wg.Wait()
	// NOTE: ファイルへの書き込みをやめてバイナリをメモリに保持して比較する方が省エネかも
	// defer rt.cleanupCaptures(before, after)
}

func getPageHeight(gc *GitComparison) (int, error) {
	var height int
	if err := gc.page.RunScript("return document.body.scrollHeight;", nil, &height); err != nil {
		return 0, err
	}
	return height, nil
}

func setPageSize(gc *GitComparison, breakpoint, height int) error {
	if err := gc.page.Size(breakpoint, height); err != nil {
		return err
	}
	return nil
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

func createOutputDir(resultDir, capturesDir string) error {
	err := os.Mkdir(resultDir, os.ModePerm)
	if err != nil && err.Error() != fmt.Sprintf("mkdir %s: file exists", resultDir) {
		return err
	}
	err = os.Mkdir(capturesDir, os.ModePerm)
	if err != nil && err.Error() != fmt.Sprintf("mkdir %s: file exists", capturesDir) {
		return err
	}
	return nil
}

func saveCapture(filename string, gc *GitComparison) (*os.File, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	err = gc.page.Screenshot(filename)
	if err != nil {
		return nil, err
	}
	return f, nil
}
