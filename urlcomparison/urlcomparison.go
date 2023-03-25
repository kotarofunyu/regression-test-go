package urlcomparison

import (
	"log"
	"os"
	"strconv"

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
	createOutputDir()
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
			before, err := saveCapture("before", path, breakpoint, uc)
			if err != nil {
				log.Fatal(err)
			}
			uc.page.Navigate(uc.afterbaseurl + path)
			uc.page.Size(breakpoint, height)
			after, err := saveCapture("after", path, breakpoint, uc)
			if err != nil {
				log.Fatal(err)
			}
			comparefunc(before, after, path, breakpoint)
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

func saveCapture(timing, path string, breakpoint int, uc *UrlComparison) (string, error) {
	dest := "./captures/" + timing + "-" + path + "-" + strconv.Itoa(breakpoint) + ".png"
	_, err := os.Create(dest)
	if err != nil {
		return "", nil
	}
	err = uc.page.Screenshot(dest)
	if err != nil {
		return "", nil
	}
	return dest, nil
}
