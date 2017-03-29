#include "c_bridge.h"
#include <string.h>
#include <unicode/utypes.h>
#include <unicode/ucsdet.h>
#include <stdlib.h>
#include <unicode/ucnv.h>

// See description in c_bridge.h
const int detectCharset(void        *detector,
                        void        *input,
                        int         input_len,
                        int         *status,
                        MatchData   *matchBuffer,
                        int         matchBufferSize) {

    // Put input bytes in the detector.
    ucsdet_setText((UCharsetDetector*)detector, (char*)input, input_len, status);
    if U_FAILURE(*status) {
        return 0;
    }

    // Prepare vars for returned count and guesses.
    int matchCount;
    const UCharsetMatch **bestGuesses;

    // Perform analysis and return all guesses and their count.
    bestGuesses = ucsdet_detectAll((UCharsetDetector*)detector, &matchCount, status);
    if U_FAILURE(*status) {
        return 0;
    }

    // Fill the matchBuffer. Its size is matchBufferSize, so it is filled with
    // less or equal to matchBufferSize number of entries.
    int i;
    int retCount = matchCount > matchBufferSize ? matchBufferSize : matchCount;

    for (i = 0; i < retCount; i++) {

        const UCharsetMatch* bestGuess = bestGuesses[i];
        const char *bestGuessedCharset = NULL;
        const char *bestGuessedLanguage = NULL;

        // Fill guessed encoding
        bestGuessedCharset = ucsdet_getName(bestGuess, status);
        if U_FAILURE(*status) {
            return 0;
        }

        // Fill guessed language
        bestGuessedLanguage = ucsdet_getLanguage(bestGuess, status);
        if U_FAILURE(*status) {
            return 0;
        }

        // Fill its confidence rating
        int32_t conf = ucsdet_getConfidence(bestGuess, status);
        if U_FAILURE(*status) {
            return 0;
        }

        matchBuffer[i].confidence = conf;
        matchBuffer[i].charset = bestGuessedCharset;
        matchBuffer[i].language = bestGuessedLanguage;
    }

    // Return the number of guesses put into matchBuffer.
    return retCount;
}

// See description in c_bridge.h
int convertToUtf16(const char   *srcEncoding,
                   UChar        *dest,
                   int32_t      destCapacity,
                   const char   *src,
                   int32_t      srcLength,
                   int          *status){
    UConverter *conv;

    conv = ucnv_open(srcEncoding, status);
    if U_FAILURE(*status) {
        return 0;
    }

    /* Convert from original encoding to UTF-16 */
    int len = ucnv_toUChars(conv, dest, destCapacity, src, srcLength, status);
    if U_FAILURE(*status) {
        return 0;
    }

    ucnv_close(conv);

    return len;
}

// See description in c_bridge.h
int convertFromUtf16(const char   *destEncoding,
                     char         *dest,
                     int32_t      destCapacity,
                     const UChar  *src,
                     int32_t      srcLength,
                     int          *status){
    UConverter *conv;

    conv = ucnv_open(destEncoding, status);
    if U_FAILURE(*status) {
        return 0;
    }

    /* Convert from UTF-16 to destination encoding */
    int len = ucnv_fromUChars(conv, dest, destCapacity, src, srcLength, status);
    if U_FAILURE(*status) {
        return 0;
    }

    ucnv_close(conv);

    return len;
}
