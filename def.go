package reliableudp

import (
	"encoding/binary"
	"math/rand"
)

const (
	types_start    = iota
	types_init     //init packet which carry datalen
	types_trans    //data transimition
	types_transack //data acknowlege
	types_wanted   //client miss some packet
	types_end      //last packet
)

//
//
//
//	| byte  | uni   | byte  |  data     |
//      | types | uni   |  seq  |  bytes    |
//
// in case of types_init
//	| byte  | byte  | byte  |  data(4 bytes) |
//      | 1     | uni   |  0    |  totallen      |
// in case of types_trans
//	| byte  | byte  | byte  |  data     |
//      | 2     | uni   | seqid |  data...  |
// in case of types_end
//	| byte  | byte  | byte  |  data(4 bytes) |
//      | 3     | uni   | seqid |  totallen      |

const udp_datalen = 1400
type proto struct {
	types uint32
	uni   uint32
	seq   uint32
	data  []byte
}

func (this proto)ToByte() []byte{
	buf := make([]byte, 12 + len(this.data))
	binary.BigEndian.PutUint32(buf[0:4], this.types)
	binary.BigEndian.PutUint32(buf[4:8], this.uni)
	binary.BigEndian.PutUint32(buf[8:12], this.seq)
	copy(buf[12:] ,this.data)
	return buf
}

func newProto(data[]byte) *proto{
	p := new(proto)
	p.types = binary.BigEndian.Uint32(data[0:4])
	p.uni = binary.BigEndian.Uint32(data[4:8])
	p.seq = binary.BigEndian.Uint32(data[8:12])
	p.data = data[12:]
	return p
}

func parseProto(data[]byte) *proto{
	p := new(proto)
	p.types = binary.BigEndian.Uint32(data[0:4])
	p.uni = binary.BigEndian.Uint32(data[4:8])
	p.seq = binary.BigEndian.Uint32(data[8:12])
	p.data = data[12:]
	return p
}


type protobuffer struct {
	i     proto
	datas map[uint32]*proto
	e     proto
}

func (this *protobuffer) init(data []byte) {
	uni := rand.Int()
	datalen := len(data)
	this.i.types = types_init
	this.i.seq = 0
	this.i.uni = uint32(uni)
	this.i.data = make([]byte,4)
	binary.BigEndian.PutUint32(this.i.data, uint32(datalen))
	i := 0
	sendlen := 0
	for sendlen < datalen{
		protolens := 0
		if datalen >= sendlen+udp_datalen {
			protolens = udp_datalen
		} else {
			protolens = datalen - sendlen
		}
		p := new(proto)
		p.types = uint32(types_trans)
		p.uni = uint32(uni)
		p.seq = uint32(i)
		p.data = data[sendlen: sendlen+protolens]
		this.datas[uint32(i)] = p
		sendlen += protolens
		i++
	}
	this.e.types = types_end
	this.e.seq = 0
	this.e.uni = uint32(uni)
	this.e.data = make([]byte,4)
	binary.BigEndian.PutUint32(this.e.data, uint32(datalen))
	glog.Debug("uni:", uni)
}


func (this *protobuffer) addProto (p *proto) {
	switch p.types {
	case types_init:
		this.i = *p
	case types_end:
		this.e = *p
	case types_trans:
		this.datas[p.seq] = p
	}
}

func (this *protobuffer) traversal (f func(data[]byte)) {
	f(this.i.ToByte())
	for _ , v :=range this.datas{
		f(v.ToByte())
	}
	f(this.e.ToByte())
}

func (this *protobuffer) data () []byte {
	totallen := binary.BigEndian.Uint32(this.i.data)
	buffer := make([]byte, int(totallen))
	for _ , v :=range this.datas{
		copy(buffer[v.seq * udp_datalen:], v.data)
	}
	return buffer
}

func (this *protobuffer) isfull () bool {
	if this.i.types == 0  || this.e.types == 0  {
		return false
	}
	totallen := binary.BigEndian.Uint32(this.i.data)
	expect := totallen / udp_datalen
	if totallen % udp_datalen != 0 { expect++ }
	if len(this.datas) != int(expect) {
		return false
	}
	return true
}

func newProtobuffer() *protobuffer {
	b := new(protobuffer)
	b.datas = make(map[uint32]*proto)
	return b
}