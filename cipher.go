package csocks

type Cipher struct {
	decodeShadow []byte
	encodeShadow []byte
}

func NewCipher(shadow []byte) *Cipher {
	decodeShadow := make([]byte, 256)
	for i, v := range shadow {
		decodeShadow[v] = byte(i)
	}
	return &Cipher{decodeShadow, shadow}
}

func (cipher *Cipher) Encode(buf []byte) []byte {
	for i, v := range buf {
		buf[i] = cipher.encodeShadow[v]
	}
	return buf
}

func (cipher *Cipher) Decode(buf []byte) []byte {
	for i, v := range buf {
		buf[i] = cipher.decodeShadow[v]
	}
	return buf
}
