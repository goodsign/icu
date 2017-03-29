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
    U_ZERO_ERROR        = 0          // ICU common constant error code which means that no error occured
    U_ERROR_LIMIT       = 0x7FFFFFFF // Dirty hack, negative error codes are are being turned into large positive ints
    MatchDataBufferSize = 25         // Size of the buffer for detection results (Max count of returned guesses per detect call)
)

// Go implementation of the icu U_SUCCESS macro. Negative status codes are
// warnings, 0 is a success without warnings, > 0 is an error
func isSuccess(status int) bool {
    return status <= U_ZERO_ERROR || status >= U_ERROR_LIMIT
}

// Go implementation of the icu U_FAILURE macro.
func isFailure(status int) bool {
    return status > U_ZERO_ERROR && status < U_ERROR_LIMIT
}

// CharsetDetector provides ICU charset detection functionality.
type CharsetDetector struct {
    ptr         *C.UCharsetDetector // ICU struct needed for detection
    resBuffer   [MatchDataBufferSize]C.MatchData
    gMutex      sync.Mutex // Mutex used to guarantee thread safety for ICU calls
}

// An equivalent of MatchData C structure (see c_bridge.h)
type Match struct {
    Charset string
    Language string
    Confidence int
}

// Creates new charset detector. If it is successfully created, it
// must be closed as it needs to free native ICU resources.
func NewCharsetDetector() (*CharsetDetector, error) {
    det := new(CharsetDetector)

    var status int
    statusPtr := unsafe.Pointer(&status)

    det.ptr = C.ucsdet_open((*C.UErrorCode)(statusPtr))

    if isFailure(status) {
        return nil, fmt.Errorf("ICU Error code returned: %d", status)
    }

    return det, nil
}

func (det *CharsetDetector) GuessCharset(input []byte) (matches []Match, err error) {

    // As described in c_bridge.h, detection operations are not thread safe and
    // should be called consequently. So a mutex is used here.
    det.gMutex.Lock()
    defer det.gMutex.Unlock()

    inputLen := len(input)
    if inputLen == 0 {
        return nil, fmt.Errorf("Input data len is 0")
    }

    var status int

    // Perform detection. Guess count is the number of matches returned.
    // The matches themself are put in the result buffer
    guessCount := C.detectCharset(
        unsafe.Pointer(det.ptr),
        unsafe.Pointer(&input[0]),
        C.int(inputLen),
        (*C.int)(unsafe.Pointer(&status)),
        (*C.MatchData)(unsafe.Pointer(&det.resBuffer[0])),
        C.int(MatchDataBufferSize))

    if isSuccess(status) {
        // Convert the returned number of entries from result buffer to a slice
        // that will be returned
        count := int(guessCount)
        mt := make([]Match, count, count)

        for i := 0; i < count; i++ {
            mData := det.resBuffer[i]
            charset := C.GoString(mData.charset)
            language := C.GoString(mData.language)
            mt[i] = Match{charset, language, int(mData.confidence)}
        }

        return mt, nil
    }

    return nil, fmt.Errorf("ICU Error code returned: %d", status)
}

// Close frees native C resources
func (det *CharsetDetector) Close() {
    det.gMutex.Lock()
    defer det.gMutex.Unlock()

    if det.ptr != nil {
        C.ucsdet_close(det.ptr)
    }
}
