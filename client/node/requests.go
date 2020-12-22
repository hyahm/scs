package node

import (
	"crypto/tls"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var ReadTimeout time.Duration

func client(timeout time.Duration) *http.Client {
	// var tr *http.Transport
	// certs, err := tls.LoadX509KeyPair(rootCa, rootKey)
	// if err != nil {
	// 	tr = &http.Transport{
	// 		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// 	}
	// } else {
	// 	ca, err := x509.ParseCertificate(certs.Certificate[0])
	// 	if err != nil {
	// 		return &http.Client{Transport: tr}
	// 	}
	// 	pool := x509.NewCertPool()
	// 	pool.AddCert(ca)

	// 	tr = &http.Transport{
	// 		TLSClientConfig: &tls.Config{RootCAs: pool},
	// 	}

	// }
	// return &http.Client{Transport: tr}

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: timeout,
	}

}

func Requests(method, url, token string, body io.Reader, timeout ...time.Duration) ([]byte, error) {

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Token", token)
	var t time.Duration
	if len(timeout) == 0 {
		t = 0
	} else {
		t = timeout[0]
	}
	resp, err := client(t).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 203 {
		return nil, errors.New("token error, you can use scsctl config token <token> set server token")
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(string(b))
	}
	return b, nil

}
