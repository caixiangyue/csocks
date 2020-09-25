package csocks

import (
	"io"
	"log"
	"net"
)

type SecureTCPConn struct {
	io.ReadWriteCloser
	Cipher *cipher
}

func (secureTcpConn *SecureTCPConn) DecodeRead(bs []byte) (n int, err error) {
	n, err = secureTcpConn.Read(bs)
	if err != nil {
		return
	}
	secureTcpConn.Cipher.Decode(bs[:n])
	return
}

func (secureTcpConn *SecureTCPConn) EncodeWrite(bs []byte) (int, error) {
	secureTcpConn.Cipher.Encode(bs)
	return secureTcpConn.Write(bs)
}

func (remoteConn *SecureTCPConn) DecodeCopy(localConn io.Writer) error {
	var buf = make([]byte, 1024)
	for {
		readCount, readErr := remoteConn.DecodeRead(buf)
		if readErr != nil {
			return readErr
		}
		if readCount > 0 {
			writeCount, writeErr := localConn.Write(buf[0:readCount])
			if writeErr != nil {
				return writeErr
			}
			if readCount != writeCount {
				return io.ErrShortWrite
			}
		}
	}
}

func (localConn *SecureTCPConn) EncodeCopy(remoteConn io.ReadWriteCloser) error {
	var data = make([]byte, 1024)
	for {
		readCount, readErr := localConn.Read(data)
		if readErr != nil {
			return readErr
		}
		if readCount > 0 {
			writeCount, writeErr := (&SecureTCPConn{remoteConn, localConn.Cipher}).EncodeWrite(data[0:readCount])
			if writeErr != nil {
				return writeErr
			}
			if readCount != writeCount {
				return io.ErrShortWrite
			}
		}

	}
}

func DialTcpSecure(remoteAddr *net.TCPAddr, cipher *cipher) (*SecureTCPConn, error) {
	remoteConn, err := net.DialTCP("tcp", nil, remoteAddr)
	if err != nil {
		return nil, err
	}
	return &SecureTCPConn{remoteConn, cipher}, nil
}

func ListenSecureTCP(localTcpAddr *net.TCPAddr, cipher *cipher, handleConn func(conn *SecureTCPConn), printInfo func(listenAddr net.Addr)) error {

	listener, err := net.ListenTCP("tcp", localTcpAddr)
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
		go handleConn(&SecureTCPConn{conn, cipher})
	}

}
