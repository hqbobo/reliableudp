package reliableudp

import (
	"fmt"
	"net"
)

type Server struct {
	port int
	ip string
	local  *net.UDPAddr
	conn   *net.UDPConn
	exit bool
}

func (this *Server) recv() {
	for {
		if this.exit {
			break
		}
		buf := make([]byte, 65535)

		n, raddr, err := this.conn.ReadFromUDP(buf[0:])
		if err != nil {
			fmt.Println("from ReadFromUDP:", err)
		} else {
			fmt.Println(raddr.String(),"-recv:", n)
		}

	}

}

func (this *Server) init (){
	var err error
	if this.local, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", this.ip, this.port)); err != nil {
		fmt.Println(err)
	}
	this.conn, err = net.ListenUDP("udp", this.local)
	if err != nil {
		fmt.Println(err)
	}
	this.exit = false
	go this.recv()
}


func NewServer(ip string, port int) *Server{
	s := new(Server)
	s.ip = ip
	s.port = port
	s.init()
	return s
}

