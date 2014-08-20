package main

import (
	"fmt"
	"regexp"
	"runtime"
	"time"

	"appengine"
	"appengine/datastore"
)

func EvalRegex(ctx appengine.Context, expr string, text string, numMatches int) (*MatchResultResponse, error) {
	startRegex := time.Now()

	//Compile regex
	r, err := regexp.Compile(expr)
	if err != nil {
		return nil, fmt.Errorf("Invalid RegExp: %s -- %v", expr, err)
	}

	//Evaluate regex and get results
	matches := r.FindAllStringSubmatch(text, numMatches)
	result := &MatchResultResponse{}
	if len(matches) > 0 {
		result.Matches = matches
		result.GroupsName = r.SubexpNames()[1:]
	}

	//Record stats
	go recordStats(ctx, expr, text, time.Since(startRegex))

	return result, nil
}

type regexStats struct {
	ExprLength int           `datastore:",noindex"`
	TextLength int           `datastore:",noindex"`
	RanAt      time.Time     `datastore:",noindex"`
	Duration   time.Duration `datastore:",noindex"`
}

func recordStats(ctx appengine.Context, expr string, text string, duration time.Duration) {
	defer func() {
		if rec := recover(); rec != nil {
			stack := make([]byte, 32768)
			runtime.Stack(stack, false)
			ctx.Errorf("Recovered from panic: %v -- %s", rec, stack)
		}
	}()
	//Store stats in datastore for reporting
	stats := &regexStats{
		ExprLength: len(expr),
		TextLength: len(text),
		RanAt:      time.Now(),
		Duration:   duration,
	}
	_, err := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "regex_stats", nil), stats)
	if err != nil {
		ctx.Errorf("Error writing stats: %v", err)
	}

}
