package main

import (
	"bytes"
	"fmt"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
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
	bp := []int{1200, 768, 384}
	ul := []string{"hoge", "fuga", "foo", "bar"}
	baseurl := "http://localhost:8000/"

	initHight := 300
	for _, breakpoint := range bp {
		page.Size(breakpoint, initHight)
		for _, path := range ul {
			page.Navigate(baseurl + path)
			before := "./captures/before-" + path + ".png"
			after := "./captures/after-" + path + ".png"
			os.Create(before)
			os.Create(after)
			tp.CapturePage(before)
			tp.CapturePage(after)
			compareFiles(before, after, path, breakpoint)
		}
	}
}

func compareFiles(before, after, path string, breakpoint int) {
	diff, percent, err := diff.CompareFiles(before, after)
	if err != nil {
		log.Fatal(err)
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
