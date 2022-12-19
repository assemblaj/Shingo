package shingo_test

import (
	"net"
	"testing"
	"time"

	shingo "github.com/assemblaj/Shingo"
)

type MockClientConn struct {
	addr    *net.TCPAddr
	sendMsg shingo.IncomingMessage
	recvMsg shingo.IncomingMessage
	server  *shingo.SignalingServer
}

func NewMockClientConn(ipAddress string, port int, server *shingo.SignalingServer) MockClientConn {
	return MockClientConn{
		addr: &net.TCPAddr{
			IP:   net.ParseIP(ipAddress),
			Port: port,
		},
		server: server,
	}
}
func (mc *MockClientConn) Read(b []byte) (n int, err error) {
	bytes, err := mc.sendMsg.ToBytes()
	if err != nil {
		return 0, err
	}
	copy(b, bytes)
	return len(bytes), nil
}

func (mc *MockClientConn) Write(b []byte) (n int, err error) {
	inc, err := shingo.ParseIncomingMessage(b)
	if err != nil {
		return 0, err
	}
	mc.recvMsg = inc
	return len(b), nil
}

func (mc *MockClientConn) Close() error {
	return nil
}

func (mc *MockClientConn) LocalAddr() net.Addr {
	return nil
}

func (mc *MockClientConn) RemoteAddr() net.Addr {
	return mc.addr
}

func (mc *MockClientConn) SetDeadline(t time.Time) error {
	return nil
}

func (mc *MockClientConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (mc *MockClientConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestServerRegisterClient(t *testing.T) {
	srcIpAddress := "192.127.2.1"
	dstIpAddress := "32.123.138.68"
	ss := shingo.NewSignalingServer("localhost")
	conn := NewMockClientConn(srcIpAddress, 2000, &ss)
	msg := shingo.RegistrationMessage{
		Client: shingo.Client{
			SrcIpAddress: srcIpAddress,
			DstIpAddress: dstIpAddress,
		},
	}
	inc, err := shingo.WrapMessage(&msg, shingo.RegistrationRequest)
	if err != nil {
		t.Errorf("Error wrapping message: %s", err)
	}
	err = ss.RegisterClient(inc, &conn)
	if err != nil {
		t.Errorf("Error registering client: %s", err)
	}
	_, err = ss.ValidateClient(conn.addr)
	if err != nil {
		t.Errorf("Error in client registration: %s", err)
	}
}

func TestServerRelayClientAuth(t *testing.T) {
}

func TestServerRelayCandidate(t *testing.T) {

}

func TestServerRelayMessage(t *testing.T) {

}

func TestServerValidateClient(t *testing.T) {

}

func TestClientRegister(t *testing.T) {

}

func TestClientRequestAuth(t *testing.T) {
}

func TestClientRequestCondidate(t *testing.T) {

}

func TestClientSendMessage(t *testing.T) {

}
