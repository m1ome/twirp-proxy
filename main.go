package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	url  string
	addr string
)

func init() {
	flag.StringVar(&url, "url", "http://127.0.0.1:8080", "provide your twirp server base url")
	flag.StringVar(&addr, "addr", ":8090", "proxy binding address")
	flag.Parse()
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Allow", "POST, OPTIONS")

	if r.Method == http.MethodPost {
		req, err := http.NewRequest("POST", url+r.URL.Path, r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatalf("error creating request: %v", err)
			return
		}
		req.Header.Add("Content-Type", "application/json")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatalf("error sending request to twirp service: %v", err)
			return
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatalf("error reading response body: %v", err)
			return
		}

		log.Printf("forwarding %s, response [%d]: %s", r.URL.Path, res.StatusCode, string(body))
		w.WriteHeader(res.StatusCode)
		w.Write(body)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	log.Printf("twirp endpoint: %s", url)

	log.Printf("start listening incoming requests on: %s", addr)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(addr, nil))
}
