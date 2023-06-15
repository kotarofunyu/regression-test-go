package comparison

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

type Comparison struct {
	before func()
	after  func()
}

func SetupBrowser() (*agouti.Page, *agouti.WebDriver) {
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

func NewFileName(timing, path string, breakpoint int) string {
	return fmt.Sprintf("./captures/%s-%s-%d.png", timing, path, breakpoint)
}

func SaveCapture(filename string, p *agouti.Page) (*os.File, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	err = p.Screenshot(filename)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func CreateOutputDir(resultDir, capturesDir string) error {
	err := os.Mkdir(resultDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}
	err = os.Mkdir(capturesDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

func GetPageHeight(p *agouti.Page) (int, error) {
	var height int
	if err := p.RunScript("return document.body.scrollHeight;", nil, &height); err != nil {
		return 0, err
	}
	return height, nil
}

func SetPageSize(p *agouti.Page, breakpoint, height int) error {
	if err := p.Size(breakpoint, height); err != nil {
		return err
	}
	return nil
}

func CompareFiles(before, after, path string, breakpoint int) {
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
	var f *os.File
	f, err = os.Create(destDir + diffName)
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir("results", os.ModePerm)
			f, _ = os.Create(destDir + diffName)
		} else {
			log.Fatal(err)
		}
	}
	buf := new(bytes.Buffer)
	png.Encode(buf, diff)
	f.Write(buf.Bytes())
	fmt.Println("diff has written into " + destDir + diffName)
}
