package alsa

// #include <alsa/asoundlib.h>
//
import "C"
import (
	"reflect"
	"runtime"
	"sync/atomic"
	"unsafe"
)

type Output interface {
	Ptr() unsafe.Pointer
}

type File struct {
	inner *C.snd_output_t
}

func OpenOutput(name, mode string) (*File, error) {
	output := &File{}
	cName := C.CString(name)
	cMode := C.CString(mode)
	defer func() {
		C.free(unsafe.Pointer(cName))
		C.free(unsafe.Pointer(cMode))
	}()
	rc := C.snd_output_stdio_open(&output.inner, cName, cMode)
	if rc < 0 {
		return nil, NewError(int(rc))
	}
	runtime.SetFinalizer(output, (*File).Close)
	return output, nil
}

func AttachStdout() (*File, error) {
	return attach(C.stdout)
}

func AttachStderr() (*File, error) {
	return attach(C.stderr)
}

func attach(file *C.FILE) (*File, error) {
	output := &File{}
	rc := C.snd_output_stdio_attach(&output.inner, file, C.int(0))
	if rc < 0 {
		return nil, NewError(int(rc))
	}
	runtime.SetFinalizer(output, (*File).Close)
	return output, nil
}

func (f *File) Close() error {
	for {
		c := unsafe.Pointer(f.inner)
		if c == nil {
			break
		}
		i := f.inner
		if atomic.CompareAndSwapPointer(&c, c, nil) {
			rc := C.snd_output_close(i)
			if rc < 0 {
				panic(NewError(int(rc)))
			}
			runtime.SetFinalizer(f, nil)
		}
	}

	return nil
}

func (f *File) Ptr() unsafe.Pointer {
	return unsafe.Pointer(f.inner)
}

type Buffer struct {
	inner *C.snd_output_t
}

func OpenBuffer() (*Buffer, error) {
	buf := &Buffer{}
	rc := C.snd_output_buffer_open(&buf.inner)
	if rc < 0 {
		return nil, NewError(int(rc))
	}
	runtime.SetFinalizer(buf, (*Buffer).Close)
	return buf, nil
}

func (b *Buffer) Buf() []byte {
	var cBuf *C.char
	cLen := C.snd_output_buffer_string(b.inner, &cBuf)
	if cLen <= 0 {
		return nil
	}

	var slice []byte
	header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	header.Len = int(cLen)
	header.Cap = int(cLen)
	header.Data = uintptr(unsafe.Pointer(cBuf))
	return slice
}

//func (b *Buffer) Read(p []byte) (n int, err error) {
//	m := 0
//	for {
//		m += n
//		tmp := p[m:]
//		n, err = b.innerBuf.Read(tmp)
//		if err != nil && err != io.EOF {
//			return
//		}
//		if n == len(tmp) {
//			return n, nil
//		}
//
//		var cBuf *C.char
//		cLen := C.snd_output_buffer_steal(b.inner, &cBuf)
//		if cLen <= 0 {
//			if n > 0 {
//				return n, nil
//			} else {
//				return 0, io.EOF
//			}
//		} else {
//			var slice []byte
//			header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
//			header.Len = int(cLen)
//			header.Cap = int(cLen)
//			header.Data = uintptr(unsafe.Pointer(cBuf))
//			_, err = b.innerBuf.Write(slice)
//			if err != nil {
//				return
//			}
//		}
//	}
//}

func (b *Buffer) Close() error {
	for {
		c := unsafe.Pointer(b.inner)
		if c == nil {
			break
		}
		i := b.inner
		if atomic.CompareAndSwapPointer(&c, c, nil) {
			rc := C.snd_output_close(i)
			if rc < 0 {
				panic(NewError(int(rc)))
			}
			runtime.SetFinalizer(b, nil)
			//b.innerBuf = nil
		}
	}

	return nil
}

func (b *Buffer) Ptr() unsafe.Pointer {
	return unsafe.Pointer(b.inner)
}
