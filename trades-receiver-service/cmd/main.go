package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"insidertradesreceiver/internal/app"
	"insidertradesreceiver/internal/config"
	"insidertradesreceiver/internal/logger"
	"io"
	"log"
	"net/http"
	"os"
)

func Respond(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)

		if err != nil {
			Error(w, r, http.StatusInternalServerError, err)
		}
	}
}

func Error(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	Respond(w, r, statusCode, map[string]string{"error": err.Error()})
}


//type trade struct {
//	Name string `json:"name"`
//}

func HandlePostTrades(w http.ResponseWriter, r *http.Request) {
	//var tr trade
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("can't read request")
	}
	//decoder := json.NewDecoder(r.Body)
	//if err := decoder.Decode(&tr); err != nil {
	//	log.Println("Problem with decoding...")
	//	Error(w, r, http.StatusBadRequest, errors.New("Invalid JSON"))
	//}
	log.Println("ok")
	log.Println(string(data))
	logFile, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Println("Can't open a file")
	}
	defer logFile.Close()
	writer := bufio.NewWriter(logFile)
	n, err := writer.WriteString(string(data))
	if err != nil {
		log.Println("Can't write in file")
	}
	log.Println("Bytes wrote:", n)
	defer writer.Flush()
	Respond(w, r, http.StatusCreated, nil)
}

func HandleHomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Under Construction")
}

func HandleGetTrades(w http.ResponseWriter, r *http.Request) {

}

func buildHandler() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/insider-trades/receiver", HandlePostTrades).Methods("POST")
	router.HandleFunc("/trades", HandleGetTrades).Methods("GET")
	router.HandleFunc("/", HandleHomePage).Methods("GET")

	return router
}


func main() {
	l := logger.New()
	l.Info("Logger initialized")

	cfg, err := config.New()
	if err != nil {
		l.Fatal("Failed to load server config: %v", err)
	}

	app.Run(cfg, l)
	//handler := buildHandler()
	//http.ListenAndServe(":8080", handler)
}