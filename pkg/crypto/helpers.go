package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
	"time"
)

func GenerateRandomCertificateSerialNumber() (*big.Int, error) {
	limit := new(big.Int).Lsh(big.NewInt(1), 128)
	return rand.Int(rand.Reader, limit)
}

func SimpleSelfSignedRSAKeypair(cn string, days int) (key *rsa.PrivateKey, cert *x509.Certificate, err error) {
	key, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return key, cert, err
	}

	serialNumber, err := GenerateRandomCertificateSerialNumber()
	if err != nil {
		return key, cert, err
	}
	timeNow := time.Now()
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: cn,
		},
		NotBefore:             timeNow,
		NotAfter:              timeNow.Add(time.Duration(days) * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{cn},
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return key, cert, err
	}
	cert, err = x509.ParseCertificate(certBytes)
	return key, cert, err
}

func ReadPEMCertificateFile(path string) (*x509.Certificate, error) {
	certs, err := ReadPEMCertificatesFile(path)
	if err != nil {
		return nil, err
	}
	if len(certs) != 1 {
		return nil, errors.New("incorrect number of certificates")
	}
	return certs[0], nil
}

func ReadPEMCertificatesFile(path string) ([]*x509.Certificate, error) {
	pemData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var asn1data []byte
	rest := pemData
	for {
		var block *pem.Block
		block, rest = pem.Decode(rest)
		if block == nil || block.Type != "CERTIFICATE" {
			return nil, errors.New("failed to decode PEM block containing certificate")
		}
		asn1data = append(asn1data, block.Bytes...)
		if len(rest) == 0 {
			break
		}
	}
	return x509.ParseCertificates(asn1data)
}

func WritePEMCertificateFile(cert *x509.Certificate, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return pem.Encode(
		file,
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		})
}

func WritePEMRSAKeyFile(key *rsa.PrivateKey, path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0700)
	if err != nil {
		return err
	}
	defer file.Close()

	return pem.Encode(
		file,
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		})
}

// TopicFromCert extracts the push certificate topic from the provided certificate.
func TopicFromCert(cert *x509.Certificate) (string, error) {
	var oidASN1UserID = asn1.ObjectIdentifier{0, 9, 2342, 19200300, 100, 1, 1}
	for _, v := range cert.Subject.Names {
		if v.Type.Equal(oidASN1UserID) {
			uid, ok := v.Value.(string)
			if ok && strings.HasPrefix(uid, "com.apple.mgmt") {
				return uid, nil
			}
			return "", errors.New("invalid Push Topic (UserID OID) in certificate. Must start with 'com.apple.mgmt', was: " + uid)
		}
	}

	return "", errors.New("could not find Push Topic (UserID OID) in certificate")
}
