package main

import (
	"bytes"
	"fmt"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
	"strconv"

	diff "github.com/olegfedoseev/image-diff"
	agouti "github.com/sclevine/agouti"
)

type TestPage struct {
	*agouti.Page
}

func (tp *TestPage) CapturePage(filename string) {
	var width, height int
	if err := tp.RunScript("return document.body.scrollHeight;", nil, &height); err != nil {
		log.Fatal(err)
	}
	// TODO: widthはブレイクポイントの値を使うのでここは本当は不要
	if err := tp.RunScript("return document.body.scrollWidth;", nil, &width); err != nil {
		log.Fatal(err)
	}
	tp.Size(width, height)
	tp.Screenshot(filename)
}

type RegressionTest struct {
	testConfig TestConfig
	page       *agouti.Page
	tp         TestPage
}

type TestConfig struct {
	breakpoints []int
	baseurl     string
	paths       []string
	initheight  int
}

func (rt *RegressionTest) Run() {
	for _, breakpoint := range rt.testConfig.breakpoints {
		rt.page.Size(breakpoint, rt.testConfig.initheight)
		for _, path := range rt.testConfig.paths {
			rt.page.Navigate(rt.testConfig.baseurl + path)
			before := "./captures/before-" + path + "-" + strconv.Itoa(breakpoint) + ".png"
			after := "./captures/after-" + path + "-" + strconv.Itoa(breakpoint) + ".png"
			os.Create(before)
			os.Create(after)
			rt.tp.CapturePage(before)
			rt.tp.CapturePage(after)
			compareFiles(before, after, path, breakpoint)
		}
	}
}

func main() {
	driver := agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{
			"--headless",
		}),
	)
	err := driver.Start()
	if err != nil {
		log.Fatal(err)
	}
	defer driver.Stop()
	page, _ := driver.NewPage()
	tp := TestPage{page}
	mytestconf := TestConfig{
		breakpoints: []int{1200, 768, 384},
		baseurl:     "http://localhost:8000/",
		paths:       []string{"", "", ""},
		initheight:  300,
	}
	rt := RegressionTest{
		mytestconf,
		page,
		tp,
	}
	rt.Run()
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
	diffName := "diff-" + path + "-" + strconv.Itoa(breakpoint) + ".png"
	f, err := os.Create(diffName)
	if err != nil {
		log.Fatal(err)
	}
	buf := new(bytes.Buffer)
	png.Encode(buf, diff)
	f.Write(buf.Bytes())
	fmt.Println("diff has written into" + diffName)
}
