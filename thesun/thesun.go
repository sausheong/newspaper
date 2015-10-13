package thesun

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
   "github.com/sausheong/newspaper/paper"
)


func CheckAndLoad() paper.Paper {
	today := date()
  day := day()
	source := "sources/THESUN_" + today + ".pdf"
  fmt.Println("source is", source)
	url := url(today, day)
	if checkAndDownload(source, url) {
    fmt.Println("starting conversion ...")
		convert(source, today)
	} else {
	  fmt.Println("no conversion.")
	}
	return loadPaper(today)
}

func date() string {
	return time.Now().Format("02012006")
}

func day() string {
  return strings.ToLower(time.Now().Format("Mon"))
}

func url(date string, day string) string {  
	return fmt.Sprintf("http://thesun-epaper.com/%s/%s/files/assets/common/downloads/publication.pdf", day, date)
}

func loadPaper(date string) (p paper.Paper) {
  fmt.Println("loading paper to memory ...")
	p = paper.Paper{Name: "thesun"}
	files, err := ioutil.ReadDir("output/thesun/pages")
	if err != nil {
		fmt.Println("cannot read directory", err)
	}
	for n, f := range files {
		if strings.HasPrefix(f.Name(), date) {
			raw, err := ioutil.ReadFile("output/thesun/pages/" + f.Name())
			if err != nil {
				fmt.Println("cannot read file", err)
			}
			p.AddPage(raw)
			fmt.Print(".", n)
		}
	}
  fmt.Println("\npaper loaded.")
	return
}

// convert pdf to multiple files
func convert(source string, date string) {
	if sourceExists(source) {
		buildparams := []string{source, "output/thesun/pages/" + date}
		cmd := exec.Command("pdftopng", buildparams...)
		fmt.Println("Executing", strings.Join(cmd.Args, " "))
		var out bytes.Buffer
		cmd.Stdout, cmd.Stderr = &out, &out
		err := cmd.Run()
		if err != nil {
			msg := out.String()
			fmt.Println("convert pages:", msg)
		}

		buildparams = []string{"-r", "15", source, "output/thesun/previews/" + date}
		cmd = exec.Command("pdftopng", buildparams...)
		fmt.Println("Executing", strings.Join(cmd.Args, " "))
		cmd.Stdout, cmd.Stderr = &out, &out
		err = cmd.Run()
		if err != nil {
			msg := out.String()
			fmt.Println("convert previews:", msg)
		}
	} else {
		fmt.Println("Source not found, cannot convert")
		return
	}
}

// check if source file exists and download it if it's not
func checkAndDownload(source string, url string) bool {
	if !sourceExists(source) {
		fmt.Println("Source", source, "not found, trying to download now.")
		return downloadAndSave(url, source)
	} else {
		fmt.Println("Already downloaded source file.")
		return false
	}
}

// check if the source file exists
func sourceExists(source string) bool {
	_, err := os.Stat(source)
	if err != nil {
		fmt.Println("source does not exist")
		return false
	} else {
		fmt.Println("source exist")
		return true
	}
}

// download source file and save it
// returns true if it's downloaded properly
func downloadAndSave(url string, filename string) bool {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode == 404 {
		fmt.Println("Cannot get source at", url, err, resp.Status)
		return false
	}
	defer resp.Body.Close()

	source, err := os.Create(filename)
	if err != nil {
		fmt.Println("Cannot create file:", filename, err)
		return false
	}
	defer source.Close()
	_, err = io.Copy(source, resp.Body)
	if err != nil {
		fmt.Println("Cannot copy to file:", err)
		return false
	}
	return true
}
