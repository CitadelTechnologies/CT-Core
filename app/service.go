package app

import(
    "os/exec"
    "strconv"
    "gleipnir/errors"
    "github.com/Kern046/GleipnirServer"
    "math/rand"
    "time"
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
            Token string `json:"-"`
            Port int `json:"port"`
            Command *exec.Cmd `json:"-"`
	}
	Services []Service
)

func (sd *ServiceDefinition) initService(i int, path string, kernelPort string) Service {

    var s Service
    s.Token = sd.generateToken()
    s.Port = sd.FirstPort + i
    s.Command = exec.Command(path, "--service-port=" + strconv.Itoa(s.Port), "--kernel-port=" + kernelPort)

    err := s.Command.Start()
    errors.Check(err)

    return s
}

func (sd *ServiceDefinition) generateToken() string{

    for {
        token := getRandomString(25)

        if _, exists := Core.ServiceTokens[token]; exists != true {
            return token
        }
    }
}

func getRandomString(strlen int) string {
    rand.Seed(time.Now().UTC().UnixNano())
    const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
    result := make([]byte, strlen)
    for i := 0; i < strlen; i++ {
            result[i] = chars[rand.Intn(len(chars))]
    }
    return string(result)
}

func handleServiceMessage(message GleipnirServer.Message) {
    
    service := Core.getService(message.Emmitter)

    switch message.Command {
        case "connect": 
    }
}