package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/fvfv1013/sessionmgr"
	"github.com/fvfv1013/sessionmgr/dbg"
	pb "github.com/fvfv1013/sessionmgr/proto/pkg/sessionmgr_pb"
	"math/rand"
	"os"
	"time"
)

func main() {

	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Usage: ./webrtcdemo sender")
		fmt.Println("Usage: ./webrtcdemo receiver")
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
		fmt.Println("Usage: ./webrtcdemo sender")
		fmt.Println("Usage: ./webrtcdemo receiver")
	}
}

func startSender() {
	var sender sessionmgr.SessionManager
	var err error
	sender, err = sessionmgr.NewSessionManagerImpl()
	if err != nil {
		dbg.Fatal(dbg.ELSE, err)
	}

	// 1. create a session
	sessionID := rand.Int31()
	err = sender.CreateSession(sessionID)
	for err != nil {
		if !errors.Is(err, sessionmgr.ErrSessionID) {
			dbg.Fatal(dbg.ELSE, err)
		}
		sessionID = rand.Int31()
		err = sender.CreateSession(sessionID)
	}

	// 2. acquire offer
	offerSDP, err := sender.Offer(sessionID)
	for err != nil {
		if !errors.Is(err, sessionmgr.ErrWait) {
			dbg.Fatal(dbg.ELSE, err)
		}
		offerSDP, err = sender.Offer(sessionID)
	}

	// 3. print offer
	_ = offerSDP
	fmt.Println("offer:", offerSDP)

	// 4. read answer
	fmt.Println("input answer:")
	var answerSDP string
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscanf(reader, "%s", &answerSDP)
	err = sender.ConfirmAnswer(sessionID, answerSDP)
	if err != nil {
		dbg.Fatal(dbg.ELSE, err)
	}

	// 5. send data
	for {
		time.Sleep(2 * time.Second)
		data := []byte("Let's chat!")
		err = sender.Send(sessionID, data)
		for err != nil {
			if !errors.Is(err, sessionmgr.ErrWait) {
				dbg.Fatal(dbg.ELSE, err)
			}
			err = sender.Send(sessionID, data)
		}
		fmt.Println("send data:", string(data))
	}
}

func startReceiver() {
	var receiver sessionmgr.SessionManager
	var err error
	receiver, err = sessionmgr.NewSessionManagerImpl()
	if err != nil {
		dbg.Fatal(dbg.ELSE, err)
	}

	// 1. join session
	fmt.Println("input offerSDP:")
	var offerSDP string
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscanf(reader, "%s", &offerSDP)
	sessionID := rand.Int31()
	err = receiver.JoinSession(sessionID, offerSDP)
	for err != nil {
		if !errors.Is(err, sessionmgr.ErrSessionID) {
			dbg.Fatal(dbg.ELSE, err)
		}
		sessionID = rand.Int31()
		err = receiver.JoinSession(sessionID, offerSDP)
	}

	// 2. get answer
	answerSDP, err := receiver.Answer(sessionID)
	for err != nil {
		if !errors.Is(err, sessionmgr.ErrWait) {
			dbg.Fatal(dbg.ELSE, err)
		}
		answerSDP, err = receiver.Answer(sessionID)
	}
	_ = answerSDP
	fmt.Println("answer:", answerSDP)

	// 3. receive data
	for {
		time.Sleep(2 * time.Second)
		readys := make([]*pb.Ready, 0)
		readys, err = receiver.Ready()
		for err != nil {
			if !errors.Is(err, sessionmgr.ErrWait) {
				dbg.Fatal(dbg.ELSE, err)
			}
			readys, err = receiver.Ready()
		}
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
