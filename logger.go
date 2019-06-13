package reliableudp

import "github.com/hqbobo/log"

type logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Notice(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
	Panic(v ...interface{})
	init()
}

type defaultLogger struct {
}

func (*defaultLogger) init() {

}

func (this *defaultLogger) Debug(v ...interface{})  { log.Debug(v...) }
func (this *defaultLogger) Info(v ...interface{})   { log.Info(v...) }
func (this *defaultLogger) Notice(v ...interface{}) { log.Notice(v...) }
func (this *defaultLogger) Warn(v ...interface{})   { log.Warn(v...) }
func (this *defaultLogger) Error(v ...interface{})  { log.Error(v...) }
func (this *defaultLogger) Panic(v ...interface{})  { log.Panic(v...) }
func (this *defaultLogger) Alert(v ...interface{})  { log.Alert(v...) }
func (this *defaultLogger) Fatal(v ...interface{})  { log.Fatal(v...) }


var glog logger

func SetLogger(l logger){
	glog = l
}
