package csocks

import (
	"io/ioutil"
	"math/rand"
	"time"
)

const passwordLength = 256

type shadow [passwordLength]byte

func init() {
	rand.Seed(time.Now().Unix())
}

func RandPassword() {

	temp := make([]byte, passwordLength)

	for i := 0; i < 256; i++ {
		temp[i] = byte(i)
	}

	shadow := &shadow{}
	for i := 0; i < 256; i++ {
		index := rand.Intn(len(temp))
		shadow[i] = temp[index]
		temp = append(temp[:index], temp[index+1:]...)
	}
	ioutil.WriteFile("config", shadow[:], 0644)
}
