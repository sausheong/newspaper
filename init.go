package main

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
)

var papers map[string]Paper

func init() {
	papers = make(map[string]Paper)
	go func() {
    for {
  		checkAndLoad()
  		time.Sleep(15 * time.Minute)      
    }
	}()
}

// represents a newspaper
type Paper struct {
	Name          string
	NumOfPages    int
	Pages         [][]byte
	Previews      [][]byte
	DateRefreshed time.Time
	Sections      []Section
}

type Section struct {
	Name      string
	StartPage int
}

func (paper *Paper) AddSection(s Section) {
	paper.Sections = append(paper.Sections, s)
}

func (paper *Paper) AddPage(pg []byte) {
	paper.Pages = append(paper.Pages, pg)
}

func (paper *Paper) AddPreview(pre []byte) {
	paper.Previews = append(paper.Previews, pre)
}

func checkAndLoad() {
	today := date()
	source := "sources/TODAY_" + today + ".pdf"
  fmt.Println("source is", source)
	url := url(today)
	if checkAndDownload(source, url) {
    fmt.Println("starting conversion ...")
		convert(source, today)
	} else {
	  fmt.Println("no conversion.")
	}
	loadPaper(today)
}

func date() string {
	return time.Now().Format("020106")
}

func url(date string) string {
	return fmt.Sprintf("http://interactivepaper.todayonline.com/jrsrc/%s/%s.pdf", date, date)
}

func loadPaper(date string) {
  fmt.Println("loading paper to memory ...")
	paper := Paper{Name: "today"}
	files, err := ioutil.ReadDir("output/pages")
	if err != nil {
		fmt.Println("cannot read directory", err)
	}
	for n, f := range files {
		if strings.HasPrefix(f.Name(), date) {
			raw, err := ioutil.ReadFile("output/pages/" + f.Name())
			if err != nil {
				fmt.Println("cannot read file", err)
			}
			paper.AddPage(raw)
			fmt.Print(".", n)
		}
	}
  fmt.Println("\npaper loaded.")
	papers[paper.Name] = paper
}

// convert pdf to multiple files
func convert(source string, date string) {
	if sourceExists(source) {
		buildparams := []string{source, "output/pages/" + date}
		cmd := exec.Command("pdftopng", buildparams...)
		fmt.Println("Executing", strings.Join(cmd.Args, " "))
		var out bytes.Buffer
		cmd.Stdout, cmd.Stderr = &out, &out
		err := cmd.Run()
		if err != nil {
			msg := out.String()
			fmt.Println("convert pages:", msg)
		}

		buildparams = []string{"-r", "15", source, "output/previews/" + date}
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
