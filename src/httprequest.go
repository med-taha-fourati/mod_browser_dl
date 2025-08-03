package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const serverPort = 3000

/* 
Usage:
- GET /browse?search=keyword returns the list of songs
- POST /download/{id} downloads the thing
*/

func BrowseFunction(w http.ResponseWriter, r *http.Request) {
	search_query := r.URL.Query()
	fmt.Fprintf(w, "Query params: %v\n", search_query)

	word := search_query.Get("search")
	index := search_query.Get("index")
	fmt.Fprintf(w, "Search: %s %s\n", word, index)
	  
	var codes []string
  	var names []string
	index_int, err := strconv.Atoi(index)
	if err != nil {
		fmt.Fprintf(w, "Invalid Index, Defaulting to 1\n");
		http.Error(w, "Bad request\n", http.StatusBadRequest)
	}
	codes, names = retrieve(word, index_int)

	if (len(codes) < 1) {
		fmt.Fprintf(w, "Sorry, no results were found\n")
		http.Error(w, "Bad request\n", http.StatusBadRequest)
		return
	}

	for i, code := range codes {
		fmt.Fprintf(w, "%s | %s\n", code, names[i])
	}
}

func DownloadFunction(w http.ResponseWriter, r *http.Request) {
	if (r.Method != http.MethodPost) {
		fmt.Fprintf(w, "POST method allowed only\n")
		http.Error(w, "Method Not allowed\n", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/");
	if (len(parts) != 3 || parts[2] == "") {
		fmt.Fprintf(w, "Nothing has been given\n")
		http.Error(w, "Bad Request\n", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[2])
	fmt.Fprintf(w, "Gotten id: %d\n", id)

	if err != nil {
		fmt.Fprintf(w, "Invalid ID\n")
		http.Error(w, "Bad request\n", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Downloading in progress... Check your server console\n")
	
}

func main() {
	// go func() {
		
	// }()

	mux := http.NewServeMux()
		
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "Hi :wave: %s\n", r.URL.Path)
		})
		
		mux.HandleFunc("/browse", BrowseFunction)
		mux.HandleFunc("/download", DownloadFunction)

		server := http.Server {
				Addr: fmt.Sprintf(":%d", serverPort),
			Handler: mux,
		}

		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				fmt.Printf("error running http server: %s\n", err)
	}
	

	// requestURL := fmt.Sprintf("http://localhost:%d", serverPort)
	// res, err := http.Get(requestURL)
	// if err != nil {
	// 	fmt.Printf("error making http request: %s\n", err)
	// 	os.Exit(1)
	// }
		
	// fmt.Printf("client: got response!\n")
	// fmt.Printf("client: status code: %d\n", res.StatusCode)
	
}}
