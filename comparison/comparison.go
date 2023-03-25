package comparison

import (
	"log"

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
