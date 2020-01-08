package config

import (
	"botBoilerplate/modules"
	"encoding/json"
	"io/ioutil"
)

type Configuration struct {
	Matrix struct {
		Homeserver string `json:"homeserver"`
		Username   string `json:"username"`
		Password   string `json:"password"`
		AuthToken  string `json:"auth_token"`
	} `json:"matrix"`
	AllowedServers []struct {
		Homeserver string `json:"homeserver"`
	} `json:"allowed-servers"`
	PrivilegedUsers    []string                  `json:"privileged-users"`
	DefaultPermissions []modules.PermissionClass `json:"default-permissions"`
	Database           string                    `json:"database"`
	Ansible            struct {
		Inventory string `json:"inventory"`
	} `json:"ansible"`
	Iptables struct {
		IpNets []struct {
			Name string `json:"name"`
			Net  string `json:"net"`
		} `json:"ip-nets"`
		Domains []string `json:"domains"`
	} `json:"iptables"`
}

var Config Configuration

func Load(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &Config)
}

func Save(path string) error {
	data, err := json.MarshalIndent(Config, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0600)
}
