package cfg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//Config is a struct containing Server structs
type Config struct {
	Servers []Server `json:"Servers"`
	Admin   []string `json:"Admin"`
}

// Server struct for server structs on load. Convert both Info fields to structs when you learn how to make them variable
type Server struct {
	ID          int               `json:"ID"`
	Game        string            `json:"Game"`
	Name        string            `json:"Name"`
	ServerInfo  map[string]string `json:"ServerInfo"`
	ServiceInfo map[string]string `json:"ServiceInfo"`
	Owner       string            `json:"Owner"`
	Authorized  map[string]bool
}

func LoadConfig() map[int]Server {
	configFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}

	byteValue, _ := ioutil.ReadAll(configFile)

	var config Config

	json.Unmarshal([]byte(byteValue), &config)

	servers := make(map[int]Server)

	for i := 0; i < len(config.Servers); i++ {

		var authorized = make(map[string]bool)

		//Add Owner and all Admin to Authorized
		authorized[config.Servers[i].Owner] = true

		for _, admin := range config.Admin {
			authorized[admin] = true
		}

		config.Servers[i].Authorized = authorized

		servers[config.Servers[i].ID] = config.Servers[i]
	}

	fmt.Println(servers)

	return servers
}
