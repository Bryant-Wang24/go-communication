package model

import (
	"client/message"
	"net"
)

type CurUser struct {
	Conn net.Conn
	message.User
}
