package fast_io

import (
	"io"
	"sync"
)

var (
	pool = sync.Pool{
		New: func() interface{} {
			// 64k
			return make([]byte, 64*1024)
		},
	}
)

func FastCopy(_io ...io.ReadWriteCloser) {
	if 2 <= len(_io) {
		exit := make(chan struct{})
		//
		go func() {
			//
			if buffer := getBuffer(); nil != buffer {
				io.CopyBuffer(_io[0], _io[1], buffer)
				//
				putBuffer(buffer)
			}
			//
			_io[0].Close()
			_io[1].Close()
			//
			close(exit)
		}()
		//
		if buffer := getBuffer(); nil != buffer {
			io.CopyBuffer(_io[1], _io[0], buffer)
			//
			putBuffer(buffer)
		}
		//
		_io[0].Close()
		_io[1].Close()
		//
		<-exit
	}
}

func getBuffer() []byte {
	if buffer, ok := pool.Get().([]byte); ok {
		return buffer
	}
	return nil
}

func putBuffer(buffer []byte) {
	pool.Put(buffer)
}
