package sessionmgr

import (
	"fmt"
	"github.com/fvfv1013/sessionmgr/dbg"
	"math/rand"
	"testing"
	"time"
)

func TestLifeControl(t *testing.T) {
	err := dbg.Init(dbg.STDOUT)
	if err != nil {
		return
	}
	defer dbg.Close()
	mgr, err := NewSessionManagerImpl()
	if err != nil {
		t.Fatal(err)
	}
	startNum := 10
	for i := 0; i < startNum; i++ {
		err := mgr.CreateSession(int32(i))
		if err != nil {
			t.Fatal(err)
		}
	}
	if len(mgr.sessionBook) != startNum {
		t.Fatalf("expected %d sessions, got %d", startNum, len(mgr.sessionBook))
	}

	// set lifecycle to 40 second for test
	mgr.config.SessionLifeCycle = 40
	mgr.enableLifeControl()

	// wait and check session number
	time.Sleep(20 * time.Second)
	if len(mgr.sessionBook) != startNum {
		t.Fatalf("expected %d sessions, got %d", startNum, len(mgr.sessionBook))
	}
	// trigger random session
	luckydog := rand.Int31() % int32(startNum)
	_, err = mgr.session(luckydog)
	if err != nil {
		t.Error(err)
	}

	// wait life control work
	time.Sleep(30 * time.Second)
	if len(mgr.sessionBook) != 1 {
		t.Errorf("expected %d sessions, got %d", 1, len(mgr.sessionBook))
	}

	// wait longer time
	time.Sleep(40 * time.Second)
	if mgr.sessionBook[luckydog] != nil {
		fmt.Println(time.Since(mgr.sessionBook[luckydog].LastUsed))
	}
	if len(mgr.sessionBook) != 0 {
		t.Errorf("expected %d sessions, got %d", 0, len(mgr.sessionBook))
	}
}

func TestDiscard(t *testing.T) {
	err := dbg.Init(dbg.STDOUT)
	if err != nil {
		return
	}
	defer dbg.Close()
	mgr, err := NewSessionManagerImpl()
	if err != nil {
		t.Fatal(err)
	}
	err = mgr.CreateSession(0)
	if err != nil {
		t.Fatal(err)
	}
	err = mgr.Discard()
	if err != nil {
		return
	}
	// drop session can be very slow
	time.Sleep(10 * time.Second)
	// there should be no sessions
	if len(mgr.sessionBook) != 0 {
		t.Errorf("expected %d sessions, got %d", 0, len(mgr.sessionBook))
	}
}
