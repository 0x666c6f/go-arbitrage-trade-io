package http

import (
	"bytes"
	"crypto/tls"
	"github.com/florianpautot/go-arbitrage-trade-io/model"
	"github.com/florianpautot/go-arbitrage-trade-io/utils"
	"io/ioutil"
	"net/http"
)

//HTTPGet :
func HTTPGet(url string, args string, auth bool) ([]byte, error) {

	req, err := http.NewRequest("GET", url+args, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type","application/json")

	if auth {
		req.Header.Add("Key",model.GlobalConfig.APIKey)
		req.Header.Add("Sign", utils.GenerateSignature(args,model.GlobalConfig.APISecret))
	}

	client := &http.Client{Transport:&http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify : true},}}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return body, nil
}

//HTTPPost :
func HTTPPost(url string, data []byte) ([]byte, error){
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Add("content-type","application/json")

	req.Header.Add("Key",model.GlobalConfig.APIKey)
	req.Header.Add("Sign", utils.GenerateSignature(string(data),model.GlobalConfig.APISecret))

	client := &http.Client{Transport:&http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify : true},}}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()
	return body, nil
}


//HTTPDelete :
func HTTPDelete(url string, args string) ([]byte, error){
	req, err := http.NewRequest("DELETE", url+args, nil)
	req.Header.Add("content-type","application/json")
	req.Header.Add("Key",model.GlobalConfig.APIKey)
	req.Header.Add("Sign", utils.GenerateSignature(args,model.GlobalConfig.APISecret))


	client := &http.Client{Transport:&http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify : true},}}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return body, nil
}
