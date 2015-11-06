

package Utils

import (
	"net/http"
	"io/ioutil"
	"net"
	"time"
)


func RequestUrl(url string) (string, error) {

	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*50)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * 50))
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * 30,
		},
	}
	rsp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil

}