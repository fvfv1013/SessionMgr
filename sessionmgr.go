package sessionmgr

import (
	"github.com/pion/webrtc/v4"
	"sessionmgr/conf"
	"sessionmgr/dbg"
	pb "sessionmgr/proto/pkg/ready_pb"
	"sessionmgr/util"
	"sync"
	"sync/atomic"
	"time"
)

type SessionManagerImpl struct {
	mu           sync.Mutex
	config       *conf.Configuration
	sessionBook  map[int32]*Session
	readyChannel chan *pb.Ready
	discarded    atomic.Bool // not protected by mu
}

func NewSessionManagerImpl(ConfPath string) (*SessionManagerImpl, error) {
	config, err := conf.LoadConfig(ConfPath)
	if err != nil {
		return nil, err
	}
	s := &SessionManagerImpl{
		mu:           sync.Mutex{},
		config:       config,
		sessionBook:  make(map[int32]*Session),
		readyChannel: make(chan *pb.Ready, config.CacheSize),
		discarded:    atomic.Bool{},
	}
	s.enableLifeControl()
	return s, nil
}

func (s *SessionManagerImpl) CreateSession(SessionID int32) error {
	if s.discarded.Load() {
		return ErrCall
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, existed := s.sessionBook[SessionID]; existed {
		dbg.Println(dbg.MANAGER, "int repeated")
		return ErrID
	}
	session, err := NewSession(&s.config.WebrtcConf)
	if err != nil {
		return err
	}
	session.RecentActive()
	s.sessionBook[SessionID] = session
	if err = s.initA(SessionID); err != nil {
		return err
	}
	dbg.Println(dbg.MANAGER, "create session: ", SessionID)
	return nil
}

func (s *SessionManagerImpl) Offer(SessionID int32) (string, error) {
	if s.discarded.Load() {
		dbg.Println(dbg.MANAGER, ErrCall)
		return "", ErrCall
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	session, err := s.session(SessionID)
	if err != nil {
		return "", err
	}
	if ready := session.OfferReady(); !ready {
		return "", ErrWait
	}
	sdpBase64, err := session.Offer()
	if err != nil {
		return "", err
	}
	dbg.Println(dbg.MANAGER, "get offer: ", sdpBase64)
	return sdpBase64, nil
}

func (s *SessionManagerImpl) JoinSession(SessionID int32, sdpBase64 string) error {
	if s.discarded.Load() {
		dbg.Println(dbg.MANAGER, ErrCall)
		return ErrCall
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, existed := s.sessionBook[SessionID]; existed {
		dbg.Println(dbg.MANAGER, "SessionID repeated")
		return ErrID
	}
	if err := s.joinSession(SessionID, sdpBase64); err != nil {
		return err
	}
	dbg.Println(dbg.MANAGER, "join session: ", SessionID)
	return nil
}

func (s *SessionManagerImpl) Answer(SessionID int32) (string, error) {
	if s.discarded.Load() {
		dbg.Println(dbg.MANAGER, ErrCall)
		return "", ErrCall
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	session, err := s.session(SessionID)
	if err != nil {
		return "", err
	}
	if ready := session.OfferReady(); !ready {
		return "", ErrWait
	}
	sdpBase64, err := session.Answer()
	if err != nil {
		return "", err
	}
	dbg.Println(dbg.MANAGER, "get answer: ", sdpBase64)
	return sdpBase64, nil
}

func (s *SessionManagerImpl) ConfirmAnswer(SessionID int32, sdpBase64 string) error {
	if s.discarded.Load() {
		dbg.Println(dbg.MANAGER, ErrCall)
		return ErrCall
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	session, err := s.session(SessionID)
	if err != nil {
		return err
	}
	if err := session.ConfirmAnswer(sdpBase64); err != nil {
		return err
	}
	dbg.Println(dbg.MANAGER, "confirm answer: ", sdpBase64)
	return nil
}

func (s *SessionManagerImpl) Send(SessionID int32, dAtA []byte) error {
	if s.discarded.Load() {
		dbg.Println(dbg.MANAGER, ErrCall)
		return ErrCall
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	session, err := s.session(SessionID)
	if err != nil {
		return err
	}
	if err = session.Send(dAtA); err != nil {
		return err
	}
	dbg.Println(dbg.MANAGER, "send data: ", string(dAtA))
	return nil
}

func (s *SessionManagerImpl) Ready() ([]*pb.Ready, error) {
	rlist := make([]*pb.Ready, 0)
	for len(s.readyChannel) > 0 {
		rlist = append(rlist, <-s.readyChannel)
	}
	dbg.Println(dbg.MANAGER, "get ready list: ", rlist)
	return rlist, nil
}

func (s *SessionManagerImpl) DropSession(SessionID int32) error {
	if s.discarded.Load() {
		dbg.Println(dbg.MANAGER, ErrCall)
		return ErrCall
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dropSession(SessionID)
	dbg.Println(dbg.MANAGER, "drop session: ", SessionID)
	return nil
}

func (s *SessionManagerImpl) ReloadConfig(ConfPath string) error {
	if s.discarded.Load() {
		dbg.Println(dbg.MANAGER, ErrCall)
		return ErrCall
	}

	config, err := conf.LoadConfig(ConfPath)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	dbg.Println(dbg.MANAGER, "config reloaded")
	s.config = config
	return nil
}

func (s *SessionManagerImpl) Discard() error {
	s.discarded.Store(true)
	dbg.Println(dbg.MANAGER, "discard manager")
	return nil
}

func (s *SessionManagerImpl) lifeControl() {
	timeout := time.Second * time.Duration(s.config.SessionLifeCycle)
	ticker := time.NewTicker(timeout)
	for {
		select {
		case <-ticker.C:
			dbg.Println(dbg.MANAGER, "life control triggered")
			s.mu.Lock()
			deadSession := make([]int32, 0)
			for SessionID, session := range s.sessionBook {
				if time.Since(session.LastUsed) > timeout {
					deadSession = append(deadSession, SessionID)
				}
			}
			for _, sessionID := range deadSession {
				s.dropSession(sessionID)
				dbg.Println(dbg.MANAGER, "actively drop session:", sessionID)
			}
			s.mu.Unlock()
		default:
			if s.discarded.Load() {
				s.mu.Lock()
				deadSession := make([]int32, 0)
				for SessionID, _ := range s.sessionBook {
					deadSession = append(deadSession, SessionID)
				}
				for _, sessionID := range deadSession {
					s.dropSession(sessionID)
					dbg.Println(dbg.MANAGER, "actively drop session:", sessionID)
				}
				s.mu.Unlock()
				return
			}
		}
	}
}

func (s *SessionManagerImpl) dropSession(SessionID int32) {
	session := s.sessionBook[SessionID]
	if session == nil {
		return
	}
	err := session.Connection.Close()
	if err != nil {
		dbg.Println(dbg.MANAGER, "drop session error: ", err)
	}
	delete(s.sessionBook, SessionID)
}

func (s *SessionManagerImpl) initA(SessionID int32) error {
	// 1. passively drop session
	if err := s.moniterLost(SessionID); err != nil {
		return err
	}
	if err := s.reportCandidate(SessionID); err != nil {
		return err
	}
	// 2. create dataCh
	if err := s.createDataCh(SessionID); err != nil {
		return err
	}
	// 3. set local state
	if err := s.prepareOffer(SessionID); err != nil {
		return err
	}
	return nil
}

func (s *SessionManagerImpl) prepareOffer(SessionID int32) error {
	session := s.sessionBook[SessionID]
	if session == nil {
		return ErrLost
	}

	initOffer, err := session.Connection.CreateOffer(nil)
	if err != nil {
		dbg.Println(dbg.MANAGER, err)
		return err
	}
	if err = session.Connection.SetLocalDescription(initOffer); err != nil {
		dbg.Println(dbg.MANAGER, err)
		return err
	}
	return nil
}

func (s *SessionManagerImpl) createDataCh(SessionID int32) error {
	session := s.sessionBook[SessionID]
	if session == nil {
		return ErrLost
	}
	dataCh, err := session.Connection.CreateDataChannel("data", nil)
	if err != nil {
		dbg.Println(dbg.MANAGER, err)
		return err
	}
	session.DataCh = dataCh
	dataCh.OnMessage(func(msg webrtc.DataChannelMessage) {
		dbg.Println(dbg.SESSION, "[A] receive:", string(msg.Data))
		if s.discarded.Load() {
			return
		}
		s.readyChannel <- &pb.Ready{
			SessionID: SessionID,
			DAtA:      msg.Data,
		}
	})
	return nil
}

func (s *SessionManagerImpl) moniterLost(SessionID int32) error {
	session := s.sessionBook[SessionID]
	if session == nil {
		return ErrLost
	}
	session.Connection.OnConnectionStateChange(func(connectionState webrtc.PeerConnectionState) {
		dbg.Println(dbg.SESSION, "connection state changed:", connectionState)
		if s.discarded.Load() {
			return
		}
		switch connectionState {
		case webrtc.PeerConnectionStateClosed, webrtc.PeerConnectionStateDisconnected, webrtc.PeerConnectionStateFailed:
			s.mu.Lock()
			defer s.mu.Unlock()
			s.dropSession(SessionID)
			dbg.Println(dbg.SESSION, "passively drop session", SessionID)
		default:
		}
	})
	return nil
}

func (s *SessionManagerImpl) joinSession(SessionID int32, sdpBase64 string) error {
	if err := util.ValidateSDP(sdpBase64); err != nil {
		dbg.Println(dbg.SESSION, err)
		return err
	}
	offer, err := util.DecodeSDP(sdpBase64)
	if err != nil {
		dbg.Println(dbg.SESSION, err)
		return err
	}
	if _, existed := s.sessionBook[SessionID]; existed {
		dbg.Println(dbg.MANAGER, "int repeated")
		return ErrID
	}

	session, err := NewSession(&s.config.WebrtcConf)
	if err != nil {
		return err
	}
	s.sessionBook[SessionID] = session
	if err = s.initB(SessionID, offer); err != nil {
		return err
	}
	return nil
}

func (s *SessionManagerImpl) initB(SessionID int32, offer *webrtc.SessionDescription) error {
	// 1. passively drop session
	if err := s.moniterLost(SessionID); err != nil {
		return err
	}
	if err := s.reportCandidate(SessionID); err != nil {
		return err
	}
	// 2. transmit
	if err := s.waitDataCh(SessionID); err != nil {
		return err
	}
	// 3. set sdp
	if err := s.prepareAnswer(SessionID, offer); err != nil {
		return err
	}
	return nil
}

func (s *SessionManagerImpl) prepareAnswer(SessionID int32, offer *webrtc.SessionDescription) error {
	session := s.sessionBook[SessionID]
	if session == nil {
		dbg.Println(dbg.MANAGER, ErrLost)
		return ErrLost
	}
	if err := session.Connection.SetRemoteDescription(*offer); err != nil {
		dbg.Println(dbg.SESSION, err)
		return err
	}
	initAnswer, err := session.Connection.CreateAnswer(nil)
	if err != nil {
		dbg.Println(dbg.SESSION, err)
		return err
	}
	if err = session.Connection.SetLocalDescription(initAnswer); err != nil {
		dbg.Println(dbg.SESSION, err)
		return err
	}
	return nil
}

func (s *SessionManagerImpl) waitDataCh(SessionID int32) error {
	session := s.sessionBook[SessionID]
	if session == nil {
		dbg.Println(dbg.MANAGER, ErrLost)
		return ErrLost
	}
	session.Connection.OnDataChannel(func(channel *webrtc.DataChannel) {
		dbg.Println(dbg.ICE, "[B] get dataCh")
		if s.discarded.Load() {
			return
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		session.DataCh = channel
		session.DataCh.OnMessage(func(msg webrtc.DataChannelMessage) {
			dbg.Println(dbg.SESSION, "[B] receive message:", string(msg.Data))
			s.readyChannel <- &pb.Ready{
				SessionID: SessionID,
				DAtA:      msg.Data,
			}
		})
	})
	return nil
}

func (s *SessionManagerImpl) session(SessionID int32) (*Session, error) {
	session := s.sessionBook[SessionID]
	if session == nil {
		dbg.Println(dbg.SESSION, ErrLost)
		return nil, ErrLost
	}
	session.RecentActive()
	return session, nil
}

func (s *SessionManagerImpl) enableLifeControl() {
	go s.lifeControl()
}

func (s *SessionManagerImpl) reportCandidate(SessionID int32) error {
	session := s.sessionBook[SessionID]
	if session == nil {
		dbg.Println(dbg.MANAGER, ErrLost)
		return ErrLost
	}
	session.ReportCandidate()
	return nil
}
