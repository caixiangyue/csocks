package csocks

import (
	"io/ioutil"
	"testing"
)

func TestRandPassword(t *testing.T) {
	config, _ := ioutil.ReadFile("./shadow/config")
	cipher := NewCipher(config)

	expected := []byte{5, 1, 0}
	actual := []byte{5, 1, 0}

	cipher.Encode(actual)
	cipher.Decode(actual)

	if len(actual) != len(expected) {
		t.Error("编解码错误!")
	}
	for i, v := range expected {
		if v != actual[i] {
			t.Error("编解码错误!")
		}
	}
}
