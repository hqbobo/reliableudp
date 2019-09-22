package reliableudp

import (
	"fmt"
	"net"
)

type Server struct {
	port  int
	ip    string
	local *net.UDPAddr
	conn  *net.UDPConn
	exit  bool
	mng   sessionMng
}

type session struct {
	addr net.Addr
	conn net.Conn
	buf  map[uint32]*protobuffer
}

func (this *session) Add(data []byte, reader udpreader) *proto {
	p := parseProto(data)
	if val, ok := this.buf[p.uni]; ok {
		val.addProto(p)
		if val.isfull() {
			if reader != nil {
				 go reader.OnRecive(val.data(), this.addr)
			}
			delete(this.buf,p.uni)

		}
	} else {
		this.buf[p.uni] = newProtobuffer()
		this.buf[p.uni].addProto(p)
	}
	return p
}

func newSession(addr net.Addr) *session {
	s := new(session)
	s.addr = addr
	s.buf = make(map[uint32]*protobuffer, 0)
	conn, err := net.Dial("udp", addr.String())
	if err != nil {
		glog.Warn(err)
	}
	s.conn = conn
	return s
}

type sessionMng struct {
	sessions map[string]*session
	reader   udpreader
}

func (this *sessionMng) Add(addr net.Addr, data []byte) {
	if val, ok := this.sessions[addr.String()]; ok {
		p := val.Add(data, this.reader)
		if p.types == types_trans {
			glog.Debug("send ",p.seq,"-", p.uni)
			this.sendAck(addr.(*net.UDPAddr) , newSeqAck(p.seq, p.uni).ToByte())
		}
	} else {
		sess := newSession(addr)
		this.sessions[addr.String()] = sess
		sess.Add(data, this.reader)
	}
}

func (this *sessionMng) sendAck(addr *net.UDPAddr, data []byte){
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", addr.IP, addr.Port +1))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	_ , e := conn.Write(data)
	if e != nil {
		glog.Warn(e)
	}
}


func (this *Server) recv() {
	for {
		if this.exit {
			break
		}
		buf := make([]byte, udp_datalen+12)
		n, raddr, err := this.conn.ReadFromUDP(buf[0:])
		if err != nil {
			glog.Warn("from ReadFromUDP:", err)
		} else {
			glog.Debug(raddr.String(), "-recv:", n)
			this.mng.Add(raddr, buf[0:n])
		}
	}
}
func (this *Server) init() {
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

func NewServer(ip string, port int, reader udpreader) *Server {
	s := new(Server)
	s.ip = ip
	s.port = port
	s.mng.sessions = make(map[string]*session)
	s.mng.reader = reader
	s.init()
	return s
}
