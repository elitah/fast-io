package fast_io

import (
	"errors"
	"io"
	"sync"
)

var (
	blackSize = 32
)

var (
	ENOBUF = errors.New("no buffer can be used")

	pool = sync.Pool{
		New: func() interface{} {
			return make([]byte, blackSize*1024)
		},
	}
)

func SetBlockSize(size int) int {
	if 0 < size {
		blackSize = size
	}
	return blackSize
}

func Copy(dst io.Writer, src io.Reader) (int64, error) {
	// 通过WriterTo接口快速拷贝
	if wt, ok := src.(io.WriterTo); ok {
		return wt.WriteTo(dst)
	}

	// 通过ReadFrom接口快速拷贝
	if rt, ok := dst.(io.ReaderFrom); ok {
		return rt.ReadFrom(src)
	}

	// 通过pool获取缓存
	if buffer := getBuffer(); nil != buffer {
		// 缓冲释放
		defer putBuffer(buffer)
		// 复制数据
		return io.CopyBuffer(dst, src, buffer)
	}

	return 0, ENOBUF
}

func FastCopy(_io ...io.ReadWriteCloser) (cnt_recv, cnt_send int64) {
	if 2 <= len(_io) {
		exit := make(chan struct{})
		// 协程
		go func() {
			// 拷贝
			cnt_recv, _ = Copy(_io[0], _io[1])
			// 关闭
			_io[0].Close()
			_io[1].Close()
			// 关闭通道
			close(exit)
		}()
		//
		cnt_send, _ = Copy(_io[1], _io[0])
		//
		_io[0].Close()
		_io[1].Close()
		//
		<-exit
	}

	return
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
