package main

import(
    "fmt"
    "flag"
    "runtime"
    "gleipnir/node"
)

type Server struct{
    ip string
    port int
    memoryStats runtime.MemStats
    nbCores int
	nodes map[int]node.Node
}

func main(){

	var server Server
    server.initialize()
	
    fmt.Printf("Server is ready at %s:%d\n",  server.ip, server.port)
	
	server.launchNodes()

}

func (s *Server) initialize(){

	s.ip = *flag.String("ip", "127.0.0.1", "The Server IP")
    s.port = *flag.Int("port", 9999, "The Server Port")
    flag.Parse()

	s.nodes = make(map[int]node.Node)
    s.nbCores = runtime.NumCPU()
    runtime.ReadMemStats(&s.memoryStats)

    fmt.Printf("Server is now initialized with %d/%d bytes \n", s.memoryStats.Alloc, s.memoryStats.TotalAlloc)
}

func (s *Server) launchNodes(){

    fmt.Printf("%d cores are available on this machine\n", s.nbCores)
	
	s.createNode()

}

func (s *Server) createNode(){

    var newNode node.Node
	newNode.Id = s.generateNodeId()
    newNode.Clients = make(map[int]node.Client, 0)
	
    s.nodes[newNode.Id] = newNode

	fmt.Printf("The node %d is now initialized\n", newNode.Id)
	
}

func (s *Server) generateNodeId() int {

    for i := 1; ; i++ {
	
	    if _, ok := s.nodes[i]; !ok {
		
		    return i
		
		}
	
	}
}