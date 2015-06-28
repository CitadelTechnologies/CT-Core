package service

import(
	"os"
	"strings"
	"gleipnir/error"
)

type(
	ServiceDefinition struct {
		Name string `json:"name"`
		Path string `json:"path"`
		FirstPort int `json:"first_port"`
		NbInstances int `json:"nb_instances"`
		MaxInstances int `json:"max_instances"`
	}
	ServiceDefinitions []ServiceDefinition
	Service struct{
		Port int
		Process os.Process
		ProcessAttr os.ProcessAttr
	}
	Services []Service
)

func LaunchServices(k *Kernel) {

	k.usedPorts := make([]int, 0)

	for sd in Config.ServiceDefinitions {

		for i := 0; i < sd.MaxInstances; ++i {

			service := initService(sd)
			k.Services[sd.Name][] = service
			k.usedPorts[] = service.Port
			sd.NbInstances++

		}
	}
}

func initService(sd *ServiceDefinition) *Service {

	path := Gopath + PathSeparator + "src" + PathSeparator + strings.Replace(sd.Path, ":", PathSeparator, -1)

	var s Service

	s.Port = sd.FirstPort + i
	s.Process, err := os.StartProcess("go run " + path, [(string)s.Port], s.ProcessAttr)

	error.Check(err)


	return s
}