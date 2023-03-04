package main

import (
	"bytes"
	"fmt"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"os"

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
	page.Navigate("http://localhost:8000/")
	before := "./captures/before.png"
	after := "./captures/after.png"
	os.Create(before)
	os.Create(after)
	tp.CapturePage(before)
	tp.CapturePage(after)
	compareFiles(before, after)
}

func compareFiles(before, after string) {
	diff, percent, err := diff.CompareFiles(before, after)
	if err != nil {
		log.Fatal(err)
	}
	if percent == 0.0 {
		fmt.Println("Image is same!")
		return
	}
	f, err := os.Create("diff.png")
	if err != nil {
		log.Fatal(err)
	}
	buf := new(bytes.Buffer)
	png.Encode(buf, diff)
	f.Write(buf.Bytes())
	fmt.Println("diff has written into diff.png")
}
