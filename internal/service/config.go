package service

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type config struct {
	Duration uint   `json:"Duration"`
	OffLevel int    `json:"OffLevel"`
	OnLevel  int    `json:"OnLevel"`
	Comments string `json:"_comments"`
}

// creates service configuration file in ~/.smartlight/
func createConfig(c config) error {
	b, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}
	d, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	d += "/.smartlight"
	_, err = os.Stat(d)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		err = os.Mkdir(d, 0755)
		if err != nil {
			return err
		}
	}
	d += "/service.json"
	err = ioutil.WriteFile(d, b, 0644)
	if err != nil {
		return err
	}
	return nil
}

// loads service configuration file from ~/.smartlight/
func loadConfig() (*config, error) {
	d, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	d += "/.smartlight/service.json"
	_, err = os.Stat(d)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(d)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	c := &config{}
	err = json.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
