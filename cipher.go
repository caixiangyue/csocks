package csocks

type cipher struct {
	decodeShadow []byte
	encodeShadow []byte
}

func NewCipher(shadow []byte) *cipher {
	decodeShadow := make([]byte, 256)
	for i, v := range shadow {
		decodeShadow[v] = byte(i)
	}
	return &cipher{decodeShadow, shadow}
}

func (cipher *cipher) Encode(buf []byte) {
	for i, v := range buf {
		buf[i] = cipher.encodeShadow[v]
	}
}

func (cipher *cipher) Decode(buf []byte) {
	for i, v := range buf {
		buf[i] = cipher.decodeShadow[v]
	}
}
