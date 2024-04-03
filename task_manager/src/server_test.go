package test

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	_, err := io.WriteString(w, "This is my website!\n")
	if err != nil {
		return
	}
}

func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	_, err := io.WriteString(w, "Hello, HTTPs!\n")
	if err != nil {
		return
	}
}

func getESElementById(w http.ResponseWriter, r *http.Request, es string, id string) {
	resp, err := http.Get(es + "crawl_data/_search?q=_id:" + id)
	if err != nil {
		fmt.Printf("http.Get() failed with %s\n", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return
	}
}

func main() {
	elasticsearchHost := os.Getenv("ELASTICSEARCH_HOST")
	if elasticsearchHost == "" {
		elasticsearchHost = "localhost"
	}

	http.HandleFunc("/", getRoot)
	http.HandleFunc("/hello", getHello)
	url := fmt.Sprintf("https://%s:9200/", elasticsearchHost)
	http.HandleFunc("/es/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("got /es/ request\n")
		id := r.URL.Path[len("/es/"):]
		getESElementById(w, r, url, id)
	})
	fmt.Printf("Starting server at port 3333\n")
	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		fmt.Printf("http.ListenAndServe() failed with %s\n", err)
	}
}
