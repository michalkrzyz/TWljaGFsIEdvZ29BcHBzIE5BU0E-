package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"url-collector/collector"
	"url-collector/collector/handlers/nasaapi"
	"url-collector/daterange"
	"url-collector/interfaces"
)

const (
	defaultUrlCollectorPort = "8090"
)

type ErrorResponse struct {
	Msg string `json:"error"`
}

type UrlListResponse struct {
	Urls []string `json:"urls"`
}

func pictures(w http.ResponseWriter, r *http.Request, urlCollector interfaces.UrlCollector) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()
	encoder := json.NewEncoder(w)

	from := q.Get("from")
	if from == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		encoder.Encode(ErrorResponse{Msg: "Missing 'from' URL Query parameter"})
		return
	}

	to := q.Get("to")
	if to == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		encoder.Encode(ErrorResponse{Msg: "Missing 'to' URL Query parameter"})
		return
	}

	fromTime, toTime, convertionErr := daterange.ConvertRangeParametersToTime(from, to)
	if convertionErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(ErrorResponse{Msg: fmt.Sprintf("Invalid range parameters: %s", convertionErr)})
		return
	}

	urlList, retriveErr := urlCollector.GetUrlList(fromTime, toTime)
	if retriveErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(ErrorResponse{Msg: fmt.Sprintf("Problem retreiving data: %s", retriveErr)})
		return
	}
	encoder.Encode(UrlListResponse{Urls: urlList})
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func main() {
	port := os.Getenv("URL_COLLECTOR_PORT")
	if port == "" {
		port = defaultUrlCollectorPort
	}

	urlCollectorHandler := nasaapi.NasaApiUrlCollectorHandlerFactory{}.New()
	urlCollector := collector.BaseUrlCollectorFactory{}.New(urlCollectorHandler)

	fmt.Printf("Starting url-collector using port: %s\n", port)

	http.HandleFunc("/pictures", func(w http.ResponseWriter, r *http.Request) { pictures(w, r, urlCollector) })
	http.HandleFunc("/health", healthcheck)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	fmt.Println(err)
}
