package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/hyahm/scs/alert"
	"github.com/hyahm/scs/internal"
)

type SCSClient struct {
	Domain string
	Token  string
	Pname  string
	Name   string
}

func NewClient() *SCSClient {
	return &SCSClient{
		Domain: "https://127.0.0.1:11111",
		Token:  os.Getenv("TOKEN"),
		Pname:  os.Getenv("PNAME"),
		Name:   os.Getenv("NAME"),
	}
}

func client() *http.Client {
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
		Timeout: 5 * time.Second,
	}

}

func (sc *SCSClient) Requests(url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, sc.Domain+url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Token", sc.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 203 {
		return nil, errors.New("token error")
	}
	return ioutil.ReadAll(resp.Body)
}

func (sc *SCSClient) CanNotStop(name ...string) ([]byte, error) {
	temp := sc.Name
	if len(name) > 0 {
		temp = name[0]
	}
	return sc.Requests("/cannotstop/"+temp, nil)
}

func (sc *SCSClient) CanStop(name ...string) ([]byte, error) {
	temp := sc.Name
	if len(name) > 0 {
		temp = name[0]
	}
	return sc.Requests("/canstop/"+temp, nil)
}

func (sc *SCSClient) Log(name string) ([]byte, error) {
	return sc.Requests("/log/"+name, nil)
}

func (sc *SCSClient) Env(name string) ([]byte, error) {
	return sc.Requests("/env/"+name, nil)
}

func (sc *SCSClient) Reload() ([]byte, error) {
	return sc.Requests("/-/reload", nil)
}

func (sc *SCSClient) Kill(args ...string) ([]byte, error) {
	l := len(args)
	switch l {
	case 1:
		return sc.Requests("/kill/"+args[0], nil)
	default:
		return sc.Requests(fmt.Sprintf("/kill/%s/%s", args[0], args[1]), nil)
	}
}

func (sc *SCSClient) Update(args ...string) ([]byte, error) {
	l := len(args)
	switch l {
	case 0:
		return sc.Requests("/update", nil)
	case 1:
		return sc.Requests("/update/"+args[0], nil)
	default:
		return sc.Requests(fmt.Sprintf("/update/%s/%s", args[0], args[1]), nil)
	}
}

func (sc *SCSClient) Restart(args ...string) ([]byte, error) {
	l := len(args)
	switch l {
	case 0:
		return sc.Requests("/restart", nil)
	case 1:
		return sc.Requests("/restart/"+args[0], nil)
	default:
		return sc.Requests(fmt.Sprintf("/restart/%s/%s", args[0], args[1]), nil)
	}
}

func (sc *SCSClient) Start(args ...string) ([]byte, error) {
	l := len(args)
	switch l {
	case 0:
		return sc.Requests("/start", nil)
	case 1:
		return sc.Requests("/start/"+args[0], nil)
	default:
		return sc.Requests(fmt.Sprintf("/start/%s/%s", args[0], args[1]), nil)
	}
}

func (sc *SCSClient) Stop(args ...string) ([]byte, error) {
	l := len(args)
	switch l {
	case 0:
		return sc.Requests("/stop", nil)
	case 1:
		return sc.Requests("/stop/"+args[0], nil)
	default:
		return sc.Requests(fmt.Sprintf("/stop/%s/%s", args[0], args[1]), nil)
	}
}

func (sc *SCSClient) Remove(args ...string) ([]byte, error) {
	l := len(args)
	switch l {
	case 0:
		return sc.Requests("/remove", nil)
	case 1:
		return sc.Requests("/remove/"+args[0], nil)
	default:
		return sc.Requests(fmt.Sprintf("/remove/%s/%s", args[0], args[1]), nil)
	}
}

func (sc *SCSClient) Repo() ([]byte, error) {
	return sc.Requests("/get/repo", nil)
}

func (sc *SCSClient) Search(derivative, serviceName string) ([]byte, error) {
	return sc.Requests(fmt.Sprintf("/search/%s/%s", derivative, serviceName), nil)
}

func (sc *SCSClient) Script(s *internal.Script) ([]byte, error) {
	send, _ := json.Marshal(s)
	return sc.Requests("/script", bytes.NewReader(send))
}

func (sc *SCSClient) DelScript(pname string) ([]byte, error) {
	return sc.Requests("/delete/"+pname, nil)
}

func (sc *SCSClient) Status(args ...string) ([]byte, error) {
	l := len(args)
	switch l {
	case 0:
		return sc.Requests("/status", nil)
	case 1:
		return sc.Requests("/status/"+args[0], nil)
	default:
		return sc.Requests(fmt.Sprintf("/status/%s/%s", args[0], args[1]), nil)
	}
}

func (sc *SCSClient) Probe() ([]byte, error) {
	return sc.Requests("/probe", nil)
}

func (sc *SCSClient) Alert(alert *alert.RespAlert) ([]byte, error) {
	send, _ := json.Marshal(alert)
	return sc.Requests("/set/alert", bytes.NewReader(send))
}
