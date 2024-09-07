package ciphers

import (
	"blockchain-transactions/internal/logger"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	openssl "github.com/Luzifer/go-openssl/v4"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

var secretKey string

func init() {
	//TODO get SecretKey
	secretKey = "204812730425442A472D2F423F452847"
}

func Encrypt(strToEncrypt string) string {

	o := openssl.New()

	enc, err := o.EncryptBytes(secretKey, []byte(strToEncrypt), openssl.BytesToKeyMD5)
	if err != nil {
		fmt.Printf("An error occurred: %s\n", err)
		return ""
	}

	return string(enc)
}

func Decrypt(strToDecrypt string) string {
	o := openssl.New()
	dec, err := o.DecryptBytes(secretKey, []byte(strToDecrypt), openssl.BytesToKeyMD5)
	if err != nil {
		fmt.Printf("An error occurred: %s\n", err)
		return err.Error()
	}
	return string(dec)
}

func GetSecret() string {
	return secretKey
}

func EncryptRSAOAEP(secretMessage string, publicKey rsa.PublicKey) string {
	label := []byte(secretKey)
	rng := rand.Reader
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, &publicKey, []byte(secretMessage), label)
	if err != nil {
		logger.Error.Printf("No se pudo cifrar el mensaje: error: " + err.Error())
		return ""
	}
	return base64.StdEncoding.EncodeToString(ciphertext)
}

func DecryptRSAOAEP(cipherText string, privateKey rsa.PrivateKey) string {
	ct, _ := base64.StdEncoding.DecodeString(cipherText)
	label := []byte(secretKey)
	rng := rand.Reader
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rng, &privateKey, ct, label)
	if err != nil {
		logger.Error.Printf("No se pudo decifrar el mensaje: error: " + err.Error())
		return ""
	}
	return string(plaintext)
}

func RsaPublicStringToRsaPublic(public string) *rsa.PublicKey {
	blockRsa, _ := pem.Decode([]byte(public))
	if blockRsa == nil {
		return nil
	}
	publicRsaPem, err := x509.ParsePKIXPublicKey(blockRsa.Bytes)
	if err != nil {
		return nil
	}

	publicRsa, ok := publicRsaPem.(*rsa.PublicKey)
	if !ok {
		return nil
	}
	return publicRsa
}

func RsaPrivateStringToRsaPrivate(public string) *rsa.PrivateKey {
	blockRsa, _ := pem.Decode([]byte(public))
	if blockRsa == nil {
		return nil
	}
	privateRsaPem, err := x509.ParsePKCS1PrivateKey(blockRsa.Bytes)
	if err != nil {
		return nil
	}

	return privateRsaPem
}

func GenerateKeyPairEcdsa() (string, string, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", err
	}

	pemPrivateKey, err := EncodePrivate(privateKey)
	if err != nil {
		return "", "", err
	}

	publicKey := privateKey.PublicKey
	pemPublicKey, err := EncodePublic(&publicKey)
	if err != nil {
		return "", "", err
	}

	return pemPrivateKey, pemPublicKey, nil
}

func EncodePrivate(privateKey *ecdsa.PrivateKey) (string, error) {
	encoded, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return "", err
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: encoded})
	return string(pemEncoded), nil
}

func EncodePublic(pubKey *ecdsa.PublicKey) (string, error) {

	encoded, err := x509.MarshalPKIXPublicKey(pubKey)

	if err != nil {
		return "", err
	}
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: encoded})

	return string(pemEncodedPub), nil
}

func DecodePrivate(pemEncodedPrivate string) (*ecdsa.PrivateKey, error) {
	blockPrivate, _ := pem.Decode([]byte(pemEncodedPrivate))
	x509EncodedPrivate := blockPrivate.Bytes
	privateKey, err := x509.ParseECPrivateKey(x509EncodedPrivate)
	return privateKey, err
}

func DecodePublic(pemEncodedPub string) (*ecdsa.PublicKey, error) {
	blockPub, _ := pem.Decode([]byte(pemEncodedPub))
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub)
	publicKey := genericPublicKey.(*ecdsa.PublicKey)
	return publicKey, err
}

func SignWithEcdsa(hash []byte, private ecdsa.PrivateKey) (string, error) {
	sign, err := ecdsa.SignASN1(rand.Reader, &private, hash)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sign), nil
}

func VerifySignWithEcdsa(hash []byte, pubKey ecdsa.PublicKey, sign []byte) (bool, error) {
	return ecdsa.VerifyASN1(&pubKey, hash, sign), nil
}

func StringToHashSha256(value string) string {
	h := sha256.New()
	h.Write([]byte(value))
	hash := h.Sum(nil)
	return fmt.Sprintf("%x", hash)
}

func GetPrivateKeyFormatJWK(privateKey string) error {
	privateKeyEcdsa, err := DecodePrivate(privateKey)
	if err != nil {
		return err
	}
	key, err := jwk.FromRaw(privateKeyEcdsa)
	if err != nil {
		fmt.Printf("failed to create ECDSA key: %s\n", err)
		return err
	}
	if _, ok := key.(jwk.ECDSAPrivateKey); !ok {
		fmt.Printf("expected jwk.ECDSAPrivateKey, got %T\n", key)
		return err
	}

	key.Set(jwk.KeyIDKey, uuid.New().String())

	/*buf, err := json.MarshalIndent(key, "", "  ")
	if err != nil {
		fmt.Printf("failed to marshal key into JSON: %s\n", err)
		return err
	}*/

	test := key.(jwk.ECDSAPrivateKey)
	fmt.Println(string(test.D()))

	res, _, err := jwk.EncodeX509(key)
	if err != nil {

	}
	res2, err := x509.ParseECPrivateKey([]byte(res))
	if err != nil {

	}
	fmt.Println(res2)
	return nil
}

func ParseJwkToEcdsa(jwkJSON []byte) *ecdsa.PrivateKey {
	return nil
}
