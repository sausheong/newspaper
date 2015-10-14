package paper

import "time"

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

type Page struct {
  Page []byte `json:"page"`
  Num int `json:"num"`  
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