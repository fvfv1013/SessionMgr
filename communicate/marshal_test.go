package communicate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	pb "sessionmgr/proto/pkg/ready_pb"
	"sessionmgr/proto/pkg/return_pb"
	"testing"
)

func TestMarshal(t *testing.T) {
	rds := return_pb.ReturnReady{
		ReadyList: []*pb.Ready{
			{SessionID: 99, DAtA: []byte("test")},
		},
	}
	bin, err := rds.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	newRds := return_pb.ReturnReady{}
	err = newRds.Unmarshal(bin)
	if err != nil {
		t.Fatal(err)
	}
	if newRds.ReadyList[0].SessionID != 99 || !bytes.Equal(newRds.ReadyList[0].DAtA, []byte("test")) {
		t.Fatalf("not equal")
	}
}

func TestZeroByte(t *testing.T) {
	rds := return_pb.ReturnReady{
		ReadyList: []*pb.Ready{
			{SessionID: 99, DAtA: []byte("t\000tt")},
		},
	}
	bin, err := rds.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	newRds := return_pb.ReturnReady{}
	err = newRds.Unmarshal(bin)
	if err != nil {
		t.Fatal(err)
	}
	if newRds.ReadyList[0].SessionID != 99 || !bytes.Equal(newRds.ReadyList[0].DAtA, []byte("t\000tt")) {
		t.Fatalf("not equal")
	}
}

func TestZeroIntJson(t *testing.T) {
	rds := return_pb.ReturnReady{
		ReadyList: []*pb.Ready{
			{SessionID: 0},
		},
	}
	str, err := json.Marshal(rds)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(str))
	newRds := return_pb.ReturnReady{}
	err = json.Unmarshal(str, &newRds)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(newRds)
	if newRds.ReadyList[0].SessionID != 0 {
		t.Fatalf("not equal")
	}
}

func TestZeroStringJson(t *testing.T) {
	rd := pb.Ready{
		SessionID: 0,
		DAtA:      make([]byte, 10),
	}
	rd.DAtA[5] = byte(0)
	rd.DAtA[6] = byte(1)
	rds := return_pb.ReturnReady{
		ReadyList: []*pb.Ready{
			&rd,
		},
	}
	str, err := json.Marshal(rds)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(str))
	newRds := return_pb.ReturnReady{}
	err = json.Unmarshal(str, &newRds)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(newRds)
	if newRds.ReadyList[0].SessionID != 0 {
		t.Fatalf("not equal")
	}
}

type ArraySample struct {
	Arr []int `json:"Arr"`
}

func TestZeroArrayJson(t *testing.T) {
	arr := []int{1, 2, 3, 4, 0, 5}
	arrWrap := ArraySample{arr}
	str, err := json.Marshal(arrWrap)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(str))
	newArr := ArraySample{}
	err = json.Unmarshal(str, &newArr)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(newArr.Arr)
	if !reflect.DeepEqual(newArr.Arr, arr) {
		t.Errorf("arr not equal")
	}
}
