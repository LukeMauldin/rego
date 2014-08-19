package main

import (
	"encoding/json"
	"net/http"

	"appengine"
)

// MatchResultResponse is a json type
type MatchResultResponse struct {
	Matches    [][]string `json:"matches"`
	GroupsName []string   `json:"groupsName"`
}

func regExpHandler(rw http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)

	inputData := struct {
		Regexp          string
		Text            string
		FindAllSubmatch bool
	}{}

	err := json.NewDecoder(req.Body).Decode(&inputData)
	if err != nil {
		ctx.Errorf("Error parsing input: %v", err)
		http.Error(rw, "Error parsing input", http.StatusInternalServerError)
		return
	}

	numMatches := -1
	if !inputData.FindAllSubmatch {
		numMatches = 1
	}
	result, err := EvalRegex(ctx, inputData.Regexp, inputData.Text, numMatches)
	if err != nil {
		ctx.Errorf("%v", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(rw).Encode(result)
	if err != nil {
		ctx.Errorf("JSON encoding err: %v", err)
		http.Error(rw, "JSON encoding error", http.StatusInternalServerError)
		return
	}
}

func init() {
	// Regex testing service
	http.HandleFunc("/test_regexp/", regExpHandler)

	// Main handler (index.html)
	http.Handle("/", http.RedirectHandler("/assets/html/index.html", http.StatusMovedPermanently))
}
