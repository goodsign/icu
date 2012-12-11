#ifndef __C_BRIDGE_H__
#define __C_BRIDGE_H__

// C_BRIDGE is a bridge between go and native pure c functions used to 
// operate with ICU library code.

#include <unicode/utypes.h>
#include <unicode/ucsdet.h>

// MatchData contains information about one 'guess' of the
// encoding detector. It contains the guessed charset (ICU string identifiers,
// see ICU documentation for them) and a confidence coefficient, which is a
// number between 0 and 100 (100 is the best).
typedef struct MatchData {
  const char* charset;
  const char* language;
  short int confidence;
} MatchData;

// detectCharset performs the detection (guessing) operation using a given detector (ICU internals),
// input data (bytes), input length and error status pointer (Read ICU docs abour error codes).
//
// After the detection is performed, all possible matches are put into the matchBuffer. If there are
// more results than matchBufferSize, then only matchBufferSize entries are put (So no overflow can
// ever happen).
//
// The results of this function are put into the matchBuffer, so it MUST NOT be called asynchronously.
// Caller should guarantee thread safety and perform locks while working with it.
const int detectCharset(void       *detector, 
                        void       *input, 
                        int        input_len, 
                        int        *status, 
                        MatchData  *matchBuffer, 
                        int        matchBufferSize);

// convertToUtf16 performs conversion from any encoding to utf16. Utf16 is the ICU standard so
// it is easier to convert to/from it.
// 
// The results of this function are put into the dest buffer, so it MUST NOT be called asynchronously.
// Caller should guarantee thread safety and perform locks while working with it.
int convertToUtf16(const char   *srcEncoding,
                   UChar        *dest, 
                   int32_t      destCapacity,
                   const char   *src,
                   int32_t      srcLength,
                   int          *status);

// convertFromUtf16 performs conversion from utf16 to any other encoding. Utf16 is the ICU standard so
// it is easier to convert to/from it.
// 
// The results of this function are put into the dest buffer, so it MUST NOT be called asynchronously.
// Caller should guarantee thread safety and perform locks while working with it.
int convertFromUtf16(const char   *destEncoding,
                     char         *dest, 
                     int32_t      destCapacity,
                     const UChar  *src,
                     int32_t      srcLength,
                     int          *status);


#endif //__C_BRIDGE_H__