package kernel

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
	}
	ServiceDefinitions []ServiceDefinition
	Service struct{
		Port int
		Command *exec.Cmd
	}
	Services []Service
)

func initService(sd ServiceDefinition, i int, path string) Service {

	var s Service
	s.Port = sd.FirstPort + i
	s.Command = exec.Command(path, strconv.Itoa(s.Port))

	err := s.Command.Start()
	errors.Check(err)

	return s
}