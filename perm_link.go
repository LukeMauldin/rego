package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/gob"
	"fmt"

	"appengine"
	"appengine/datastore"
)

const permLinkLength = 16

type PermLinkStore struct {
	Input MatchInput `datastore:",noindex"`
	Hash  string
}

func CreatePermLink(ctx appengine.Context, input *MatchInput) (string, error) {
	//Create SHA-1 hash
	var bytes bytes.Buffer
	enc := gob.NewEncoder(&bytes)
	err := enc.Encode(input)
	if err != nil {
		return "", err
	}
	hashOutput := sha1.Sum(bytes.Bytes())
	strHash := fmt.Sprintf("%x", hashOutput[:])
	strHash = strHash[:permLinkLength]

	//Only store if hash doesn't already exist
	num, err := datastore.NewQuery("perm_link").Filter("Hash =", strHash).Count(ctx)
	if err != nil {
		return "", err
	}
	if num > 0 {
		return strHash, nil
	}

	//If hash does not exist, store in datastore
	entity := &PermLinkStore{*input, strHash}
	_, err = datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "perm_link", nil), entity)
	if err != nil {
		return "", err
	}

	return strHash, nil
}

func RetrievePermLink(ctx appengine.Context, permLink string) (*MatchInput, error) {
	if len(permLink) != permLinkLength {
		return nil, fmt.Errorf("Invalid share link")
	}

	results := make([]*PermLinkStore, 0)
	_, err := datastore.NewQuery("perm_link").Filter("Hash =", permLink).GetAll(ctx, &results)
	if err != nil {
		return nil, err
	}
	if len(results) != 1 {
		return nil, fmt.Errorf("No records found")
	}

	return &results[0].Input, nil
}
