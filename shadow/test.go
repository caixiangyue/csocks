package main

import (
	"csocks"
	"fmt"
	"io/ioutil"
)

func main() {
	config, _ := ioutil.ReadFile("config")
	cipher := csocks.NewCipher(config)

	test := []byte{5, 1, 0}
	cipher.Encode(test)
	fmt.Println(test[:])
	cipher.Decode(test)
	fmt.Println(test[:])
}
