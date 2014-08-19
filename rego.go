package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"appengine"
	"appengine/datastore"
)

// MatchResultResponse is a json type
type MatchResultResponse struct {
	Matches    [][]string `json:"matches"`
	GroupsName []string   `json:"groupsName"`
}

type RegexStats struct {
	RegexLength      int
	TextStringLength int
	MatchDuration    time.Duration
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

	startRegex := time.Now()

	r, err := regexp.Compile(inputData.Regexp)
	if err != nil {
		ctx.Errorf("Invalid RegExp : %s \n", inputData.Regexp)
		http.Error(rw, fmt.Sprintf("Invalid RegExp : %s", inputData.Regexp), http.StatusInternalServerError)
		return
	}

	numMatches := -1
	if !inputData.FindAllSubmatch {
		numMatches = 1
	}
	matches := r.FindAllStringSubmatch(inputData.Text, numMatches)

	m := &MatchResultResponse{}
	if len(matches) > 0 {
		m.Matches = matches
		m.GroupsName = r.SubexpNames()[1:]
	}

	regexDuration := time.Since(startRegex)
	ctx.Infof("regex duration: %v", regexDuration)
	stats := &RegexStats{
		RegexLength:      len(inputData.Regexp),
		TextStringLength: len(inputData.Text),
		MatchDuration:    regexDuration,
	}
	writeStatsEntry(ctx, stats)

	enc := json.NewEncoder(rw)
	err = enc.Encode(m)
	if err != nil {
		ctx.Errorf("JSON encoding err: %v", err)
		http.Error(rw, "JSON encoding error", http.StatusInternalServerError)
		return
	}
}

func writeStatsEntry(ctx appengine.Context, stats *RegexStats) {
	_, err := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "regex_stats", nil), stats)
	if err != nil {
		ctx.Errorf("Error writing stats: %v", err)
	}
}

func init() {
	// Regex testing service
	http.HandleFunc("/test_regexp/", regExpHandler)

	// Main handler (index.html)
	http.Handle("/", http.RedirectHandler("/assets/html/index.html", http.StatusMovedPermanently))
}
