package app

import(
	"io/ioutil"
	"encoding/json"
	"runtime"
	"gleipnir/errors"
	"strings"
        "strconv"
	"os"
        "math"
)

type(
	Configuration struct {
                MemoryLimit string `json:"memory_limit"`
		PathSeparator string
		Gopath string
                ServerData Server `json:"server"`
		ServiceDefinitions ServiceDefinitions `json:"services"`
	}
	Kernel struct {
		UsedPorts []int `json:"-"`
                MaxMemory int `json:"max_memory"`
		Memory runtime.MemStats `json:"memory"`

                CpusNumber int `json:"cpus_number"`
                UsedCpus int `json:"used_cpus"`

		Server Server `json:"-"`
		Services map[string]Services `json:"services"`
		Configuration Configuration `json:"-"`

		IsRunning bool `json:"is_running"`
	}
)

var Core Kernel

func init(){

	Core.Initialize()
	defer Core.Shutdown()
	//Core.Run()

}

func (c *Configuration) definePaths() {

	if runtime.GOOS == "windows" {
		c.PathSeparator = "\\"
	} else {
		c.PathSeparator = "/"
	}

	c.Gopath = os.Getenv("GOPATH")

}

func (k *Kernel) Initialize() {

        k.IsRunning = false
        k.refreshProfile()
	k.loadConfig()

        k.CpusNumber = runtime.NumCPU()
        k.UsedCpus = runtime.GOMAXPROCS(0)

	k.Server = k.Configuration.ServerData
	k.Server.Launch()
}

func (k *Kernel) refreshProfile() {

        runtime.ReadMemStats(&k.Memory)
        runtime.GC()

}

/**
 *	Read the config.json file to extract data about services and kernel
 */
func (k *Kernel) loadConfig() {

	k.Configuration.definePaths()

	data, err := ioutil.ReadFile("config.json")
	errors.Check(err)
	
	err = json.Unmarshal([]byte(data), &k.Configuration)
	errors.Check(err)

        var memoryLimit int
        memoryLimit, err = strconv.Atoi(k.Configuration.MemoryLimit[:len(k.Configuration.MemoryLimit)-1])
        errors.Check(err)

        unit := k.Configuration.MemoryLimit[len(k.Configuration.MemoryLimit)-1:]
        units := map[string]int{
            "O": 0,
            "K": 1,
            "M": 2,
            "G": 3,
        }

        k.MaxMemory = memoryLimit * int(math.Pow(1024, float64(units[unit])))
}

func (k *Kernel) Run() {
    
	k.launchServices()
        k.IsRunning = true
        runtime.ReadMemStats(&k.Memory)
}

/**
 *  Get the defined services for the kernel and launch them
 *  In separate processes
 */
func (k *Kernel) launchServices() {
        // The configuration contains the services definitions
	c := k.Configuration
	k.Services = make(map[string]Services)

        // For each one, we instanciate the service in the kernel
	for _, sd := range c.ServiceDefinitions {

                // Each service can be initialized in several processes
                // As long as material resources consumption stays stable (to implement)
		for i := 0; i < sd.MaxInstances; i++ {

                        // The executable file is contained in the config as "project:executable"
			path := c.Gopath + c.PathSeparator + "src" + c.PathSeparator + strings.Replace(sd.Path, ":", c.PathSeparator, -1)

			service := initService(sd, i, path)
                        // If this service has already been initialized, we just append an item to the Services struct
                        // Otherwise we declare a new Services struct with the service inside
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

func (k *Kernel) Shutdown() {

        if k.IsRunning == false {

            return

        }

        k.Server.Shutdown()

	for _, services := range k.Services {

		for _, service := range services {

			err := service.Command.Wait()
			errors.Check(err)

		}
	}

        k.IsRunning = false

}