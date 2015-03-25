package main

import(
    "fmt"
    "flag"
    "runtime"
    //"gleipnir/node"
)

type Server struct{
    ip string
    port int
    memoryStats runtime.MemStats
    nbCores int
}

func main(){

    server := Initialize()

    fmt.Printf("Server is ready at %s:%d\n",  server.ip, server.port)

}

func Initialize() *Server{

    serverIp := flag.String("ip", "127.0.0.1", "The Server IP")
    serverPort := flag.Int("port", 9999, "The Server Port")

    flag.Parse()

    server := Server{
        ip: *serverIp,
        port: *serverPort,
    }

    server.nbCores = runtime.NumCPU()
    runtime.ReadMemStats(&server.memoryStats)

    fmt.Printf("Server is now initialized with %d/%d bytes \n", server.memoryStats.Alloc, server.memoryStats.TotalAlloc)

    return &server
}
