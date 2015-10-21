package app

import(
    "net"
    "github.com/Kern046/GleipnirServer"

)

func (s *Service) connect(message GleipnirServer.Message, conn net.Conn) {

    sendResponse(conn, 200, "Connected")

}