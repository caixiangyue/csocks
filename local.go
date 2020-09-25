package csocks

import (
	"log"
	"net"
)

type CLocal struct {
	localTCPAddr  *net.TCPAddr
	remoteTCPAddr *net.TCPAddr
	cipher        *cipher
}

func NewCLocal(localAddr string, remoteAddr string, cipher *cipher) (*CLocal, error) {
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
	return ListenSecureTCP(local.localTCPAddr, local.cipher, local.handleConn, printInfo)
}

func (local *CLocal) handleConn(localConn *SecureTCPConn) {
	defer localConn.Close()

	remoteConn, err := DialTcpSecure(local.remoteTCPAddr, local.cipher)
	if err != nil {
		log.Println(err)
		return
	}

	defer remoteConn.Close()

	go func() {
		if err := remoteConn.DecodeCopy(localConn); err != nil {
			//log.Println(err)
			remoteConn.Close()
			localConn.Close()
		}

	}()

	if err := localConn.EncodeCopy(remoteConn); err != nil {
		//log.Println(err)
		remoteConn.Close()
		localConn.Close()
	}
}
