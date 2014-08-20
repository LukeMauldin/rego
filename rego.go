package main

import (
	"encoding/json"
	"net/http"

	"appengine"
)

type MatchInput struct {
	Expr       string `datastore:",noindex"`
	Text       string `datastore:",noindex"`
	NumMatches int    `datastore:",noindex"`
}

// MatchResultResponse is a json type
type MatchResultResponse struct {
	Matches    [][]string `json:"matches"`
	GroupsName []string   `json:"groupsName"`
}

func handlerEvalRegexp(rw http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)

	inputData := &MatchInput{}
	err := json.NewDecoder(req.Body).Decode(inputData)
	if err != nil {
		ctx.Errorf("Error parsing input: %v", err)
		http.Error(rw, "Error parsing input", http.StatusInternalServerError)
		return
	}

	result, err := EvalRegex(ctx, inputData)
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

func handlerShareRegexp(rw http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)

	inputData := &MatchInput{}
	err := json.NewDecoder(req.Body).Decode(inputData)
	if err != nil {
		ctx.Errorf("Error parsing input: %v", err)
		http.Error(rw, "Error parsing input", http.StatusInternalServerError)
		return
	}

	result, err := CreatePermLink(ctx, inputData)
	if err != nil {
		ctx.Errorf("%v", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = rw.Write([]byte(result))
	if err != nil {
		ctx.Errorf("%v", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handlerLoadRegexp(rw http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)

	key := req.URL.Query().Get("key")
	if key == "" {
		http.RedirectHandler("/assets/html/index.html", http.StatusMovedPermanently)
		return
	}

	matchInput, err := RetrievePermLink(ctx, key)
	if err != nil {
		ctx.Errorf("%v", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	//TODO: figure out how to get loadRegexp results to evaluate and display to user (need to replace defaults of $scope.regexpInput)
	//TODO: OR have JS type the presence of ?key= and replace the defaults on page load
	err = json.NewEncoder(rw).Encode(matchInput)
	if err != nil {
		ctx.Errorf("JSON encoding err: %v", err)
		http.Error(rw, "JSON encoding error", http.StatusInternalServerError)
		return
	}
}

func init() {
	// Regex eval service
	http.HandleFunc("/eval_regexp/", handlerEvalRegexp)

	//Regex share service
	http.HandleFunc("/share_regexp/", handlerShareRegexp)

	//Regex load share service
	http.HandleFunc("/load_regexp", handlerLoadRegexp)

	// Main handler (index.html)
	http.Handle("/", http.RedirectHandler("/assets/html/index.html", http.StatusMovedPermanently))
}
