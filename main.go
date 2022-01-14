package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
)

const (
	nameURL = "https://names.mcquay.me/api/v0"
	jokeURL = "http://api.icndb.com/jokes/random"
)

type Options struct {
	Port string `short:"p" long:"port" description:"port on which app will be listening" default:"8080"`
}

// Name holds response structure from random names API
type Name struct {
	First string `json:"first_name"`
	Last  string `json:"last_name"`
}

// Joke holds response structure from random jokes API
type Joke struct {
	Type  string `json:"type"`
	Value struct {
		ID         int      `json:"id"`
		Joke       string   `json:"joke"`
		Categories []string `json:"categories"`
	} `json:"value"`
}

func main() {
	opts := &Options{}

	p := flags.NewParser(opts, flags.Default)
	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	fmt.Printf("listening on localhost:%s\n", opts.Port)
	log.Fatal(http.ListenAndServe("localhost:"+opts.Port, Router()))
}

func Router() *chi.Mux {
	mux := chi.NewRouter()
	mux.Get("/", JokeHandler)
	mux.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})

	return mux
}

func JokeHandler(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	resp, err := client.Get(nameURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		sendErr(w, err, "error sending request to get random name")
		return
	}
	name := &Name{}
	err = json.NewDecoder(resp.Body).Decode(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		sendErr(w, err, "error decoding random name response")
		return
	}

	u, err := url.Parse(jokeURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		sendErr(w, err, "error parsing joke url")
		return
	}
	val := u.Query()
	val.Add("firstName", name.First)
	val.Add("lastName", name.Last)
	u.RawQuery = val.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		sendErr(w, err, "error creating request to get random joke")
		return
	}

	resp, err = client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		sendErr(w, err, "error making joke request")
		return
	}

	joke := &Joke{}
	err = json.NewDecoder(resp.Body).Decode(joke)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		sendErr(w, err, "error parsing joke response")
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprint(joke.Value.Joke)))
}

func sendErr(w http.ResponseWriter, err error, msg string) {
	if err == nil {
		log.Printf("[ERROR] %s", errors.New(msg))
	} else {
		log.Printf("[ERROR] %s", errors.Wrap(err, msg))
	}

	_, _ = w.Write([]byte(msg))
}
