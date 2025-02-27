package main

import (
	"encoding/json"
	"io"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var receiptPoints = make(map[uuid.UUID]int)

type item struct {
	ShortDescription string `json:"ShortDescription"`
	Price            string `json:"price"`
}

type receiptStruct struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []item `json:"items"`
	Total        string `json:"total"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/receipts/process", process).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", getPoints).Methods("GET")
	http.ListenAndServe(":8000", r)
}

func (rS *receiptStruct) calculatePoints() int {
	returnPoints := 0

	// One point for every alphanumeric character in the retailer name.
	returnPoints += len(regexp.MustCompile(`[a-zA-Z0-9]`).FindAllString(rS.Retailer, -1))

	// 50 points if the total is a round dollar amount with no cents.
	totalFloat, _ := strconv.ParseFloat(rS.Total, 64)
	if (int(totalFloat*100) % 100) == 0 {
		returnPoints += 50
	}

	// 25 points if the total is a multiple of 0.25.
	if (int(totalFloat*100) % 25) == 0 {
		returnPoints += 25
	}

	// 5 points for every two items on the receipt.
	returnPoints += 5 * (len(rS.Items) / 2)

	// If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
	for _, item := range rS.Items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			priceFloat, _ := strconv.ParseFloat(item.Price, 64)
			returnPoints += int(math.Ceil(priceFloat * 0.2))
		}
	}

	// 6 points if the day in the purchase date is odd.
	if int(rS.PurchaseDate[len(rS.PurchaseDate)-1])%2 == 1 {
		returnPoints += 6
	}

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	hour, _ := strconv.ParseInt(rS.PurchaseTime[0:2], 10, 64)
	if hour >= 14 && hour < 16 {
		returnPoints += 10
	}

	return returnPoints
}

func (rS *receiptStruct) validateReceipt() bool {
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

	_, found := receiptPoints[receiptUUID]
	if !found {
		pointsAwarded := newReceipt.calculatePoints()
		receiptPoints[receiptUUID] = pointsAwarded
	}

	responseStruct := struct {
		ID string `json:"id"`
	}{
		ID: receiptUUID.String(),
	}

	responseJson, err := json.Marshal(responseStruct)
	if err != nil {
		sendInvalidResponse()
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.Write(responseJson)
}

func getPoints(response http.ResponseWriter, request *http.Request) {
	sendInvalidResponse := func() {
		response.WriteHeader(400)
		response.Write([]byte("No receipt found for that ID."))
	}

	id := mux.Vars(request)["id"]
	receiptUUID, _ := uuid.Parse(id)
	returnPoints, found := receiptPoints[receiptUUID]
	if !found {
		sendInvalidResponse()
	}

	responseStruct := struct {
		Points int `json:"points"`
	}{
		Points: returnPoints,
	}

	responseJson, err := json.Marshal(responseStruct)
	if err != nil {
		sendInvalidResponse()
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.Write(responseJson)
}
