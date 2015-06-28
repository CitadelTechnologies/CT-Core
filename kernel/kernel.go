package kernel

import(
	"fmt"
	"io/ioutil"
	"encoding/json"
	"runtime"
	"gleipnir/server"
	"gleipnir/errors"
	"strings"
	"os"

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
		Services map[string]*Services
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

}

func (k *Kernel) Shutdown() {

	for name, services := range k.Services {

		fmt.Println("Shutdown " + name)

		for _, service := range services {

			procState, err := service.Process.Wait()
			errors.Check(err)

			fmt.Printf("%d status : %d\n", procState.Pid(), procState.Exited())

		}

	}

}

func (k *Kernel) launchServices() {

	c := k.Configuration

	for _, sd := range c.ServiceDefinitions {

		fmt.Printf("Max instances : %v\n", sd.MaxInstances)
		fmt.Printf("gopath : %v\n", c.Gopath)
		fmt.Printf("path separator : %v\n", c.PathSeparator)
		for i := 0; i < sd.MaxInstances; i++ {

			path := c.Gopath + c.PathSeparator + "src" + c.PathSeparator + strings.Replace(sd.Path, ":", c.PathSeparator, -1)
			fmt.Printf("path : %v\n", path)
			fmt.Printf("name : %v\n", sd.Name)
			fmt.Printf("map : %v\n", k.Services[sd.Name])
			os.Exit(1)
			service := initService(sd, i, path)
			k.Services[sd.Name] = new(Services)
			k.Services[sd.Name] = service
			k.UsedPorts = append(k.UsedPorts, service.Port)
			sd.NbInstances++

		}
	}
}