package sessionmgr

import (
	"errors"
	pb "github.com/fvfv1013/sessionmgr/proto/pkg/sessionmgr_pb"
)

type SessionManager interface {
	// CreateSession for offer side to create a session
	CreateSession(SessionID int32) error
	// Offer return offer BASE64
	Offer(SessionID int32) (string, error)
	// JoinSession for answer side to join a session described by SDP
	JoinSession(SessionID int32, sdpBase64 string) error
	// Answer can be called after JoinSession
	Answer(SessionID int32) (string, error)
	// ConfirmAnswer confirms a session description
	ConfirmAnswer(SessionID int32, sdpBase64 string) error
	// Send add dAtA to send queue, it is not a obstructive function
	Send(SessionID int32, dAtA []byte) error
	// Ready return a list of received messages and where are they from
	Ready() ([]*pb.Ready, error)
	// DropSession allow user to drop a session
	// Warning: don't call DropSession easily, because it is very slow; not-used session will be shutdown automatically
	DropSession(SessionID int32) error
	// ReloadConfig will force SessionManager reload config from conf.json
	// warning: it may not work immediately
	ReloadConfig() error
	// Discard a SessionManager
	Discard() error
}

var ErrSessionID = errors.New("SessionID invalid")
var ErrCall = errors.New("manager has been discarded")
var ErrLost = errors.New("session lost")
var ErrWait = errors.New("service is not prepared")
var ErrSDP = errors.New("sdp invalid")
