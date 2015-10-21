package app

import(
    "net"
    "net/http"
    "encoding/json"
    "github.com/Kern046/GleipnirServer"
    "github.com/gorilla/websocket"
    "bytes"
)

type(
    Server struct {
        HttpPort string `json:"http_port"`

        TcpListener net.Listener `json:"-"`
        TcpPort string `json:"tcp_port"`

        WsConn WsConnection `json:"-"`
        WsPort string `json:"ws_port"`
    }
    WsConnection  struct {
        // The websocket connection.
        Connection *websocket.Conn
        // Buffered channel of outbound messages.
        OutputBuffer chan []byte
        // The hub.
        Hub *WsHub `json:"-"`
    }
    WsHandler struct {
        Hub *WsHub
    }
    WsHub struct {
        // Registered connections.
        Connections map[*WsConnection]bool
        // Inbound messages from the connections.
        Broadcast chan []byte
        // Register requests from the connections.
        Register chan *WsConnection
        // Unregister requests from connections.
        Unregister chan *WsConnection
    }
)

func newHub() *WsHub {
    return &WsHub{
        Broadcast:   make(chan []byte),
        Register:    make(chan *WsConnection),
        Unregister:  make(chan *WsConnection),
        Connections: make(map[*WsConnection]bool),
    }
}

func(s *Server) Launch() {

    Core.ServerWaitGroup.Add(3)
    go func() { defer Core.ServerWaitGroup.Done(); s.ListenTcp()}()
    go func() { defer Core.ServerWaitGroup.Done(); s.ListenHttp()}()
    go func() { defer Core.ServerWaitGroup.Done(); s.ListenWebsocket()}()
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
    handleServiceMessage(message, conn)
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

func (c *WsConnection) Read() {
    for {
        _, message, err := c.Connection.ReadMessage()
        if err != nil {
            break
        }
        c.Hub.Broadcast <- message
    }
    c.Connection.Close()
}

func (c *WsConnection) Write() {
    for message := range c.OutputBuffer {
        err := c.Connection.WriteMessage(websocket.TextMessage, message)
        if err != nil {
            break
        }
    }
    c.Connection.Close()
}

func (wsh WsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    upgrader := &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }
    wsc := &WsConnection{OutputBuffer: make(chan []byte, 256), Connection: ws, Hub: wsh.Hub}
    wsc.Hub.Register <- wsc
    defer func() { wsc.Hub.Unregister <- wsc }()
    go wsc.Write()
    wsc.Read()
}

func (s *Server) Shutdown() {


}