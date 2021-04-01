// +build linux,!nopcre

package fastregex

/*
#cgo LDFLAGS: -lpcre
#cgo CFLAGS: -I/opt/local/include
#include <pcre.h>
#include <string.h>
*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

const Version = "libpcre"

var JITCompile = true

type PCRERegexp struct {
	p          *C.pcre
	extra      *C.pcre_extra
	groupCount int
}

func Compile(expr string) (Regexp, error) {
	cExpr := C.CString(expr)
	defer C.free(unsafe.Pointer(cExpr))

	var err *C.char
	var errOffset C.int
	compiled := C.pcre_compile(cExpr, 0, &err, &errOffset, nil)
	if compiled == nil {
		return nil, &CompileError{
			expr,
			C.GoString(err),
			int(errOffset),
		}
	}

	var study *C.pcre_extra
	if JITCompile {
		study = C.pcre_study(compiled, C.PCRE_STUDY_JIT_COMPILE, &err)
		if err != nil {
			return nil, &CompileError{
				expr,
				C.GoString(err),
				0,
			}
		}
	}

	var groupCount C.int
	C.pcre_fullinfo(compiled, nil, C.PCRE_INFO_CAPTURECOUNT, unsafe.Pointer(&groupCount))

	ret := &PCRERegexp{
		p:          compiled,
		extra:      study,
		groupCount: int(groupCount) + 1,
	}
	runtime.SetFinalizer(ret, func(f *PCRERegexp) {
		if f.extra != nil {
			C.pcre_free_study(f.extra)
		}
		C.free(unsafe.Pointer(f.p))
	})
	return ret, nil
}

func MustCompile(expr string) Regexp {
	c, err := Compile(expr)
	if err != nil {
		panic(err)
	}
	return c
}

func (s *PCRERegexp) GroupCount() int {
	return s.groupCount
}

func (s *PCRERegexp) Match(b []byte) bool {
	bPtr := (*C.char)(unsafe.Pointer(&b[0]))
	ret := C.pcre_exec(s.p, s.extra, bPtr, C.int(len(b)), 0, 0, nil, 0)
	return ret >= 0
}

func (s *PCRERegexp) MatchString(str string) bool {
	bPtr := *(**C.char)(unsafe.Pointer(&str)) // This style of syntax prevents a string-copy
	ret := C.pcre_exec(s.p, s.extra, bPtr, C.int(len([]byte(str))), 0, 0, nil, 0)
	return ret >= 0
}

func (s *PCRERegexp) FindSubmatchIndex(b []byte) []int {
	groups := make([]C.int, 3*s.groupCount)
	bPtr := (*C.char)(unsafe.Pointer(&b[0]))
	ret := C.pcre_exec(s.p, s.extra, bPtr, C.int(len(b)), 0, 0, &groups[0], C.int(len(groups)))
	if ret < 0 {
		return nil
	}

	converted := make([]int, s.groupCount*2)
	for i := 0; i < len(converted); i++ {
		converted[i] = int(groups[i])
	}

	return converted
}

type CompileError struct {
	Expr    string
	Message string
	Offset  int
}

var _ error = &CompileError{}

func (s *CompileError) Error() string {
	return fmt.Sprintf("error in expression '%s', offset %d: %s", s.Expr, s.Offset, s.Message)
}
