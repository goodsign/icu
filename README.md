About
==========

Cgo binding for icu4c C library detection and conversion functions. Guaranteed compatibility with version 50.1.

Installation
==========

Installation consists of several simple steps. They may be a bit different on your target system (e.g. require more permissions) so adapt them to the parameters of your system.

### Install build-essential

Make sure you have **build-essential** installed. Otherwise icu would fail on the configuration stage.

Installation example using apt-get (Ubuntu):

```
sudo apt-get install build-essential
```

### Install pkg-config

Make sure you have **pkg-config** installed.

Installation example using apt-get (Ubuntu):

```
sudo apt-get install pkg-config
```

### Get icu4c C library code

Download and unarchive original icu4c archive from [icu download section](http://site.icu-project.org/download).

Example (for version 50.1):

```
wget http://download.icu-project.org/files/icu4c/50.1/icu4c-50_1-src.tgz
tar -zxvf icu4c-50_1-src.tgz
mv -i ./icu ~/where-you-store-libs
```

NOTE: If this link is not working or there are some problems with downloading, there is a stable version 50.1 snapshot saved in [Github Downloads](https://github.com/downloads/goodsign/icu/icu4c-50_1-src.tgz).

### Build and install icu4c C library

From the directory, where you unarchived icu4c, run:

```
cd source
./configure
make
sudo make install
sudo ldconfig
```

### Install Go wrapper

```
go get github.com/goodsign/icu
go test github.com/goodsign/icu (must PASS)
```

Installation notes
==========

* Make sure that you have your local library paths set correctly and that installation was successful. Otherwise, **go build** or **go test** may fail.

* icu4c is installed in your local library directory (e.g. **/usr/local/lib**) and puts its libraries there. This path should be registered in your system (using ldconfig or exporting LD_LIBRARY_PATH, etc.) or the linker would fail.

* icu4c installs its header files to local include folders (e.g. **/usr/local/include/unicode**) so there is no need to have additional .h files with this package, but the system must be properly set up to detect .h files in those directories.

Usage
==========

Note: check icu documentation for returned encoding identifiers.

Detector
----------

```go
// Create detector
detector, err := NewCharsetDetector()
    
if err != nil {
    //... Handle error ...
}
defer detector.Close()

// Guess encoding
encMatches, err := detector.GuessCharset(encodedText)

if err != nil {
    //... Handle error ...
}

// Get charset with max confidence (goes first)
maxenc := encMatches[0].Charset

// Use maxenc. 
// ...
```

Converter
----------

```go
...

// Create converter
converter := NewCharsetConverter(DefaultMaxTextSize)

// Convert to utf-8
converted, err := converter.ConvertToUtf8(encodedText, maxenc)

if nil != err {
    //... Handle error ...
}
```

Usage notes
==========

* Check **NewCharsetConverter** func comments for details on max text size parameter.
* Often you would use detector and converter in pair. So, the 'converter' usage example actually continues the 'detector' example and uses the 'maxenc' result from it.

More info
----------

For more information on icu refer to the original [website](http://site.icu-project.org/), which contains links on theory and other details.

icu4c Licence
==========

ICU is released under a nonrestrictive open source license that is suitable for use with both commercial software and with other open source or free software.

[LICENCE file](https://github.com/goodsign/icu/blob/master/LICENCE_icu)

Licence
==========

The goodsign/icu binding is released under the [BSD Licence](http://opensource.org/licenses/bsd-license.php)

[LICENCE file](https://github.com/goodsign/icu/blob/master/LICENCE)