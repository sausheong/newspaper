package main

import "testing"

func TestDownload(t *testing.T) {
  downloadAndSave("http://interactivepaper.todayonline.com/jrsrc/111015/111015.pdf", "test.pdf")
}

func TestConvert(t *testing.T) {
  convert("sources/TODAY_101015.pdf", "101015")
}