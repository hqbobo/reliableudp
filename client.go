package reliableudp

import (
	"fmt"
	"net"
)
type Client struct {
	port int
	ip string
	conn net.Conn
}

func (this *Client) init (){
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", this.ip, this.port))
	if err != nil {
		fmt.Println(err)
	}
	this.conn = conn
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
