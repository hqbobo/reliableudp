package reliableudp

import (
	"github.com/hqbobo/log"
	"net"
)

type udpreader interface {
	OnRecive(data []byte, saddr net.Addr)
}

func init() {
	log.InitLog(log.AllLevels...)
	log.SetCode(true)
	log.SetPathFilter("github.com/hqbobo/")
	glog = &defaultLogger{}
}