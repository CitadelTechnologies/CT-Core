package app

import(
    "net"
    "github.com/CitadelTechnologies/CT-Client"

)

func (s *Service) connect(message ctclient.Message, conn net.Conn) {

    sendResponse(conn, 200, "Connected")

}