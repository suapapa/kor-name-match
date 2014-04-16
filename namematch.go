package main

import (
	"fmt"
	"net/http"

	"github.com/suapapa/go_hangul"
)

const namesForm = `
<html>
	<body>
		<h1>이름궁합</h1>
		<form action="/match/" method="GET">
			<div>내 이름 <input type="text" name="name1"></div>
			<div>상대방 이름 <input type="text" name="name2"></div>
			<div><input type="submit" value="Match"></div>
		</form>
	</body>
</html>
`

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, namesForm)
}

func matchHandler(w http.ResponseWriter, r *http.Request) {
	nameMatch(w, r.FormValue("name1"), r.FormValue("name2"))
}

func nameMatch(w http.ResponseWriter, name1, name2 string) {
	fmt.Fprintf(w, "<html>%s %s<br>", name1, name2)
	r1 := []rune(name1)
	r2 := []rune(name2)
	if len(r1) != 3 || len(r2) != 3 {
		http.Error(w, "석 자 이름을 넣으세요", http.StatusBadRequest)
		return
	}

	rc := []rune{r1[0], r2[0], r1[1], r2[1], r1[2], r2[2]}
	rn := make([]int, 6)
	for i, r := range rc {
		rn[i] = hangul.Stroke(r)
	}
	match(w, rn)
	fmt.Fprintf(w, "</html>")
}

func match(w http.ResponseWriter, in []int) {
	fmt.Fprintf(w, "%v <br>", in)
	r := make([]int, len(in)-1)
	for i := 0; i < len(r); i++ {
		r[i] = (in[i] + in[i+1]) % 10
	}

	if len(in) > 2 {
		match(w, r)
	}
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/match/", matchHandler)

	http.ListenAndServe(":8080", nil)
}
