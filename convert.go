package icu

// #cgo pkg-config: icu-i18n
// #include "c_bridge.h"
// #include "stdlib.h"
import "C"
import (
    "fmt"
    "sync"
    "unsafe"
)

const (
    DefaultMaxTextSize = 1024 * 1024    // Default value for the max text length in conversion operations
    utf8MaxCharSize = 4
    utf16MaxCharSize = 4
)

var (
    Utf8CString = C.CString("UTF-8")
)

// CharsetConverter provides ICU charset conversion functionality.
type CharsetConverter struct {
    utf16Buffer   []byte
    utf8Buffer    []byte
    maxTextSize   int
    cMutex        sync.Mutex // Mutex used to guarantee thread safety for ICU calls
}

// NewCharsetConverter creates a new charset converter. It doesn't need to be closed as
// it doesn't allocate any resources.
//
// For better performance, conversion buffers are not allocated on each operation. Instead they
// are created in memory once and then used. 'maxTextSize' sets the size of these buffers.
// ICU library would return error if any processed text is longer than this parameter.
//
// NOTE:
//
// UTF8 uses 1 to 4 bytes for each symbol.
// UTF16 uses 2 bytes to 4 bytes for each symbol.
//
// So, to guarantee successful conversion of text with size = 'maxTextSize' we need:
//     maxTextSize * 8 bytes    (utf8 buffer + utf16 buffer).
func NewCharsetConverter(maxTextSize int) (*CharsetConverter) {
    conv := new(CharsetConverter)

    conv.utf16Buffer = make([]byte, utf16MaxCharSize * maxTextSize)
    conv.utf8Buffer = make([]byte, utf8MaxCharSize * maxTextSize)

    return conv
}

// ConvertToUtf8 converts input bytes encoded with srcEncoding to UTF-8.
func (conv *CharsetConverter) ConvertToUtf8(input []byte, srcEncoding string) ([]byte, error) {
    // As described in c_bridge.h, conversion operations are not thread safe and
    // should be called consequently. So a mutex is used here.
    conv.cMutex.Lock()
    defer conv.cMutex.Unlock()

    inputLen := len(input)
    if inputLen == 0 {
        return nil, fmt.Errorf("Nil length of input")
    }

    var status int

    encCString := C.CString(srcEncoding)
    inputCString := C.CString(string(input))

    defer C.free(unsafe.Pointer(encCString))
    defer C.free(unsafe.Pointer(inputCString))

    convLen := C.convertToUtf16(
            encCString,
            (*C.UChar)(unsafe.Pointer(&conv.utf16Buffer[0])),
            C.int32_t(len(conv.utf16Buffer)),
            inputCString,
            C.int32_t(len(input)),
            (*C.int)(unsafe.Pointer(&status)))

    if isSuccess(status) {
        nConvLen := C.convertFromUtf16(
            Utf8CString,
            (*C.char)(unsafe.Pointer(&conv.utf8Buffer[0])),
            C.int32_t(len(conv.utf8Buffer)),
            (*C.UChar)(unsafe.Pointer(&conv.utf16Buffer[0])),
            C.int32_t(convLen),
            (*C.int)(unsafe.Pointer(&status)))

        if isSuccess(status) {
            resStr := conv.utf8Buffer[:nConvLen]
            return ([]byte)(resStr), nil
        }
    }

    return nil, fmt.Errorf("ICU Error code returned: %d", status)
}
