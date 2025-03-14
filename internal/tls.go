package internal

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io"
	"math/big"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/hyahm/golog"
)

func CreateTLS() {
	fi, err := os.Stat("keys")
	if err != nil {
		os.MkdirAll("keys", 0755)
	} else {
		if !fi.IsDir() {
			panic("exsit file")
		}
	}
	sn := time.Now().Unix()
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(sn),
		Subject: pkix.Name{
			Organization: []string{"scs"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		SubjectKeyId:          []byte{1, 2, 3, 4, 5},
		BasicConstraintsValid: true,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
	privCa, _ := rsa.GenerateKey(rand.Reader, 1024)
	CreateCertificateFile("ca", ca, privCa, ca, nil)
	server := &x509.Certificate{
		SerialNumber: big.NewInt(sn),
		Subject: pkix.Name{
			Organization:       []string{"scs"},
			OrganizationalUnit: []string{"hyahm"},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	hosts := make([]net.IP, 0)
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		golog.Error(err)
		return
	}
	for _, value := range addrs {
		if ipnet, ok := value.(*net.IPNet); ok {
			if ipnet.IP.To4() != nil {
				hosts = append(hosts, ipnet.IP)
			}
		}
	}

	r, err := http.Get("http://ip.hyahm.com")
	if err == nil {
		b, err := io.ReadAll(r.Body)
		if err == nil {
			hosts = append(hosts, net.ParseIP(string(b)))
		}
	}
	server.IPAddresses = append(server.IPAddresses, hosts...)

	privSer, _ := rsa.GenerateKey(rand.Reader, 1024)
	CreateCertificateFile("server", server, privSer, ca, privCa)
	// client := &x509.Certificate{
	// 	SerialNumber: big.NewInt(sn),
	// 	Subject: pkix.Name{
	// 		Organization:       []string{"scs"},
	// 		OrganizationalUnit: []string{"client"},
	// 	},
	// 	NotBefore:    time.Now(),
	// 	NotAfter:     time.Now().AddDate(10, 0, 0),
	// 	SubjectKeyId: []byte{1, 2, 3, 4, 7},
	// 	ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	// 	KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	// }
	// privCli, _ := rsa.GenerateKey(rand.Reader, 1024)
	// CreateCertificateFile("client", client, privCli, ca, privCa)

}

func CreateCertificateFile(name string, cert *x509.Certificate, key *rsa.PrivateKey, caCert *x509.Certificate, caKey *rsa.PrivateKey) {
	name = filepath.Join("keys", name)
	priv := key
	pub := &priv.PublicKey
	privPm := priv
	if caKey != nil {
		privPm = caKey
	}
	ca_b, err := x509.CreateCertificate(rand.Reader, cert, caCert, pub, privPm)
	if err != nil {
		golog.Error("create failed ", err)
		return
	}
	ca_f := name + ".pem"
	var certificate = &pem.Block{Type: "CERTIFICATE",
		Headers: map[string]string{},
		Bytes:   ca_b}
	ca_b64 := pem.EncodeToMemory(certificate)
	os.WriteFile(ca_f, ca_b64, 0600)

	priv_f := name + ".key"
	priv_b := x509.MarshalPKCS1PrivateKey(priv)
	os.WriteFile(priv_f, priv_b, 0600)
	var privateKey = &pem.Block{Type: "PRIVATE KEY",
		Headers: map[string]string{},
		Bytes:   priv_b}
	priv_b64 := pem.EncodeToMemory(privateKey)
	os.WriteFile(priv_f, priv_b64, 0600)
}
