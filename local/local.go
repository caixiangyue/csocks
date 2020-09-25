package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/caixiangyue/csocks"
)

var localAddr, serverAddr string

func init() {
	flag.StringVar(&localAddr, "l", "", "local address")
	flag.StringVar(&serverAddr, "s", "", "server address")
}

func main() {
	flag.Parse()

	config, _ := ioutil.ReadFile("config")
	cLocal, err := csocks.NewCLocal(localAddr, serverAddr, csocks.NewCipher(config))

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	cLocal.Listen(func(listenAddr net.Addr) { log.Printf("CLocalSocks启动啦%s", listenAddr.String()) })
}
