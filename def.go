package reliableudp

import (
	"encoding/binary"
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
//	| byte  | byte  |  data     |
//      | types |  seq  |  bytes    |
//
// in case of types_init
//	| byte  | byte  |  data(4 bytes) |
//      | 1     |  0    |  totallen      |
// in case of types_trans
//	| byte  | byte  |  data     |
//      | 2     | seqid |  data...  |
// in case of types_end
//	| byte  | byte  |  data(4 bytes) |
//      | 3     | seqid |  totallen      |

const udp_datalen = 1450

type proto struct {
	types int
	seq   int
	data  []byte
}

type protobuffer struct {
	i     proto
	datas map[int][]proto
	e     proto
}

func (this *protobuffer) init(data []byte) {
	this.i.types = types_init
	this.i.seq = 0
	binary.BigEndian.PutUint32(this.i.data, uint32(len(data)))



	this.e.types = types_end
	this.e.seq = 0
	binary.BigEndian.PutUint32(this.e.data, uint32(len(data)))
}
