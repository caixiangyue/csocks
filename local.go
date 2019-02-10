package csocks

import (
	"log"
	"net"
)

type CLocal struct {
	LocalTCPAddr  *net.TCPAddr
	RemoteTCPAddr *net.TCPAddr
	cipher        *Cipher
}

func NewCLocal(localAddr string, remoteAddr string, cipher *Cipher) (*CLocal, error) {
	localTCPAddr, err := net.ResolveTCPAddr("tcp", localAddr)

	if err != nil {
		return nil, err
	}
	remoteTCPAddr, err := net.ResolveTCPAddr("tcp", remoteAddr)

	if err != nil {
		return nil, err
	}
	return &CLocal{localTCPAddr, remoteTCPAddr, cipher}, nil
}

func (local *CLocal) Listen(printInfo func(listenAddr net.Addr)) error {
	listener, err := net.ListenTCP("tcp", local.LocalTCPAddr)
	if err != nil {
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
		conn.SetLinger(0)
		go local.handleConn(conn)
	}
}

func (local *CLocal) handleConn(conn *net.TCPConn) {
	defer conn.Close()

	remoteConn, err := net.DialTCP("tcp", nil, local.RemoteTCPAddr)
	if err != nil {
		log.Fatalln(err)
		return
	}

	defer remoteConn.Close()

	go func() {
		var buf = make([]byte, 1024)
		for {
			readCount, readErr := remoteConn.Read(buf)
			local.cipher.Decode(buf)
			if readErr != nil {
				return
			}
			if readCount > 0 {
				writeCount, writeErr := conn.Write(buf[0:readCount])
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
		readCount, readErr := conn.Read(data)
		if readErr != nil {
			return
		}
		if readCount > 0 {
			writeCount, writeErr := remoteConn.Write(local.cipher.Encode(data[0:readCount]))
			log.Println("写服务器数据:", data[0:readCount])
			if writeErr != nil {
				return
			}
			if readCount != writeCount {
				return
			}
		}

	}
}
