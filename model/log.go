package model

type Logger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Infofln(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Debugfln(format string, args ...interface{})
	Warning(args ...interface{})
	Warningf(format string, args ...interface{})
	Warningfln(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorfln(format string, args ...interface{})
	SetLevel(int)
	GetLevel() int
}
