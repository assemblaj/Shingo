package shingo

import "encoding/json"

type Message interface {
	ToBytes() ([]byte, error)
	FromBytes([]byte) error
}

type MessageType byte

const (
	RegistrationRequest MessageType = iota
	NotRegisteredResponse
	RegisteredResponse
	AuthRequest
	CandidateRequest
	AuthResponse
	CandidateResponse
)

func ParseIncomingMessage(buf []byte) (IncomingMessage, error) {
	var incoming IncomingMessage
	err := incoming.FromBytes(buf)
	if err != nil {
		return IncomingMessage{}, err
	}
	return incoming, nil
}

func jsonMessageToBytes(msg Message) ([]byte, error) {
	msgBuf, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	} else {
		return msgBuf, nil
	}
}

func jsonMessageFromBytes(buf []byte, message Message) error {
	err := json.Unmarshal(buf, message)
	if err != nil {
		return err
	}
	return nil
}

type IncomingMessage struct {
	Data string      `json:"message"`
	Type MessageType `json:"messageType"`
}

func (im *IncomingMessage) ToBytes() ([]byte, error) {
	return jsonMessageToBytes(im)
}

func (im *IncomingMessage) FromBytes(buf []byte) error {
	err := jsonMessageFromBytes(buf, im)
	if err != nil {
		return err
	}
	return nil
}

type RegistrationResponseMessage struct {
	Client Client `json:"client"`
}

func (rr *RegistrationResponseMessage) ToBytes() ([]byte, error) {
	return jsonMessageToBytes(rr)
}

func (rr *RegistrationResponseMessage) FromBytes(buf []byte) error {
	err := jsonMessageFromBytes(buf, rr)
	if err != nil {
		return err
	}
	return nil
}

type AuthResponseMessage struct {
	Ufrag      string `json:"ufrag"`
	Pwd        string `json:"pwd"`
	Registered bool   `json:"registered"`
}

func (ar *AuthResponseMessage) ToBytes() ([]byte, error) {
	return jsonMessageToBytes(ar)
}

func (ar *AuthResponseMessage) FromBytes(buf []byte) error {
	err := jsonMessageFromBytes(buf, ar)
	if err != nil {
		return err
	}
	return nil
}

type CandidateResponseMessage struct {
	Candidate string `json:"candidate"`
}

func (cr *CandidateResponseMessage) ToBytes() ([]byte, error) {
	return jsonMessageToBytes(cr)
}

func (cr *CandidateResponseMessage) FromBytes(buf []byte) error {
	err := jsonMessageFromBytes(buf, cr)
	if err != nil {
		return err
	}
	return nil
}

type RegistrationMessage struct {
	Client Client `json:"client"`
}

func (rm *RegistrationMessage) ToBytes() ([]byte, error) {
	return jsonMessageToBytes(rm)
}

func (rm *RegistrationMessage) FromBytes(buf []byte) error {
	err := jsonMessageFromBytes(buf, rm)
	if err != nil {
		return err
	}
	return nil
}

type AuthRequestMessage struct {
	Client Client `json:"client"`
	Ufrag  string `json:"ufrag"`
	Pwd    string `json:"pwd"`
}

func (ar *AuthRequestMessage) ToBytes() ([]byte, error) {
	return jsonMessageToBytes(ar)
}

func (ar *AuthRequestMessage) FromBytes(buf []byte) error {
	err := jsonMessageFromBytes(buf, ar)
	if err != nil {
		return err
	}
	return nil
}

type CandidateRequestMessage struct {
	Client    Client `json:"client"`
	Candidate string `json:"candidate"`
}

func (cr *CandidateRequestMessage) ToBytes() ([]byte, error) {
	return jsonMessageToBytes(cr)
}

func (cr *CandidateRequestMessage) FromBytes(buf []byte) error {
	err := jsonMessageFromBytes(buf, cr)
	if err != nil {
		return err
	}
	return nil
}
