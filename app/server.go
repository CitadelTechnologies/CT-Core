package app

import(
    "fmt"
    "net"
    "net/http"
    //"encoding/json"
)

type Server struct {
	HttpPort string `json:"http_port"`
	WsPort string `json:"ws_port"`
	TcpPort string `json:"tcp_port"`
        
        TcpListener net.Listener `json:"-"`
}

func(s *Server) Launch() {

        s.ListenTcp()
        s.ListenHttp()
        s.ListenWebsocket()
}

func(s *Server) ListenTcp() {

    var err error
    if s.TcpListener, err = net.Listen("tcp", ":" + s.TcpPort); err != nil {
        panic(err)
    }
    for {
        var conn net.Conn
        if conn, err = s.TcpListener.Accept(); err != nil {
            panic(err)
        }
        go s.authenticateTcpConnection(conn)
    }

}

func (s *Server) authenticateTcpConnection(conn net.Conn) {

    buffer := make([]byte, 4096)

    if _, err := conn.Read(buffer); err != nil {
        panic(err)
    }
    fmt.Printf("%s\n", buffer)

}

func(s *Server) ListenHttp() {

    http.HandleFunc("/status", sendStatus)
    http.HandleFunc("/run", runKernel)
    http.HandleFunc("/shutdownServices", shutdownServices)
    http.HandleFunc("/shutdown", shutdownKernel)
    http.ListenAndServe(":" + s.HttpPort, nil)
}

func(s *Server) ListenWebsocket() {

    

}

func (s *Server) Shutdown() {


}