package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sessionmgr"
	"sessionmgr/communicate"
	"sessionmgr/dbg"
	"sessionmgr/proto/pkg/return_pb"
	"time"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Usage: ./example(.exe) sender")
		fmt.Println("Usage: ./example(.exe) receiver")
		PressAnyKey()
		return
	}
	err := dbg.Init(dbg.STDOUT)
	if err != nil {
		return
	}
	defer dbg.Close()
	cmd := args[0]
	switch cmd {
	case "sender":
		startSender()
	case "receiver":
		startReceiver()
	default:
		fmt.Println("Usage: ./example(.exe) sender")
		fmt.Println("Usage: ./example(.exe) receiver")
	}
}

func startSender() {
	var sender *communicate.CommString
	var RetStr string
	var Ret *return_pb.Return
	var err error
	sender = communicate.NewCommString("conf.json")

	// 1. create a session
	sessionID := rand.Int31()
	RetStr = sender.CreateSession(sessionID)
	Ret = Return(RetStr)
	err = communicate.ErrorUnwrap(Ret.Err)
	for err != nil {
		if !errors.Is(err, sessionmgr.ErrID) {
			dbg.Fatal(dbg.ELSE, err)
		}
		sessionID = rand.Int31()
		RetStr = sender.CreateSession(sessionID)
		Ret = Return(RetStr)
		err = communicate.ErrorUnwrap(Ret.Err)
	}

	// 2. acquire offer
	RetStr = sender.Offer(sessionID)
	Ret = Return(RetStr)
	err = communicate.ErrorUnwrap(Ret.Err)
	for err != nil {
		if !errors.Is(err, sessionmgr.ErrWait) {
			dbg.Fatal(dbg.ELSE, err)
		}
		RetStr = sender.Offer(sessionID)
		Ret = Return(RetStr)
		err = communicate.ErrorUnwrap(Ret.Err)
	}

	// 3. print offer
	offerSDP := Ret.OfferReturn.OfferBase64
	_ = offerSDP
	fmt.Println("offer:", offerSDP)

	// 4. read answer
	fmt.Println("input answer:")
	var answerSDP string
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscanf(reader, "%s", &answerSDP)
	RetStr = sender.ConfirmAnswer(sessionID, answerSDP)
	Ret = Return(RetStr)
	err = communicate.ErrorUnwrap(Ret.Err)
	if err != nil {
		dbg.Fatal(dbg.ELSE, err)
	}

	// 5. send data
	for {
		time.Sleep(2 * time.Second)
		data := []byte("Let's chat!")
		RetStr = sender.Send(sessionID, data)
		Ret = Return(RetStr)
		err = communicate.ErrorUnwrap(Ret.Err)
		for err != nil {
			if !errors.Is(err, sessionmgr.ErrWait) {
				dbg.Fatal(dbg.ELSE, err)
			}
			RetStr = sender.Send(sessionID, data)
			Ret = Return(RetStr)
			err = communicate.ErrorUnwrap(Ret.Err)
		}
		fmt.Println("send data:", string(data))
	}
}

func startReceiver() {
	var receiver *communicate.CommString
	var RetStr string
	var Ret *return_pb.Return
	var err error
	receiver = communicate.NewCommString("conf.json")

	// 1. join session
	fmt.Println("input offerSDP:")
	var offerSDP string
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscanf(reader, "%s", &offerSDP)
	sessionID := rand.Int31()
	RetStr = receiver.JoinSession(sessionID, offerSDP)
	Ret = Return(RetStr)
	err = communicate.ErrorUnwrap(Ret.Err)
	for err != nil {
		if !errors.Is(err, sessionmgr.ErrID) {
			dbg.Fatal(dbg.ELSE, err)
		}
		sessionID = rand.Int31()
		RetStr = receiver.JoinSession(sessionID, offerSDP)
		Ret = Return(RetStr)
		err = communicate.ErrorUnwrap(Ret.Err)
	}

	// 2. get answer
	RetStr = receiver.Answer(sessionID)
	Ret = Return(RetStr)
	err = communicate.ErrorUnwrap(Ret.Err)
	for err != nil {
		if !errors.Is(err, sessionmgr.ErrWait) {
			dbg.Fatal(dbg.ELSE, err)
		}
		RetStr = receiver.Answer(sessionID)
		Ret = Return(RetStr)
		err = communicate.ErrorUnwrap(Ret.Err)
	}
	answerSDP := Ret.AnswerReturn.AnswerBase64
	_ = answerSDP
	fmt.Println("answer:", answerSDP)

	// 3. receive data
	for {
		time.Sleep(2 * time.Second)
		RetStr = receiver.Ready()
		Ret = Return(RetStr)
		err = communicate.ErrorUnwrap(Ret.Err)
		for err != nil {
			if !errors.Is(err, sessionmgr.ErrWait) {
				dbg.Fatal(dbg.ELSE, err)
			}
			RetStr = receiver.Ready()
			Ret = Return(RetStr)
			err = communicate.ErrorUnwrap(Ret.Err)
		}
		readys := Ret.ReadyReturn.ReadyList
		for _, ready := range readys {
			fmt.Println("get message:", ready.SessionID, string(ready.DAtA))
		}
	}
}

func PressAnyKey() {
	fmt.Println("Press Any Key to Continue...")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadByte()
}

func Return(RetStr string) *return_pb.Return {
	Ret := return_pb.Return{}
	err := json.Unmarshal([]byte(RetStr), &Ret)
	if err != nil {
		dbg.Fatal(dbg.ELSE, err)
	}
	return &Ret
}
