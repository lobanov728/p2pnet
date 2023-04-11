package protocol

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"log"
)

const bits = 2048
const encChunk = bits / 8
const bytesInChunk = 4

func GenerateKeys() (privateKey []byte, publicKey []byte) {
	privKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return x509.MarshalPKCS1PrivateKey(privKey), x509.MarshalPKCS1PublicKey(&privKey.PublicKey)
}

func EncryptMessage(publicKey, msg []byte) ([]byte, error) {
	rsaPubKey, err := x509.ParsePKCS1PublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	msgLen := len(msg)
	resp := make([]byte, 0)

	for i := 0; i < msgLen; i += bytesInChunk {
		end := i + bytesInChunk

		if end > msgLen {
			end = msgLen
		}

		subMsg := msg[i:end]
		encMsg, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPubKey, subMsg)
		if err != nil {
			return nil, err
		}
		resp = append(resp, encMsg...)
	}

	return resp, nil
}

func DecryptMessage(privateKey, msg []byte) ([]byte, error) {
	rsaPrivKey, err := x509.ParsePKCS1PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	msgLen := len(msg)
	resp := make([]byte, 0)

	for i := 0; i < msgLen; i += encChunk {
		end := i + encChunk
		if end > msgLen {
			end = msgLen
		}

		subMsg := msg[i:end]
		encMsg, err := rsa.DecryptPKCS1v15(rand.Reader, rsaPrivKey, subMsg)
		if err != nil {
			return nil, err
		}
		resp = append(resp, encMsg...)
	}

	return resp, nil
}
