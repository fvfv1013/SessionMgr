package sessionmgr

import (
	"github.com/pion/webrtc/v4"
	"sessionmgr/dbg"
	"sessionmgr/util"
	"time"
)

// Session describe a talk
type Session struct {
	Connection *webrtc.PeerConnection
	DataCh     *webrtc.DataChannel
	LastUsed   time.Time
}

// NewSession create session
func NewSession(config *webrtc.Configuration) (*Session, error) {
	conn, err := webrtc.NewPeerConnection(*config)
	if err != nil {
		return nil, err
	}
	s := &Session{
		Connection: conn,
		DataCh:     nil,
		LastUsed:   time.Now(),
	}
	return s, nil
}

func (s *Session) OfferReady() bool {
	if state := s.Connection.ICEGatheringState(); state == webrtc.ICEGatheringStateComplete {
		return true
	}
	return false
}

func (s *Session) Offer() (string, error) {
	offer := s.Connection.LocalDescription()
	dbg.Println(dbg.SESSION, "offer details:\n", offer)
	sdpBase64, err := util.EncodeSDP(offer)
	if err != nil {
		dbg.Println(dbg.MANAGER, err)
		return "", err
	}
	return sdpBase64, nil
}

func (s *Session) AnswerReady() bool {
	if state := s.Connection.ICEGatheringState(); state == webrtc.ICEGatheringStateComplete {
		return true
	}
	return false
}

func (s *Session) Answer() (string, error) {
	answer := s.Connection.LocalDescription()
	dbg.Println(dbg.SESSION, "answer details:\n", answer)
	sdpBase64, err := util.EncodeSDP(answer)
	if err != nil {
		dbg.Println(dbg.MANAGER, err)
		return "", err
	}
	return sdpBase64, nil
}

func (s *Session) ConfirmAnswer(sdpBase64 string) error {
	if err := util.ValidateSDP(sdpBase64); err != nil {
		dbg.Println(dbg.SESSION, err)
		return err
	}
	answer, err := util.DecodeSDP(sdpBase64)
	if err != nil {
		dbg.Println(dbg.SESSION, err)
		return err
	}
	if err = s.Connection.SetRemoteDescription(*answer); err != nil {
		return err
	}
	return nil
}

func (s *Session) Send(dAtA []byte) error {
	if state := s.DataCh.ReadyState(); state != webrtc.DataChannelStateOpen {
		return ErrWait
	}
	if err := s.DataCh.Send(dAtA); err != nil {
		return err
	}
	return nil
}

func (s *Session) RecentActive() {
	s.LastUsed = time.Now()
}

func (s *Session) ReportCandidate() {
	s.Connection.OnICECandidate(func(c *webrtc.ICECandidate) {
		dbg.Println(dbg.ICE, "candidate found:", c)
	})
}
