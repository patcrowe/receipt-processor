package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var receipts map[uuid.UUID]int

type item struct {
	shortDescription string `json:"ShortDescription"`
	price            string `json:"price"`
}

type receiptStruct struct {
	retailer     string `json:"retailer"`
	purchaseDate string `json:"purchaseDate"`
	purchaseTime string `json:"purchaseDate"`
	items        []item `json:"items"`
	price        string `json:"price"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/receipts/process", process).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", getPoints).Methods("GET")
}

func (*receiptStruct) calculatePoints() int {

	return 0
}

func (*receiptStruct) validateReceipt() bool {
	return true
}

func process(response http.ResponseWriter, request *http.Request) {
	sendInvalidResponse := func() {
		response.WriteHeader(400)
		response.Write([]byte("The receipt is invalid."))
	}

	rBody, err := io.ReadAll(request.Body)
	if err != nil {
		sendInvalidResponse()
		return
	}

	newReceipt := receiptStruct{}

	json.Unmarshal(rBody, &newReceipt)
	if !newReceipt.validateReceipt() {
		sendInvalidResponse()
		return
	}

	receiptUUID := uuid.NewSHA1(uuid.Max, rBody)

	fmt.Println(receiptUUID)

}

func getPoints(response http.ResponseWriter, request *http.Request) {

}
