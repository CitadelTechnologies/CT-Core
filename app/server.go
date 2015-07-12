package app

import(
    "net/http"
)

type Server struct {
	HttpPort string `json:"http_port"`
	WsPort string `json:"ws_port"`
	TcpPort string `json:"tcp_port"`
}

func(s *Server) Launch() {

        s.ListenTcp()
        s.ListenHttp()
        s.ListenWebsocket()
}

func(s *Server) ListenTcp() {

    

}

func(s *Server) ListenHttp() {

    http.HandleFunc("/status", sendStatus)
    http.ListenAndServe(":" + s.HttpPort, nil)
}

func(s *Server) ListenWebsocket() {

    

}

func (s *Server) Shutdown() {


}