package communicate

import (
	"encoding/json"
)

type CommString struct {
	communicate *Communicate
}

func NewCommString(ConfPath string) *CommString {
	communicate, err := NewCommunicate(ConfPath)
	if err != nil {
		return &CommString{}
	}
	return &CommString{
		communicate: communicate,
	}
}

func (c *CommString) CreateSession(SessionID int32) string {
	returnValue := c.communicate.CreateSession(SessionID)
	jsonBytes, err := json.Marshal(returnValue)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func (c *CommString) Offer(SessionID int32) string {
	returnValue := c.communicate.Offer(SessionID)
	jsonBytes, err := json.Marshal(returnValue)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func (c *CommString) JoinSession(SessionID int32, sdpBase64 string) string {
	returnValue := c.communicate.JoinSession(SessionID, sdpBase64)
	jsonBytes, err := json.Marshal(returnValue)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func (c *CommString) Answer(SessionID int32) string {
	returnValue := c.communicate.Answer(SessionID)
	jsonBytes, err := json.Marshal(returnValue)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func (c *CommString) ConfirmAnswer(SessionID int32, sdpBase64 string) string {
	returnValue := c.communicate.ConfirmAnswer(SessionID, sdpBase64)
	jsonBytes, err := json.Marshal(returnValue)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func (c *CommString) Send(SessionID int32, dAtA []byte) string {
	returnValue := c.communicate.Send(SessionID, dAtA)
	jsonBytes, err := json.Marshal(returnValue)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func (c *CommString) Ready() string {
	returnValue := c.communicate.Ready()
	jsonBytes, err := json.Marshal(returnValue)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func (c *CommString) DropSession(SessionID int32) string {
	returnValue := c.communicate.DropSession(SessionID)
	jsonBytes, err := json.Marshal(returnValue)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func (c *CommString) ReloadConfig(ConfPath string) string {
	returnValue := c.communicate.ReloadConfig(ConfPath)
	jsonBytes, err := json.Marshal(returnValue)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func (c *CommString) Discard() string {
	returnValue := c.communicate.Discard()
	jsonBytes, err := json.Marshal(returnValue)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}
