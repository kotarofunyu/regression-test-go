package urlcomparison

import (
	"log"

	"github.com/kotarofunyu/regression-test-go/comparison"
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

func (uc *UrlComparison) Run() {
	if err := comparison.CreateOutputDir("results/", "captures/"); err != nil {
		log.Fatal(err)
	}
	for _, path := range uc.paths {
		for _, breakpoint := range uc.breakpoints {
			uc.page.Navigate(uc.beforebaseurl + path)
			height, err := comparison.GetPageHeight(uc.page)
			if err != nil {
				log.Fatal(err)
			}
			if err := comparison.SetPageSize(uc.page, breakpoint, height); err != nil {
				log.Fatal(err)
			}
			beforefilename := comparison.NewFileName("before", path, breakpoint)
			bf, err := comparison.SaveCapture(beforefilename, uc.page)
			if err != nil {
				log.Fatal(err)
			}
			uc.page.Navigate(uc.afterbaseurl + path)
			uc.page.Size(breakpoint, height)
			afterfilename := comparison.NewFileName("after", path, breakpoint)
			af, err := comparison.SaveCapture(afterfilename, uc.page)
			if err != nil {
				log.Fatal(err)
			}
			comparison.CompareFiles(bf.Name(), af.Name(), path, breakpoint)
		}
	}
}
