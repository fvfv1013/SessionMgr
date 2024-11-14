package communicate

import (
	"sessionmgr"
	"sessionmgr/proto/pkg/return_pb"
)

type Communicate struct {
	sessionMgr *sessionmgr.SessionManagerImpl
}

func NewCommunicate(ConfPath string) (*Communicate, error) {
	sessionMgr, err := sessionmgr.NewSessionManagerImpl(ConfPath)
	if err != nil {
		return nil, err
	}
	return &Communicate{
		sessionMgr: sessionMgr,
	}, nil
}

func (c *Communicate) CreateSession(SessionID int32) *return_pb.Return {
	err := c.sessionMgr.CreateSession(SessionID)
	var returnValue return_pb.Return
	if err != nil {
		returnValue.Err = ErrorWrap(err)
		return &returnValue
	}
	returnValue.CreateSessionReturn = &return_pb.ReturnCreateSession{}
	return &returnValue
}

func (c *Communicate) Offer(SessionID int32) *return_pb.Return {
	sdpBase64, err := c.sessionMgr.Offer(SessionID)
	var returnValue return_pb.Return
	if err != nil {
		returnValue.Err = ErrorWrap(err)
		return &returnValue
	}
	returnValue.OfferReturn = &return_pb.ReturnOffer{
		OfferBase64: sdpBase64,
	}
	return &returnValue
}

func (c *Communicate) JoinSession(SessionID int32, sdpBase64 string) *return_pb.Return {
	err := c.sessionMgr.JoinSession(SessionID, sdpBase64)
	var returnValue return_pb.Return
	if err != nil {
		returnValue.Err = ErrorWrap(err)
		return &returnValue
	}
	returnValue.JoinSessionReturn = &return_pb.ReturnJoinSession{}
	return &returnValue
}

func (c *Communicate) Answer(SessionID int32) *return_pb.Return {
	sdpBase64, err := c.sessionMgr.Answer(SessionID)
	var returnValue return_pb.Return
	if err != nil {
		returnValue.Err = ErrorWrap(err)
		return &returnValue
	}
	returnValue.AnswerReturn = &return_pb.ReturnAnswer{
		AnswerBase64: sdpBase64,
	}
	return &returnValue
}

func (c *Communicate) ConfirmAnswer(SessionID int32, sdpBase64 string) *return_pb.Return {
	err := c.sessionMgr.ConfirmAnswer(SessionID, sdpBase64)
	var returnValue return_pb.Return
	if err != nil {
		returnValue.Err = ErrorWrap(err)
		return &returnValue
	}
	returnValue.ConfirmAnswerReturn = &return_pb.ReturnConfirmAnswer{}
	return &returnValue
}

func (c *Communicate) Send(SessionID int32, dAtA []byte) *return_pb.Return {
	err := c.sessionMgr.Send(SessionID, dAtA)
	var returnValue return_pb.Return
	if err != nil {
		returnValue.Err = ErrorWrap(err)
		return &returnValue
	}
	returnValue.SendReturn = &return_pb.ReturnSend{}
	return &returnValue
}

func (c *Communicate) Ready() *return_pb.Return {
	rlist, err := c.sessionMgr.Ready()
	var returnValue return_pb.Return
	if err != nil {
		returnValue.Err = ErrorWrap(err)
		return &returnValue
	}
	returnValue.ReadyReturn = &return_pb.ReturnReady{
		ReadyList: rlist,
	}
	return &returnValue
}

func (c *Communicate) DropSession(SessionID int32) *return_pb.Return {
	err := c.sessionMgr.DropSession(SessionID)
	var returnValue return_pb.Return
	if err != nil {
		returnValue.Err = ErrorWrap(err)
		return &returnValue
	}
	returnValue.DropSessionReturn = &return_pb.ReturnDropSession{}
	return &returnValue
}

func (c *Communicate) ReloadConfig(ConfPath string) *return_pb.Return {
	err := c.sessionMgr.ReloadConfig(ConfPath)
	var returnValue return_pb.Return
	if err != nil {
		returnValue.Err = ErrorWrap(err)
		return &returnValue
	}
	returnValue.ReloadConfigReturn = &return_pb.ReturnReloadConfig{}
	return &returnValue
}

func (c *Communicate) Discard() *return_pb.Return {
	err := c.sessionMgr.Discard()
	var returnValue return_pb.Return
	if err != nil {
		returnValue.Err = ErrorWrap(err)
		return &returnValue
	}
	returnValue.DiscardReturn = &return_pb.ReturnDiscard{}
	return &returnValue
}
