package main

import(
    "fmt"
    "flag"
    "gleipnir/node"
)

type Server struct{
    ip string
    port int
}

func main(){

    serverIp := flag.String("ip", "127.0.0.1", "The Server IP")
    serverPort := flag.Int("port", 9999, "The Server Port")

    flag.Parse()

    server := Server{
        ip: *serverIp,
        port: *serverPort,
    }

    c := make(chan string)

    go node.Initialize(c)

    msg := <- c

    fmt.Printf(msg + "\n")
    fmt.Printf("Server initialized at " + server.ip + ":" + fmt.Sprintf("%d", server.port) + "\n")
}
