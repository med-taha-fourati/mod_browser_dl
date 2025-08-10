package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"encoding/json"
)

const serverPort = 3000

/* 
Usage:
- GET /browse?search=keyword returns the list of songs
- POST /download/{id} downloads the thing
*/

func ModFileDetails(w http.ResponseWriter, r *http.Request) {
	//TODO: fetch the metadata of the song
	
}

func BrowseFunction(w http.ResponseWriter, r *http.Request) {
	search_query := r.URL.Query()
	//fmt.Fprintf(w, "Query params: %v\n", search_query)

	word := search_query.Get("search")
	//fmt.Fprintf(w, "Search: %s %s\n", word, index)

	var codes []string
  	var names []string
		
	codes, names = automaticDepaginator(word)

	if (len(codes) < 1) {
		fmt.Fprintf(w, "Sorry, no results were found\n")
		http.Error(w, "Bad request\n", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	jsonData := []map[string]interface{}{}

	for i, code := range codes {
		jsonData = append(jsonData, map[string]interface{}{
			"name": names[i],
			"code": code,
		})
		//fmt.Fprintf(w, "%s | %s\n", code, names[i])
	}

	jsonHeader := map[string]interface{}{
		"result": jsonData,
	}

	jsonString, err := json.Marshal(jsonHeader)

	if err != nil {
		fmt.Fprintf(w, "Unknown error occured\n")
		http.Error(w, "Internal Server Error\n", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", string(jsonString))

}

func DownloadFunction(w http.ResponseWriter, r *http.Request) {
	if (r.Method != http.MethodPost) {
		fmt.Fprintf(w, "POST method allowed only\n")
		http.Error(w, "Method Not allowed\n", http.StatusMethodNotAllowed)
		return
	}

	search_query := r.URL.Query()
	word := search_query.Get("search")

	id, err := strconv.Atoi(word)
	fmt.Fprintf(w, "Gotten id: %d\n", id)

	if err != nil {
		fmt.Fprintf(w, "Invalid ID\n")
		http.Error(w, "Bad request\n", http.StatusBadRequest)
		return
	}

	x := word
	result_name, url_name := lookupFileName(x);
	//fmt.Println(result_name, url_name)
	fmt.Fprintf(w, "Downloading in progress... Check your server console\n")

	downloadFile(result_name, url_name)
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
