package kor_name_match

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/suapapa/go_hangul"
)

func init() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/match/", matchHandler)
}

var templates = template.Must(template.ParseFiles("template/index.html"))

func rootHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

type MatchResult struct {
	N1      []rune
	N2      []rune
	Prog    string
	Percent int
}

func newMatchResult(name1, name2 string) (*MatchResult, error) {
	//TODO: N1 and N2 should be string
	r := &MatchResult{
		N1: []rune(name1),
		N2: []rune(name2),
	}

	if len(r.N1) != 3 || len(r.N2) != 3 {
		return nil, errors.New("Name should be three characters")
	}

	return r, nil
}

func matchHandler(w http.ResponseWriter, r *http.Request) {
	nr, err := newMatchResult(r.FormValue("name1"), r.FormValue("name2"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nameMatch(nr)

	// TODO: executeTemplate
	fmt.Fprintln(w, nr)
}

func nameMatch(nr *MatchResult) {
	r1 := nr.N1
	r2 := nr.N2
	rc := []rune{r1[0], r2[0], r1[1], r2[1], r1[2], r2[2]}
	nr.Prog += fmt.Sprintln(string(rc))

	rn := make([]int, 6)
	for i, r := range rc {
		rn[i] = hangul.Stroke(r)
	}

	nr.Percent = match(rn, nr)
}

func match(in []int, nr *MatchResult) int {
	nr.Prog += fmt.Sprintln(in)
	r := make([]int, len(in)-1)
	for i := 0; i < len(r); i++ {
		r[i] = (in[i] + in[i+1]) % 10
	}

	if len(in) > 2 {
		match(r, nr)
	}

	// TODO: bugfix
	return in[0]*10 + in[1]
}
