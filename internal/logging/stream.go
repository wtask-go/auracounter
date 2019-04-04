package logging

import (
	"io"
	"time"
)

// stream - wrapper for standart go logger to implement loggin.Interface
type stream struct {
	// mu  sync.Mutex  // mutex not needed due to write operation inside protected log.Output operation
	out io.Writer // access to writer, used inside logger
	// logger         log.Logger // without prefix and time gen
	flags          int // logger flags, only Llongfile and Lshortfile
	prefix         []byte
	timeFunc       func() time.Time       // also may define timezone
	timeFormatFunc func(time.Time) string // also may convert time zone
}

func (s stream) Write(bytes []byte) (int, error) {
	return s.out.Write(append(s.getPrefix(), bytes...))
}

func (s stream) getPrefix() []byte {
	// NOTE when timeFunc() is used here, it is catching more later time than Facade method was called
	return append(s.prefix, []byte(s.timeFormatFunc(s.timeFunc()))...)
}
