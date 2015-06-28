package kernel

import(
	"io/ioutil"
	"encoding/json"
	"runtime"
	"gleipnir/server"
	"gleipnir/service"
	"gleipnir/error"

)

type(
	Configuration struct {
		Port int `json:"port"`
		ServiceDefinitions service.ServiceDefinitions `json:"services"`
	}
	Kernel struct {
		Memory runtime.MemStats
		Server server.Server
		Services Services
		Configuration Configuration
	}
)

func (k *Kernel) loadConfig() {
	data, err := ioutil.ReadFile("config.json")
	error.Check(err)
	
	err = json.Unmarshal([]byte(data), &k.Configuration)
	error.Check(err)
}

func (k *Kernel) Init() {

	k.loadConfig()

}

func (k *Kernel) Run() {

	service.LaunchServices()
	
	k.Server := Server{Port: k.Configuration.Port}
	k.Server.Launch()

}