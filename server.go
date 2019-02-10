package csocks

import (
	"encoding/binary"
	"log"
	"net"
)

type CServer struct {
	listenAddr *net.TCPAddr
	cipher     *Cipher
}

func NewCServer(listenAddr string, cipher *Cipher) (*CServer, error) {
	listenTCPAddr, err := net.ResolveTCPAddr("tcp", listenAddr)

	if err != nil {
		return nil, err
	}
	return &CServer{listenTCPAddr, cipher}, nil
}

func (server *CServer) Listen(printInfo func(listenAddr net.Addr)) error {
	listener, err := net.ListenTCP("tcp", server.listenAddr)
	log.Println(server.listenAddr)
	if err != nil {
		log.Println(err)
		return err
	}
	defer listener.Close()
	if printInfo != nil {
		printInfo(listener.Addr())
	}
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}
		go server.handleConn(conn)
	}
}
func (server *CServer) handleConn(conn *net.TCPConn) {
	defer conn.Close()
	conn.SetLinger(0)
	buf := make([]byte, 256)

	_, err := conn.Read(buf)
	server.cipher.Decode(buf)
	if err != nil || buf[0] != 0x05 {
		return
	}
	conn.Write(server.cipher.Encode([]byte{0x05, 0x00}))
	n, err := conn.Read(buf)
	server.cipher.Decode(buf)

	if err != nil && n < 7 {
		return
	}
	if buf[1] != 0x01 {
		return
	}

	var dIP []byte
	switch buf[3] {
	case 0x01:
		//	IP V4 address: X'01'
		dIP = buf[4 : 4+net.IPv4len]
	case 0x03:
		ipAddr, err := net.ResolveIPAddr("ip", string(buf[5:n-2]))

		if err != nil {
			return
		}
		dIP = ipAddr.IP
	case 0x04:
		//	IP V6 address: X'04'
		dIP = buf[4 : 4+net.IPv6len]
	default:
		return
	}
	dPort := buf[n-2:]
	dstAddr := &net.TCPAddr{
		IP:   dIP,
		Port: int(binary.BigEndian.Uint16(dPort)),
	}
	log.Println(dstAddr.String())

	dstServer, err := net.DialTCP("tcp", nil, dstAddr)
	if err != nil {
		log.Println(err)
		return
	} else {
		defer dstServer.Close()
		dstServer.SetLinger(0)
		log.Println("连接成功")
		conn.Write(server.cipher.Encode([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}))
	}

	go func() {
		var buf = make([]byte, 1024)
		for {
			readCount, readErr := conn.Read(buf)
			server.cipher.Decode(buf)
			if readErr != nil {
				return
			}
			if readCount > 0 {
				writeCount, writeErr := dstServer.Write(buf[0:readCount])
				if writeErr != nil {
					return
				}
				if readCount != writeCount {
					return
				}
			}
		}
	}()

	var data = make([]byte, 1024)
	for {
		readCount, readErr := dstServer.Read(data)
		if readErr != nil {
			return
		}
		if readCount > 0 {
			writeCount, writeErr := conn.Write(server.cipher.Encode(data[0:readCount]))
			if writeErr != nil {
				return
			}
			if readCount != writeCount {
				return
			}
		}

	}
}
