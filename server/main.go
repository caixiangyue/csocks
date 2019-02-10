package main

import (
	"csocks"
	"io/ioutil"
	"log"
	"net"
	"os"
	"runtime"
)

func main() {

	runtime.GOMAXPROCS(1)
	config, _ := ioutil.ReadFile("config")
	cServer, err := csocks.NewCServer("127.0.0.1:1085", csocks.NewCipher(config))

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	cServer.Listen(func(listenAddr net.Addr) { log.Printf("CServerrSocks启动啦%s", listenAddr.String()) })
}
