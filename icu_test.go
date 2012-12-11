package icu

import (
    "io/ioutil"
    "testing"
    "regexp"
)

var (
    IcuTestLineRx *regexp.Regexp = regexp.MustCompile(`\A\[(?P<encoded>.+)\]\s*\[(?P<utfexpected>.+)\].*\z`)
)

const (
    TestConfigPath = "defaultcfg/conf.txt"
)

func testConversion(t *testing.T, encFileName string, expFileName string) {
    // Create detector
    detector, err := NewCharsetDetector()
    
    if nil != err {
        t.Fatalf("Cannot create detector: %s", err)
    }
    defer detector.Close()

    // Create converter
    converter := NewCharsetConverter(DefaultMaxTextSize)

    // Open files

    enc, err := ioutil.ReadFile(encFileName)

    if nil != err {
        t.Error(err)
        return
    }

    exp, err := ioutil.ReadFile(expFileName)

    if nil != err {
        t.Error(err)
        return
    }

    // Guess encoding
    encMatches, err := detector.GuessCharset(enc)

    if nil != err {
        t.Error(err)
        return
    }

    // Get charset with max confidence
    maxenc := encMatches[0].Charset

    // Convert to utf-8
    converted, err := converter.ConvertToUtf8(enc, maxenc)

    if nil != err {
        t.Error(err)
        return
    }

    t.Logf("Encoded file: '%s' Expected file: [%s] Detected charset: [%s]",
           encFileName, 
           expFileName,
           maxenc)

    // Compare converted result and expected result from file.
    if string(converted) != string(exp) {
        t.Errorf("Encoded file: '%s' Expected file: [%s] Detected charset: [%s] Expected utf8: [%s] Got utf8: [%s]", 
                 encFileName, 
                 expFileName,
                 maxenc,
                 exp, 
                 string(converted))
    }
}

func TestDefault(t *testing.T) {
    testConversion(t,       "test/koi8r.txt",       "test/koi8r_to_utf.txt")
    testConversion(t,       "test/windows88591.txt","test/windows88591_to_utf.txt")
    testConversion(t,       "test/utf8.txt",        "test/utf8_to_utf.txt")
}