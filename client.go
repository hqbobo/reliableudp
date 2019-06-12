package reliableudp

import (
	"fmt"
	"net"
)

type Client struct {
	port    int
	ip      string
	conn    net.Conn
	svrconn *net.UDPConn
	exit    bool
}

func (this *Client) recv() {
	for {
		if this.exit {
			break
		}
		buf := make([]byte, 65535)

		n, raddr, err := this.svrconn.ReadFromUDP(buf[0:])
		fmt.Print("client:")
		if err != nil {
			fmt.Println("from ReadFromUDP:", err)
		} else {
			fmt.Println(raddr.String(), "-recv:", n)
		}

	}

}

func (this *Client) init() {
	var addr  *net.UDPAddr

	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", this.ip, this.port))
	if err != nil {
		fmt.Println(err)
	}
	this.conn = conn

	if addr, err = net.ResolveUDPAddr("udp", conn.LocalAddr().String()); err != nil {
		fmt.Println(err)
	}

	this.svrconn, err = net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
	}
	this.exit = false
	go this.recv()

}

func (this *Client) Send(msg []byte) (int, error) {
	fmt.Println("send-len : ", len(msg))
	return this.conn.Write(msg)
}

func NewClient(ip string, port int) *Client {
	c := new(Client)
	c.ip = ip
	c.port = port
	c.init()
	return c
}
