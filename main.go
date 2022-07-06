package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type UpdateRequest struct {
	Command string
	Data    string
}
type Move struct {
	cardValue   int
	opponentBet string
}

type Blinds struct {
	smallBlind int
	bigBlind   int
}

var m map[string]int

var moveDecision Move
var UpdateReq UpdateRequest
var blinds Blinds

var onButton bool

func main() {

	m = make(map[string]int)

	m["2"] = 2
	m["3"] = 3
	m["4"] = 4
	m["5"] = 5
	m["6"] = 6
	m["7"] = 7
	m["8"] = 8
	m["9"] = 9
	m["10"] = 10
	m["J"] = 11
	m["Q"] = 12
	m["K"] = 13
	m["A"] = 14

	http.HandleFunc("/move", moveHandler)

	http.HandleFunc("/start", startHandler)

	http.HandleFunc("/update", updateHandler)

	fmt.Printf("Starting server at port 80\n")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}

func moveHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println(onButton)

	decision := decisionEngine(moveDecision.cardValue, moveDecision.opponentBet)

	fmt.Println("Our Card value is " + strconv.Itoa(moveDecision.cardValue) + " opponent move is " + moveDecision.opponentBet + " our decision was " + decision)
	sendResponse(w, decision)
}

func startHandler(w http.ResponseWriter, r *http.Request) {

	blinds.bigBlind, _ = strconv.Atoi(r.PostFormValue("BIG_BLIND"))
	blinds.smallBlind, _ = strconv.Atoi(r.PostFormValue("SMALL_BLIND"))

	opponent := r.PostFormValue("OPPONENT_NAME")
	if len(opponent) > 0 {
		fmt.Println(opponent)
	}

	fmt.Printf("BigBlind : %v, SmallBlind: %v", blinds.bigBlind, blinds.smallBlind)
	w.WriteHeader(http.StatusOK)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	onButton = false

	parseUpdateRequest(r)

	w.WriteHeader(http.StatusOK)
}

func parseUpdateRequest(r *http.Request) {
	UpdateReq.Command = r.PostFormValue("COMMAND")
	UpdateReq.Data = r.PostFormValue("DATA")

	fmt.Println("Command is: " + UpdateReq.Command + " Data is: " + UpdateReq.Data)

	if UpdateReq.Command == "CARD" {
		moveDecision.cardValue = m[UpdateReq.Data]
	} else if UpdateReq.Command == "OPPONENT_MOVE" {
		moveDecision.opponentBet = UpdateReq.Data
	} else if UpdateReq.Command == "RECEIVE_BUTTON" {
		onButton = true
	} else if UpdateReq.Command == "OPPONENT_CARD" {
		onButton = false
		fmt.Printf("Opponent card is %s, did we lose: %t\n ************** \n", UpdateReq.Data, (m[UpdateReq.Data] > moveDecision.cardValue))
	}
}

func decisionEngine(cardValue int, opponentMove string) (move string) {
	if cardValue >= 13 {
		move = "BET:100000"
	} else if cardValue >= 12 {
		move = "BET:150"
	} else if isOpponentBettingHigh(opponentMove) && cardValue <= 11 {
		move = "FOLD"
	} else if (cardValue == 11 || cardValue == 12) && !isOpponentBettingHigh(opponentMove) {
		move = "BET:1"
	} else if cardValue <= 10 && cardValue > 7 && isOpponentStalling(opponentMove) {
		move = "BET:1"
	} else if cardValue <= 10 && cardValue > 7 {
		move = "CALL"
	} else if cardValue < 10 && onButton {
		move = "CALL"
	} else {
		move = "FOLD"
	}
	return
}

func isOpponentBettingHigh(opponentMove string) bool {
	if strings.Contains(opponentMove, "BET:") {
		strArr := strings.Split(opponentMove, ":")
		bet, _ := strconv.Atoi(strArr[1])
		return bet >= 5*blinds.bigBlind
	}
	return false
}

func isOpponentStalling(opponentMove string) bool {
	return strings.Contains(opponentMove, "CALL")
}

func sendResponse(w http.ResponseWriter, r string) {
	w.Write([]byte(r))
	w.WriteHeader(http.StatusOK)
}
