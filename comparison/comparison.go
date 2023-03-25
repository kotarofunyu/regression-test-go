package comparison

import (
	"fmt"
	"log"
	"os"

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
