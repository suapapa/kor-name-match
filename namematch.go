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

var templates = template.Must(template.ParseFiles("template/index.html", "template/result.html"))

func rootHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

type MatchResult struct {
	N1, N2  string
	Prog    string
	Percent int
}

func newMatchResult(name1, name2 string) (*MatchResult, error) {
	n1, n2 := []rune(name1), []rune(name2)

	if len(n1) != 3 || len(n2) != 3 {
		return nil, errors.New("Name should be three characters")
	}

	if !isHangulRunes(n1) || !isHangulRunes(n2) {
		return nil, errors.New("Only hangul name is supported")
	}

	return &MatchResult{
		N1: name1,
		N2: name2,
	}, nil
}

func isHangulRunes(s []rune) bool {
	for _, c := range s {
		if !hangul.IsHangul(c) {
			return false
		}
	}
	return true
}

func matchHandler(w http.ResponseWriter, r *http.Request) {
	nr, err := newMatchResult(r.FormValue("name1"), r.FormValue("name2"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nameMatch(nr)

	templates.ExecuteTemplate(w, "result.html", nr)
}

func nameMatch(nr *MatchResult) {
	r1, r2 := []rune(nr.N1), []rune(nr.N2)
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
		return match(r, nr)
	}

	return in[0]*10 + in[1]
}
