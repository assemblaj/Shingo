package shingo

import (
	"errors"
	"log"
	"net"
)

const (
	Address = "localhost:443"
)

type SignalingServer struct {
	address string
	clients map[string]*Client
}

func NewSignalingServer(address string) SignalingServer {
	return SignalingServer{
		address: address,
		clients: make(map[string]*Client),
	}
}

func (ss *SignalingServer) HandleConnection(conn net.Conn) {
	var err error
	message := make([]byte, 0)
	for {
		_, err = conn.Read(message)
		if err != nil {
			conn.Close()
			return
		}
		inc, err := ParseIncomingMessage(message)
		if err != nil {
			conn.Close()
			return
		}
		err = ss.HandleMessage(inc, conn)
		if err != nil {
			conn.Close()
			return
		}
	}
}

func (ss *SignalingServer) HandleMessage(inc IncomingMessage, conn net.Conn) error {
	var err error
	switch inc.Type {
	case RegistrationRequest:
		if err = ss.RegisterClient(inc, conn); err != nil {
			return err
		}
	case AuthRequest:
		if err = ss.RelayClientAuth(inc, conn.RemoteAddr()); err != nil {
			return err
		}
	case CandidateRequest:
		if err = ss.RelayCandidate(inc, conn.RemoteAddr()); err != nil {
			return err
		}
	}
	return nil
}

func (ss *SignalingServer) RegisterClient(inc IncomingMessage, conn net.Conn) error {
	var reg RegistrationMessage
	err := reg.FromBytes([]byte(inc.Data))
	if err != nil {
		return err
	}
	reg.Client.conn = conn
	ss.clients[reg.Client.SrcIpAddress] = &reg.Client
	return nil
}

func (ss *SignalingServer) RelayClientAuth(inc IncomingMessage, addr net.Addr) error {
	client, err := ss.ValidateClient(addr)
	if err != nil {
		return err
	}

	var authReq AuthRequestMessage
	err = authReq.FromBytes([]byte(inc.Data))
	if err != nil {
		return err
	}

	auth := AuthResponseMessage{
		Ufrag: authReq.Ufrag,
		Pwd:   authReq.Pwd,
	}
	msgBytes, err := BuildMessage(&auth, AuthResponse)
	if err != nil {
		return err
	}
	ss.RelayMessage(client, msgBytes)
	return nil
}

func (ss *SignalingServer) RelayCandidate(inc IncomingMessage, addr net.Addr) error {
	client, err := ss.ValidateClient(addr)
	if err != nil {
		return err
	}
	var candidateReq CandidateRequestMessage
	err = candidateReq.FromBytes([]byte(inc.Data))
	if err != nil {
		return err
	}
	candidate := CandidateResponseMessage{
		Candidate: candidateReq.Candidate,
	}
	msgBytes, err := BuildMessage(&candidate, CandidateResponse)
	if err != nil {
		return err
	}
	err = ss.RelayMessage(client, msgBytes)
	if err != nil {
		return err
	}
	return nil
}

func (ss *SignalingServer) RelayMessage(client *Client, message []byte) error {
	dst, ok := ss.clients[client.DstIpAddress]
	if !ok {
		return errors.New("destination client not registered to server")
	}
	_, err := dst.conn.Write(message)
	if err != nil {
		return err
	}
	return nil
}

func (ss *SignalingServer) ValidateClient(addr net.Addr) (*Client, error) {
	switch remoteAddr := addr.(type) {
	case *net.TCPAddr:
		srcIp := remoteAddr.IP.String()
		src, ok := ss.clients[srcIp]
		if !ok {
			return nil, errors.New("source client not registered to server")
		}
		return src, nil
	}
	return nil, errors.New("source client not using correct protocol")
}

type Client struct {
	SrcIpAddress string `json:"srcIp"`
	DstIpAddress string `json:"dstIp"`
	conn         net.Conn
}

type SignalingClient struct {
	client           Client
	remoteIp         string
	remotePort       string
	conn             net.Conn
	authChannel      chan AuthResponseMessage
	candidateChannel chan CandidateResponseMessage
}

func NewSignalingClient(srcIpAddress, dstIpAddress, remoteIp,
	remotePort string, conn net.Conn) SignalingClient {
	return SignalingClient{
		client: Client{
			SrcIpAddress: srcIpAddress,
			DstIpAddress: dstIpAddress,
		},
		remoteIp:         remoteIp,
		remotePort:       remotePort,
		conn:             conn,
		authChannel:      make(chan AuthResponseMessage),
		candidateChannel: make(chan CandidateResponseMessage),
	}
}

func (sc *SignalingClient) Register() error {
	rm := RegistrationMessage{
		Client: sc.client,
	}
	err := sc.SendMessage(&rm, RegistrationRequest)
	if err != nil {
		return err
	}
	return nil
}

func (sc *SignalingClient) RequestAuth(localUFrag, localPwd string) error {
	ar := AuthRequestMessage{
		Client: sc.client,
		Ufrag:  localUFrag,
		Pwd:    localPwd,
	}
	err := sc.SendMessage(&ar, AuthRequest)
	if err != nil {
		return err
	}
	return nil
}

func (sc *SignalingClient) RequestCondidate(candidate string) error {
	cr := CandidateRequestMessage{
		Client:    sc.client,
		Candidate: candidate,
	}
	err := sc.SendMessage(&cr, CandidateRequest)
	if err != nil {
		return err
	}
	return nil
}

func (sc *SignalingClient) RecieveMessage() (IncomingMessage, error) {
	buf := make([]byte, 0)
	_, err := sc.conn.Read(buf)
	if err != nil {
		return IncomingMessage{}, err
	}
	return ParseIncomingMessage(buf)
}

func (sc *SignalingClient) Listen() {
	for {
		inc, err := sc.RecieveMessage()
		if err != nil {
			sc.conn.Close()
			break
		}
		switch inc.Type {
		case AuthResponse:
			var auth AuthResponseMessage
			auth.FromBytes([]byte(inc.Data))
			sc.authChannel <- auth
		case CandidateResponse:
			var candidate CandidateResponseMessage
			candidate.FromBytes([]byte(inc.Data))
			sc.candidateChannel <- candidate
		}
	}
}

func (sc *SignalingClient) SendMessage(msg Message, messageType MessageType) error {
	msgBytes, err := BuildMessage(msg, messageType)
	if err != nil {
		return err
	}
	_, err = sc.conn.Write(msgBytes)
	if err != nil {
		return err
	}
	return nil
}

func BuildMessage(msg Message, messageType MessageType) ([]byte, error) {
	newMsg, err := WrapMessage(msg, messageType)
	if err != nil {
		return nil, err
	}
	msgBytes, err := newMsg.ToBytes()
	if err != nil {
		return nil, err
	}
	return msgBytes, nil
}

func WrapMessage(msg Message, messageType MessageType) (IncomingMessage, error) {
	bytes, err := msg.ToBytes()
	if err != nil {
		return IncomingMessage{}, err
	}
	newMsg := IncomingMessage{
		Data: string(bytes),
		Type: messageType,
	}
	return newMsg, nil
}

func main() {
	server := NewSignalingServer(Address)
	listener, err := net.Listen("tcp", server.address)
	if err != nil {
		log.Fatal("Error starting TCP Signaling Server. ")
	}
	for {
		conn, _ := listener.Accept()
		go server.HandleConnection(conn)
	}
}
