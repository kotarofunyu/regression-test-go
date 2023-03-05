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

	diff "github.com/olegfedoseev/image-diff"
	agouti "github.com/sclevine/agouti"
)

type RegressionTest struct {
	testConfig TestConfig
	page       *agouti.Page
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
			var wg sync.WaitGroup
			wg.Add(2)
			go rt.hoge(before, breakpoint, &wg)
			go rt.hoge(after, breakpoint, &wg)
			wg.Wait()
			// NOTE: ファイルへの書き込みをやめてバイナリをメモリに保持して比較する方が省エネかも
			defer rt.cleanupCaptures(before, after)
			fmt.Println("compare files")
			rt.compareFiles(before, after, path, breakpoint)
		}
	}
}

func (rt *RegressionTest) hoge(filename string, breakpoint int, wg *sync.WaitGroup) {
	os.Create(filename)
	rt.capturePage(filename, breakpoint)
	wg.Done()
}

func (rt *RegressionTest) cleanupCaptures(before, after string) {
	os.Remove(before)
	os.Remove(after)
}

func (rt *RegressionTest) capturePage(filename string, width int) {
	var height int
	if err := rt.page.RunScript("return document.body.scrollHeight;", nil, &height); err != nil {
		log.Fatal(err)
	}
	rt.page.Size(width, height)
	rt.page.Screenshot(filename)
}

func (rt *RegressionTest) compareFiles(before, after, path string, breakpoint int) {
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

func setupArgs() (baseUrl string, paths []string, breakpoints []int) {
	b := flag.String("base_url", "", "Testing target url")
	p := flag.String("paths", "", "paths")
	bp := flag.String("breakpoints", "", "breakpoints")
	flag.Parse()
	baseUrl = *b
	paths = strings.Split(*p, ",")
	for _, v := range strings.Split(*bp, ",") {
		atoi, _ := strconv.Atoi(v)
		breakpoints = append(breakpoints, atoi)
	}
	return baseUrl, paths, breakpoints
}

func main() {
	now := time.Now()
	baseUrl, paths, breakpoints := setupArgs()
	page, driver := setupBrowser()
	defer driver.Stop()
	mytestconf := TestConfig{
		breakpoints: breakpoints,
		baseurl:     baseUrl,
		paths:       paths,
		initheight:  300,
	}
	rt := RegressionTest{
		mytestconf,
		page,
	}
	rt.Run()
	fmt.Printf("Completed in: %vms\n", time.Since(now).Milliseconds())
}
