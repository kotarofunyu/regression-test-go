package urlcomparison

import (
	"bytes"
	"fmt"
	"image/png"
	"log"
	"os"
	"strconv"
	"time"

	diff "github.com/olegfedoseev/image-diff"
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
			compareFiles(before, after, path, breakpoint)
		}
	}
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
