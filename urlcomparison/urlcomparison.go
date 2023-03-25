package urlcomparison

import (
	"fmt"
	"log"
	"os"

	"github.com/sclevine/agouti"
)

type UrlComparison struct {
	beforebaseurl string
	afterbaseurl  string
	paths         []string
	breakpoints   []int
	page          *agouti.Page
}

func NewUrlComparison(beforebaseurl, afterbaseurl string, paths []string, breakpoints []int, page *agouti.Page) *UrlComparison {
	return &UrlComparison{
		beforebaseurl: beforebaseurl,
		afterbaseurl:  afterbaseurl,
		paths:         paths,
		breakpoints:   breakpoints,
		page:          page,
	}
}

func (uc *UrlComparison) Run(comparefunc func(before, after, path string, breakpoint int)) {
	if err := createOutputDir("results/", "captures/"); err != nil {
		log.Fatal(err)
	}
	for _, path := range uc.paths {
		for _, breakpoint := range uc.breakpoints {
			uc.page.Navigate(uc.beforebaseurl + path)
			height, err := getPageHeight(uc)
			if err != nil {
				log.Fatal(err)
			}
			if err := setPageSize(uc, breakpoint, height); err != nil {
				log.Fatal(err)
			}
			beforefilename := newFileName("before", path, breakpoint)
			bf, err := saveCapture(beforefilename, uc)
			if err != nil {
				log.Fatal(err)
			}
			uc.page.Navigate(uc.afterbaseurl + path)
			uc.page.Size(breakpoint, height)
			afterfilename := newFileName("after", path, breakpoint)
			af, err := saveCapture(afterfilename, uc)
			if err != nil {
				log.Fatal(err)
			}
			comparefunc(bf.Name(), af.Name(), path, breakpoint)
		}
	}
}

func getPageHeight(uc *UrlComparison) (int, error) {
	var height int
	if err := uc.page.RunScript("return document.body.scrollHeight;", nil, &height); err != nil {
		return 0, err
	}
	return height, nil
}

func setPageSize(uc *UrlComparison, breakpoint, height int) error {
	if err := uc.page.Size(breakpoint, height); err != nil {
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

func saveCapture(filename string, uc *UrlComparison) (*os.File, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	err = uc.page.Screenshot(filename)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func newFileName(timing, path string, breakpoint int) string {
	return fmt.Sprintf("./captures/%s-%s-%d.png", timing, path, breakpoint)
}
