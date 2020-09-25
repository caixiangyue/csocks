package csocks

import (
	"encoding/binary"
	"log"
	"net"
)

type CServer struct {
	localTcpAddr *net.TCPAddr
	cipher       *cipher
}

func NewCServer(listenAddr string, cipher *cipher) (*CServer, error) {
	listenTCPAddr, err := net.ResolveTCPAddr("tcp", listenAddr)

	if err != nil {
		return nil, err
	}
	return &CServer{listenTCPAddr, cipher}, nil
}

func (server *CServer) Listen(printInfo func(listenAddr net.Addr)) error {
	return ListenSecureTCP(server.localTcpAddr, server.cipher, server.handleConn, printInfo)
}
func (server *CServer) handleConn(localConn *SecureTCPConn) {
	defer localConn.Close()
	buf := make([]byte, 256)

	_, err := localConn.DecodeRead(buf)
	if err != nil || buf[0] != 0x05 {
		return
	}
	localConn.EncodeWrite([]byte{0x05, 0x00})
	n, err := localConn.DecodeRead(buf)

	if err != nil && n < 7 {
		return
	}
	if buf[1] != 0x01 {
		return
	}

	var dIP []byte
	switch buf[3] {
	case 0x01:
		dIP = buf[4 : 4+net.IPv4len]
	case 0x03:
		ipAddr, err := net.ResolveIPAddr("ip", string(buf[5:n-2]))

		if err != nil {
			return
		}
		dIP = ipAddr.IP
	case 0x04:
		dIP = buf[4 : 4+net.IPv6len]
	default:
		return
	}
	dPort := buf[n-2:]
	dstAddr := &net.TCPAddr{
		IP:   dIP,
		Port: int(binary.BigEndian.Uint16(dPort)),
	}

	dstServer, err := net.DialTCP("tcp", nil, dstAddr)
	if err != nil {
		log.Println(err)
		return
	} else {
		defer dstServer.Close()
		dstServer.SetLinger(0)
		localConn.EncodeWrite([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	}

	go func() {
		if err := localConn.DecodeCopy(dstServer); err != nil {
			localConn.Close()
			dstServer.Close()
		}
	}()

	(&SecureTCPConn{dstServer, localConn.Cipher}).EncodeCopy(localConn)
}
