package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/ray31245/tinder/pkg/match"
)

type Person struct {
	Name   string
	Height uint
	Gender string
}

func welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "welcome to my tinder")
}

var MatchSystem match.MatchSystem = match.NewMatchSystem()

func AddSinglePersonAndMatch(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "require a name")
		return
	}
	height := parseIntParm(r, "height")
	if height <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "height must greater than zero")
		return
	}
	gender, ok := match.GenderToValueMap[r.URL.Query().Get("gender")]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "system does not has this gender")
		return
	}
	wantedDate := parseIntParm(r, "wantedDate")
	if wantedDate <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "wantedDate must greater than zero")
		return
	}

	p := MatchSystem.AddSinglePerson(name, uint(height), gender, wantedDate)
	matchPersons, _ := MatchSystem.MatchPerson(p.ID, p.WantedDate)
	MatchSystem.RecordNewMatch(&p, matchPersons)
	matchPersons = match.Display(&p, p.WantedDate)

	res := struct {
		ID    string
		Match []Person
	}{
		ID:    p.ID.String(),
		Match: []Person{},
	}
	for _, p := range matchPersons {
		res.Match = append(res.Match, Person{Name: p.Name, Height: p.Height, Gender: match.GenderToStringMap[p.Gender]})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func RemoveSinglePerson(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "require a id")
		return
	}
	MatchSystem.RemovePerson(id)
}

func QuerySinglePeople(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "require a id")
		return
	}
	N := parseIntParm(r, "N")
	if N <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "N must greater than zero")
		return
	}
	p, ok := MatchSystem.PersonList[id]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Person not found")
		return
	}
	remaining := p.WantedDate - len(p.MatchPerson)
	if N > remaining {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "remaining %d, but want %d", remaining, N)
		return
	}

	res := struct {
		ID    string
		Match []Person
	}{
		ID:    p.ID.String(),
		Match: []Person{},
	}

	if N <= len(p.NewMatchPerson) {
		for _, v := range match.Display(p, N) {
			res.Match = append(res.Match, Person{Name: v.Name, Height: v.Height, Gender: match.GenderToStringMap[v.Gender]})
		}
	} else {
		matchPersons, _ := MatchSystem.MatchPerson(p.ID, N-len(p.NewMatchPerson))
		MatchSystem.RecordNewMatch(p, matchPersons)
		for _, v := range match.Display(p, N) {
			res.Match = append(res.Match, Person{Name: v.Name, Height: v.Height, Gender: match.GenderToStringMap[v.Gender]})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func parseIntParm(r *http.Request, param string) int {
	val := r.URL.Query().Get(param)
	if val == "" {
		return 0
	}
	num, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return num
}

func main() {
	http.HandleFunc("/", welcome)
	http.HandleFunc("/add", AddSinglePersonAndMatch)
	http.HandleFunc("/remove", RemoveSinglePerson)
	http.HandleFunc("/query", QuerySinglePeople)

	port := "8089"

	log.Println("Tinder match server is running on port " + port)
	http.ListenAndServe(":"+port, nil)
}
