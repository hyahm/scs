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

	"github.com/hyahm/scs/client/alert"
)

var ErrPnameIsEmpty = errors.New("pname is empty")
var ErrNameIsEmpty = errors.New("name is empty")

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

func (sc *SCSClient) CanNotStop() ([]byte, error) {
	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.Requests("/cannotstop/"+sc.Name, nil)
}

func (sc *SCSClient) CanStop() ([]byte, error) {

	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.Requests("/canstop/"+sc.Name, nil)
}

func (sc *SCSClient) Log() ([]byte, error) {
	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.Requests("/log/"+sc.Name, nil)
}

func (sc *SCSClient) Env() ([]byte, error) {
	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.Requests("/env/"+sc.Name, nil)
}

func (sc *SCSClient) Reload() ([]byte, error) {
	return sc.Requests("/-/reload", nil)
}

func (sc *SCSClient) KillPname() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	return sc.Requests("/kill/"+sc.Pname, nil)
}

func (sc *SCSClient) KillName() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.Requests(fmt.Sprintf("/kill/%s/%s", sc.Pname, sc.Name), nil)
}

func (sc *SCSClient) UpdateAll() ([]byte, error) {
	return sc.Requests("/update", nil)
}

func (sc *SCSClient) UpdatePname() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	return sc.Requests("/update/"+sc.Pname, nil)
}

func (sc *SCSClient) UpdateName() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.Requests(fmt.Sprintf("/update/%s/%s", sc.Pname, sc.Name), nil)
}

func (sc *SCSClient) RestartAll() ([]byte, error) {
	return sc.Requests("/restart", nil)
}

func (sc *SCSClient) RestartPname() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	return sc.Requests("/restart/"+sc.Pname, nil)
}

func (sc *SCSClient) RestartName() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.Requests(fmt.Sprintf("/restart/%s/%s", sc.Pname, sc.Name), nil)
}

func (sc *SCSClient) StartAll() ([]byte, error) {
	return sc.Requests("/start", nil)
}

func (sc *SCSClient) StartPname() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	return sc.Requests("/start/"+sc.Pname, nil)
}

func (sc *SCSClient) StartName() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.Requests(fmt.Sprintf("/start/%s/%s", sc.Pname, sc.Name), nil)
}

func (sc *SCSClient) StopAll() ([]byte, error) {
	return sc.Requests("/stop", nil)
}

func (sc *SCSClient) StopPname() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	return sc.Requests("/stop/"+sc.Pname, nil)
}

func (sc *SCSClient) StopName() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.Requests(fmt.Sprintf("/stop/%s/%s", sc.Pname, sc.Name), nil)
}

func (sc *SCSClient) RemoveAllScrip() ([]byte, error) {
	return sc.Requests("/remove", nil)
}

func (sc *SCSClient) RemovePnameScrip() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	return sc.Requests("/remove/"+sc.Pname, nil)
}

func (sc *SCSClient) RemoveNameScrip() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.Requests(fmt.Sprintf("/remove/%s/%s", sc.Pname, sc.Name), nil)
}

func (sc *SCSClient) Enable() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	return sc.Requests("/enable/"+sc.Pname, nil)
}

func (sc *SCSClient) Disable() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	return sc.Requests("/disable/"+sc.Pname, nil)
}

func (sc *SCSClient) Repo() ([]byte, error) {
	return sc.Requests("/get/repo", nil)
}

func (sc *SCSClient) Search(derivative, serviceName string) ([]byte, error) {
	return sc.Requests(fmt.Sprintf("/search/%s/%s", derivative, serviceName), nil)
}

func (sc *SCSClient) AddScript(s *Script) ([]byte, error) {
	send, _ := json.Marshal(s)
	return sc.Requests("/script", bytes.NewReader(send))
}

func (sc *SCSClient) DelScript(pname string) ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	return sc.Requests("/delete/"+sc.Pname, nil)
}

func (sc *SCSClient) StatusAll() ([]byte, error) {
	return sc.Requests("/status", nil)
}

func (sc *SCSClient) StatusPname() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	return sc.Requests("/status/"+sc.Pname, nil)
}

func (sc *SCSClient) StatusName() ([]byte, error) {
	if sc.Pname == "" {
		return nil, ErrPnameIsEmpty
	}
	if sc.Name == "" {
		return nil, ErrNameIsEmpty
	}
	return sc.Requests(fmt.Sprintf("/status/%s/%s", sc.Pname, sc.Name), nil)
}
func (sc *SCSClient) Probe() ([]byte, error) {
	return sc.Requests("/probe", nil)
}

func (sc *SCSClient) Alert(alert *alert.RespAlert) ([]byte, error) {
	send, _ := json.Marshal(alert)
	return sc.Requests("/set/alert", bytes.NewReader(send))
}
