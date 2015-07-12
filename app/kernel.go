package kernel

import(
	"io/ioutil"
	"encoding/json"
	"runtime"
	"gleipnir/server"
	"gleipnir/errors"
	"strings"
	"os"
	"time"
)

type(
	Configuration struct {
		PathSeparator string
		Gopath string
		Port int `json:"port"`
		ServiceDefinitions ServiceDefinitions `json:"services"`
	}
	Kernel struct {
		UsedPorts []int
		Memory runtime.MemStats
		Server server.Server
		Services map[string]Services
		Configuration Configuration
	}
)

func (k *Kernel) loadConfig() {

	k.Configuration.definePaths()

	data, err := ioutil.ReadFile("config.json")
	errors.Check(err)
	
	err = json.Unmarshal([]byte(data), &k.Configuration)
	errors.Check(err)
}

func (c *Configuration) definePaths() {

	if runtime.GOOS == "windows" {
		c.PathSeparator = "\\"
	} else {
		c.PathSeparator = "/"
	}

	c.Gopath = os.Getenv("GOPATH")

}

func (k *Kernel) Init() {

	k.loadConfig()

}

func (k *Kernel) Run() {

	k.launchServices()
	
	k.Server = server.Server{Port: k.Configuration.Port}
	k.Server.Launch()

	time.Sleep(time.Second * 5)

}

func (k *Kernel) Shutdown() {

	for _, services := range k.Services {

		for _, service := range services {

			err := service.Command.Wait()
			errors.Check(err)

		}
	}

}

func (k *Kernel) launchServices() {

	c := k.Configuration
	k.Services = make(map[string]Services)

	for _, sd := range c.ServiceDefinitions {

		for i := 0; i < sd.MaxInstances; i++ {

			path := c.Gopath + c.PathSeparator + "src" + c.PathSeparator + strings.Replace(sd.Path, ":", c.PathSeparator, -1)

			service := initService(sd, i, path)
			if _, hasName := k.Services[sd.Name]; hasName {
				k.Services[sd.Name] = append(k.Services[sd.Name], service)
			} else {
				k.Services[sd.Name] = Services{service}
			}
			k.UsedPorts = append(k.UsedPorts, service.Port)
			sd.NbInstances++
		}
	}
}