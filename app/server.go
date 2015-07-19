package app

import(
    "net"
    "net/http"
    "encoding/json"
    "github.com/Kern046/GleipnirServer"
    //"unsafe"
    "bytes"
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

    var message GleipnirServer.Message
    buffer := make([]byte, 2048)

    if _, err := conn.Read(buffer); err != nil {
        panic(err)
    }
    buffer = bytes.Trim(buffer, "\x00")
    if err := json.Unmarshal(buffer, &message); err != nil {
        panic(err)
    }
    handleServiceMessage(message)
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