package utils

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"github.com/florianpautot/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

//LoadConfig :
func LoadConfig(path string) (model.Config,error){

	config := model.Config{}

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return model.Config{},err
	}
	err = yaml.UnmarshalStrict(yamlFile, &config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return model.Config{},err

	}

	return config,nil
}

//GenerateSignature :
func GenerateSignature(input string, secret string) string {
	hmac512 := hmac.New(sha512.New, []byte(secret))
	hmac512.Write([]byte(input))
	signature := hex.EncodeToString(hmac512.Sum(nil))
	return signature;
}
