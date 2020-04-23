package tdxcrawl

import (
	"fmt"

	"testing"
	"time"
)

func TestWindow(t *testing.T) {
	n := time.Now()
	ZcTdxWindow()
	// ZcSaveTic("159998", "D:\\TdxData", "20200412", "20200424")
	ZcSaveTic("159909", "D:\\TdxData", "20200101", "20200424")

	fmt.Println(time.Since(n))
}
func TestCrawl(t *testing.T) {
	// n := time.Now()
	// ZcTdxWindow()
	// cDtickWindow("159998")
	// cHrightEnd("20200423")
	// fmt.Println(time.Since(n))
}
