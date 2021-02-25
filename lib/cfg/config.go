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
	ID          int            `json:"ID"`
	Game        string         `json:"Game"`
	Name        string         `json:"Name"`
	ServerInfo  ServerInfo     `json:"ServerInfo"`
	Commands    Commands       `json:"Commands"`
	ServiceInfo DynamicService `json:"ServiceInfo"`
	Owner       string         `json:"Owner"`
	Authorized  map[string]bool
}

type ServerInfo struct {
	Password string `json:"Password"`
	Port     string `json:"Port"`
}

type Commands struct {
	BackupToS3 CommandBackupToS3 `json:"BackupToS3"`
}

type CommandBackupToS3 struct {
	Source      string `json:"Source"`
	Destination string `json:"Destination"`
}

//DynamicService | Dynmaically pick struct for which service is in JSON
type DynamicService struct {
	Service Service
}

type Service interface {
	Name() string
	Instance() string
	Account() string
	Geolocation() string
}

func (d *DynamicService) UnmarshalJSON(data []byte) error {
	var service struct {
		Service string `json:"Service"`
	}
	if err := json.Unmarshal(data, &service); err != nil {
		return err
	}
	switch service.Service {
	case "aws":
		d.Service = new(ServiceAWS)
	case "azure":
		d.Service = new(ServiceAzure)
	case "gcp":
		d.Service = new(ServiceGCP)
	}
	return json.Unmarshal(data, d.Service)

}

//ServiceAWS | AWS Service struct for unpacking cofig file
type ServiceAWS struct {
	Service    string `json:"Service"`
	AccountID  string `json:"AccountID"`
	Region     string `json:"Region"`
	InstanceID string `json:"InstanceID"`
}

func (s ServiceAWS) Name() string {
	return s.Service
}

func (s ServiceAWS) Instance() string {
	return s.Service
}

func (s ServiceAWS) Account() string {
	return s.Service
}

func (s ServiceAWS) Geolocation() string {
	return s.Service
}

// Not implemented yet
type ServiceAzure struct {
	Service        string `json:"Service"`
	SubscriptionID string `json:"SubscriptionID"`
	Location       string `json:"Location"`
	VmName         string `json:"VmName"`
}

func (s ServiceAzure) Name() string {
	return s.Service
}

func (s ServiceAzure) Instance() string {
	return s.VmName
}

func (s ServiceAzure) Account() string {
	return s.SubscriptionID
}

func (s ServiceAzure) Geolocation() string {
	return s.Location
}

// Not implemented yet
type ServiceGCP struct {
	Service      string `json:"Service"`
	Project      string `json:"Project"`
	Zone         string `json:"Zone"`
	InstanceName string `json:"InstanceName"`
}

func (s ServiceGCP) Name() string {
	return s.Service
}

func (s ServiceGCP) Instance() string {
	return s.InstanceName
}

func (s ServiceGCP) Account() string {
	return s.Project
}

func (s ServiceGCP) Geolocation() string {
	return s.Zone
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
