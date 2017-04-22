package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Args are command line arguments.
type Args struct {
	AddedBy  string
	Title    string
	Filename string
	URL      string
}

func main() {
	args, err := getArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid argument: %s\n", err)
		flag.PrintDefaults()
		os.Exit(1)
	}

	if err := addQuote(args); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to add quote: %s\n", err)
		os.Exit(1)
	}
}

func getArgs() (*Args, error) {
	addedBy := flag.String("added-by", "", "Name of person adding the quote.")
	title := flag.String("title", "", "Title for the quote.")
	filename := flag.String("filename", "", "File containing the quote itself.")
	url := flag.String("url", "", "URL to the quote site API.")

	flag.Parse()

	if len(*addedBy) == 0 {
		return nil, fmt.Errorf("you must specify who is adding the quote")
	}

	if len(*title) == 0 {
		return nil, fmt.Errorf("you must specify a title for the quote")
	}

	if len(*filename) == 0 {
		return nil, fmt.Errorf("you must specify the file containing the quote")
	}

	if len(*url) == 0 {
		return nil, fmt.Errorf("you must specify the URL to the quote site")
	}

	return &Args{
		AddedBy:  *addedBy,
		Title:    *title,
		Filename: *filename,
		URL:      *url,
	}, nil
}

func addQuote(args *Args) error {
	quoteBytes, err := ioutil.ReadFile(args.Filename)
	if err != nil {
		return fmt.Errorf("unable to read quote from file: %s: %s", args.Filename,
			err)
	}

	type Payload struct {
		AddedBy string `json:"added_by"`
		Title   string `json:"title"`
		Quote   string `json:"quote"`
	}

	payload := Payload{
		AddedBy: args.AddedBy,
		Title:   args.Title,
		Quote:   string(quoteBytes),
	}

	json, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("unable to create JSON payload: %s", err)
	}

	url := fmt.Sprintf("%s?version=1&object=quote", args.URL)
	buf := bytes.NewBuffer(json)

	// TODO(horgh): Timeout

	resp, err := http.Post(url, "application/json", buf)
	if err != nil {
		return fmt.Errorf("unable to make HTTP request to %s: %s", url, err)
	}

	responsePayload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		_ = resp.Body.Close()
		return fmt.Errorf("unable to read response body: %s", err)
	}

	if err := resp.Body.Close(); err != nil {
		return fmt.Errorf("problem closing response body: %s", err)
	}

	fmt.Printf("%s", responsePayload)

	return nil
}
