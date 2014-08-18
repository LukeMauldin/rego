package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
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
	TestStringLength int
	MatchDuration    time.Duration
}

func regExpHandler(rw http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)

	var matches [][]string

	err := req.ParseForm()
	if err != nil {
		ctx.Errorf("Error parsing form: %v", err)
		http.Error(rw, "Error parsing form", http.StatusInternalServerError)
		return
	}
	regexpString := req.FormValue("regexp")
	testString := req.FormValue("testString")
	findAllSubmatch, _ := strconv.ParseBool(req.FormValue("findAllSubmatch"))

	//log.Printf("Regexp : %s", regexpString)
	//log.Printf("Test string : %s", testString)
	//log.Printf("Find all : %t", findAllSubmatch)

	startRegex := time.Now()

	m := &MatchResultResponse{}

	r, err := regexp.Compile(regexpString)
	if err != nil {
		ctx.Errorf("Invalid RegExp : %s \n", regexpString)
		http.Error(rw, fmt.Sprintf("Invalid RegExp : %s", regexpString), http.StatusInternalServerError)
		return
	}

	if findAllSubmatch {
		matches = r.FindAllStringSubmatch(testString, -1)
	} else {
		matches = [][]string{r.FindStringSubmatch(testString)}
	}

	//log.Println(matches)

	if len(matches) > 0 {
		m.Matches = matches
		m.GroupsName = r.SubexpNames()[1:]
	}

	regexDuration := time.Since(startRegex)
	ctx.Infof("regex duration: %v", regexDuration)
	stats := &RegexStats{
		RegexLength:      len(regexpString),
		TestStringLength: len(testString),
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
	// Main handler (index.html)
	http.Handle("/", http.RedirectHandler("/assets/html/index.html", http.StatusMovedPermanently))

	// Regex testing service
	http.HandleFunc("/test_regexp/", regExpHandler)
}
