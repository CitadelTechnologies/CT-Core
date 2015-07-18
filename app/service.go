package app

import(
	"os/exec"
	"strconv"
	"gleipnir/errors"
)

type(
	ServiceDefinition struct {
		Name string `json:"name"`
		Path string `json:"path"`
		FirstPort int `json:"first_port"`
		NbInstances int `json:"nb_instances"`
		MaxInstances int `json:"max_instances"`
                PreHeating bool `json:"preheating"`
	}
	ServiceDefinitions []ServiceDefinition
	Service struct{
		Port int `json:"port"`
		Command *exec.Cmd `json:"-"`
	}
	Services []Service
)

func initService(sd ServiceDefinition, i int, path string, kernelPort string) Service {

	var s Service
	s.Port = sd.FirstPort + i
	s.Command = exec.Command(path, "--service-port=" + strconv.Itoa(s.Port), "--kernel-port=" + kernelPort)

	err := s.Command.Start()
	errors.Check(err)

	return s
}