package main

import (
	"csocks"
	"io/ioutil"
	"log"
	"net"
	"os"
	"runtime"
)

var (
	localAddr  = "127.0.0.1:1090"
	remoteAddr = "127.0.0.1:1085"
)

func main() {

	runtime.GOMAXPROCS(1)
	config, _ := ioutil.ReadFile("config")
	cLocal, err := csocks.NewCLocal(localAddr, remoteAddr, csocks.NewCipher(config))

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	cLocal.Listen(func(listenAddr net.Addr) { log.Printf("CLocalSocks启动啦%s", listenAddr.String()) })
}
