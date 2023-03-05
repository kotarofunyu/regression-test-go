package main

import "testing"

func TestMain(t *testing.T) {

	page, driver := setupBrowser()
	defer driver.Stop()
	mytestconf := TestConfig{
		breakpoints: []int{1200},
		baseurl:     "http://localhost:8000/",
		paths:       []string{""},
		initheight:  300,
	}
	rt := RegressionTest{
		mytestconf,
		page,
	}
	rt.Run()
	// URLへアクセス
	// 指定したブレイクポイントで繰り返す
	// パスごとにbefore/afterのスクショを撮って比較する
	// 違いがあれば出力する

	//
}
