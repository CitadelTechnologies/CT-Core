package kernel

import(
	"os"
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
		Process *os.Process
		ProcessAttr *os.ProcAttr
	}
	Services []Service
)

func initService(sd ServiceDefinition, i int, path string) Service {

	var s Service
	var args []string
	var err error

	s.ProcessAttr = new(os.ProcAttr)
	s.Port = sd.FirstPort + i

	args = append(args, strconv.Itoa(s.Port))

	s.Process, err = os.StartProcess("go run " + path, args, s.ProcessAttr)

	errors.Check(err)

	return s
}