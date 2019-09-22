package reliableudp

import (
	"fmt"
	"net"
	"log"
)

type Client struct {
	port    int
	ip      string
	conn    net.Conn
	svrconn *net.UDPConn
	exit    bool
	uni     map[uint32] *protobuffer
}

func (this *Client) recv() {
	for {
		if this.exit {
			break
		}
		buf := make([]byte, 65535)

		n, raddr, err := this.svrconn.ReadFromUDP(buf[0:])
		if err != nil {
			glog.Warn("from ReadFromUDP:", raddr.String(), err)
		} else {
			p := parseProto(buf[0:n])
			switch p.types {
			case types_transack:
				glog.Debug("Ack for seq ", p.seq, " uin ", p.uni)
				if v, ok := this.uni[p.uni]; ok {
					v.ack(p.seq)
				}
			case types_wanted:
			}
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
	if addr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%d", this.conn.LocalAddr().(*net.UDPAddr).Port + 1)); err != nil {
		fmt.Println(err)
	} else {
		glog.Debug("Client recvie on ", addr.String())
	}
	this.svrconn, err = net.ListenUDP("udp", addr)
	if err != nil {
		log.Panic(err)
	}
	this.exit = false
	go this.recv()

}

func (this *Client) Send(msg []byte) (sndlen int, err error) {
	glog.Debug("send-len : ", len(msg))
	p := newProtobuffer()
	pinit:
	uni := p.init(msg)
	if _, ok := this.uni[uni]; !ok {
		this.uni[uni] = p
	} else {
		glog.Warn("uin ", uni, " is in use")
		goto pinit
	}
	//return this.conn.Write(msg)
	p.traversal(func(data []byte) {
		len , e := this.conn.Write(data)
		if e != nil {
			glog.Warn(e)
		}
		sndlen+=len
	})
	<-p.done
	return sndlen, err
}

func NewClient(ip string, port int) *Client {
	c := new(Client)
	c.ip = ip
	c.port = port
	c.uni  = make(map[uint32] *protobuffer,0)
	c.init()
	return c
}
