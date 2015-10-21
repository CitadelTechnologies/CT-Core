package app

import(
    "io/ioutil"
    "encoding/json"
    "runtime"
    "strings"
    "strconv"
    "os"
    "math"
    "sync"
    "errors"
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
        ServerWaitGroup sync.WaitGroup `json:"-"`

        Services map[string]Services `json:"services"`
        ServiceTokens map[string]*Service `json:"-"`

        Configuration Configuration `json:"-"`

        IsRunning bool `json:"is_running"`
    }
)

var Core Kernel

func init(){

	Core.Initialize()
        Core.ServerWaitGroup.Wait()
	defer Core.Shutdown()

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
	k.Services = make(map[string]Services)
	k.ServiceTokens = make(map[string]*Service)
        k.refreshProfile()
	k.loadConfig()

        k.CpusNumber = runtime.NumCPU()
        k.UsedCpus = runtime.GOMAXPROCS(0)

	k.launchServices(true)

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
	CheckError(err)
	
	err = json.Unmarshal([]byte(data), &k.Configuration)
	CheckError(err)

        var memoryLimit int
        memoryLimit, err = strconv.Atoi(k.Configuration.MemoryLimit[:len(k.Configuration.MemoryLimit)-1])
        CheckError(err)

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
    
	k.launchServices(false)
        k.IsRunning = true
        runtime.ReadMemStats(&k.Memory)
}

/**
 *  Get the defined services for the kernel and launch them
 *  In separate processes
 *  If preHeating is set to true, only the configured services will be launched
 */
func (k *Kernel) launchServices(preHeating bool) {
        // The configuration contains the services definitions
	c := k.Configuration

        // For each one, we instanciate the service in the kernel
	for _, sd := range c.ServiceDefinitions {

            if preHeating == true && sd.PreHeating == false  {
                continue
            } else if _, serviceExists := k.Services[sd.Name]; serviceExists {
                continue
            }

            // Each service can be initialized in several processes
            // As long as material resources consumption stays stable (to implement)
            for i := 0; i < sd.MaxInstances; i++ {

                    // The executable file is contained in the config as "project:executable"
                    path := c.Gopath + c.PathSeparator + "src" + c.PathSeparator + strings.Replace(sd.Path, ":", c.PathSeparator, -1)

                    service := sd.initService(i, path, c.ServerData.TcpPort)
                    // If this service has already been initialized, we just append an item to the Services struct
                    // Otherwise we declare a new Services struct with the service inside
                    if _, hasName := k.Services[sd.Name]; hasName {
                            k.Services[sd.Name] = append(k.Services[sd.Name], &service)
                    } else {
                            k.Services[sd.Name] = Services{&service}
                    }
                    k.ServiceTokens[service.Token] = &service
                    k.UsedPorts = append(k.UsedPorts, service.Port)
                    sd.NbInstances++
            }
	}
}

func (k *Kernel) Shutdown() {

    k.Server.Shutdown()
    k.ShutdownServices(true)
    os.Exit(0)

}

/**
 * If preheating is set to false, the preheating services won't be extinguished
 */
func (k *Kernel) ShutdownServices(preheating bool) {

    for key, services := range k.Services {
        mustBeDeleted := true

        for _, service := range services {
            if preheating == false && service.PreHeating == true {
                mustBeDeleted = false
                break
            }
            err := service.Command.Process.Kill()
            CheckError(err)
        }
        if mustBeDeleted == true {
            delete(k.Services, key)
        }
    }
    k.IsRunning = false
}

func (k *Kernel) getService(token string) (*Service, error) {

    pointer := k.ServiceTokens[token]
    if pointer == nil {
        return nil, errors.New("Service not found")
    }
    return pointer, nil
}