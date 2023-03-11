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
	os.Mkdir("results/", os.ModePerm)
	os.Mkdir("captures/", os.ModePerm)
	for _, path := range uc.paths {
		for _, breakpoint := range uc.breakpoints {
			before := "./captures/before-" + path + "-" + strconv.Itoa(breakpoint) + ".png"
			os.Create(before)
			after := "./captures/after-" + path + "-" + strconv.Itoa(breakpoint) + ".png"
			os.Create(after)
			uc.page.Navigate(uc.beforebaseurl + path)
			var height int
			if err := uc.page.RunScript("return document.body.scrollHeight;", nil, &height); err != nil {
				log.Fatal(err)
			}
			uc.page.Size(breakpoint, height)
			uc.page.Screenshot(before)
			uc.page.Navigate(uc.afterbaseurl + path)
			uc.page.Size(breakpoint, height)
			uc.page.Screenshot(after)
			comparefunc(before, after, path, breakpoint)
		}
	}
}
