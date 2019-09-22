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

func (this *defaultLogger) Debug(v ...interface{})  { log.Debuglogger(v...) }
func (this *defaultLogger) Info(v ...interface{})   { log.Infologger(v...) }
func (this *defaultLogger) Notice(v ...interface{}) { log.Noticelogger(v...) }
func (this *defaultLogger) Warn(v ...interface{})   { log.Warnlogger(v...) }
func (this *defaultLogger) Error(v ...interface{})  { log.Errorlogger(v...) }
func (this *defaultLogger) Panic(v ...interface{})  { log.Paniclogger(v...) }
func (this *defaultLogger) Alert(v ...interface{})  { log.Alertlogger(v...) }
func (this *defaultLogger) Fatal(v ...interface{})  { log.Fatallogger(v...) }


var glog logger

func SetLogger(l logger){
	glog = l
}
