package app

import(
    "os/exec"
    "strconv"
    "github.com/CitadelTechnologies/CT-Client"
    "math/rand"
    "time"
    "net"
    "encoding/json"
)

type(
    ServiceDefinition struct {
        Name string `json:"name"`
        ContainerName string `json:"container_name"`
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
        ConsumedMemory uint64 `json:"consumed_memory"`
        AllocatedMemory uint64 `json:"allocated_memory"`
        StartedAt time.Time `json:"started_at"`
        UpdatedAt time.Time `json:"updated_at"`
        Command *exec.Cmd `json:"-"`
        PreHeating bool `json:"-"`
    }
    Services []*Service
)

func (sd *ServiceDefinition) initService(i int, path string, kernelPort string) Service {

    var s Service
    s.Token = sd.generateToken()
    s.Port = sd.FirstPort + i
    s.PreHeating = sd.PreHeating
    s.Command = exec.Command(path, "--service-port=" + strconv.Itoa(s.Port), "--kernel-port=" + kernelPort, "--token=" + s.Token)
    err := s.Command.Start()
    CheckError(err)

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

func handleServiceMessage(message ctclient.Message, conn net.Conn) {
    service, err := Core.getService(message.Emmitter)
    if err != nil {
        sendResponse(conn, 404, err.Error())
    }
    service.updateStatus(message)
    switch message.Command {
        case "connect": service.connect(message, conn)
        default: sendResponse(conn, 404, "Command not found")
    }
}

func (s *Service) updateStatus(message ctclient.Message) {
    s.ConsumedMemory = message.Status.HeapAlloc
    s.AllocatedMemory = message.Status.HeapSys
    s.StartedAt = message.Status.StartedAt
    s.UpdatedAt = message.Status.UpdatedAt
}

func sendResponse(conn net.Conn, status int, message string) {

    response := ctclient.Response{Status: status, Message: message}
    
    buffer, err := json.Marshal(response)
    CheckError(err)
    if _, err := conn.Write(buffer); err != nil {
        panic(err)
    }
}