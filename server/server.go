package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/caixiangyue/csocks"
)

var localAddr string

func init() {
	flag.StringVar(&localAddr, "l", "0.0.0.0:23456", "local address")
}

func main() {
	config, _ := ioutil.ReadFile("config")

	cServer, err := csocks.NewCServer(localAddr, csocks.NewCipher(config))

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	cServer.Listen(func(listenAddr net.Addr) { log.Printf("CServerrSocks启动啦%s", listenAddr.String()) })
}
