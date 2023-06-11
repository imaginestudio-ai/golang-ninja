## The go doc and godoc utilities

The Go distribution comes with a plethora of tools that can make your life as a programmer easier. Two of these tools are the `go doc` subcommand and `godoc` utility, which allow you to see the documentation of existing Go functions and packages without needing an internet connection. However, if you prefer viewing the Go documentation online, you can visit [https://pkg.go.dev/](https://pkg.go.dev/). As `godoc` is not installed by default, you might need to install it by running `go install golang.org/x/tools/cmd/godoc@latest`.

The `go doc` command can be executed as a normal command-line application that displays its output on a terminal, and `godoc` as a command-line application that starts a web server. In the latter case, you need a web browser to look at the Go documentation. The first utility is similar to the UNIX `man(1)` command, but for Go functions and packages.

The number after the name of a UNIX program or system call refers to the section of the manual a manual page belongs to. Although most of the names can be found only once in the manual pages, which means that putting the section number is not required, there are names that can be located in multiple sections because they have multiple meanings, such as `crontab(1)` and `crontab(5)`. Therefore, if you try to retrieve the manual page of a name with multiple meanings without stating its section number, you will get the entry that has the smallest section number.

So, in order to find information about the `Printf()` function of the `fmt` package, you should execute the following command:

```markup
$ go doc fmt.Printf
```

Similarly, you can find information about the entire `fmt` package by running the following command:

```markup
$ go doc fmt
```

The second utility requires executing `godoc` with the `-http` parameter:

```markup
$ godoc -http=:8001
```

The numeric value in the preceding command, which in this case is `8001`, is the port number the HTTP server will listen to. As we have omitted the IP address, `godoc` is going to listen to all network interfaces.

You can choose any port number that is available provided that you have the right privileges. However, note that port numbers `0`\-`1023` are restricted and can only be used by the root user, so it is better to avoid choosing one of those and pick something else, provided that it is not already in use by a different process.

You can omit the equals sign in the presented command and put a space character in its place. So, the following command is completely equivalent to the previous one:

```markup
$ godoc -http :8001
```

After that, you should point your web browser to the `http://localhost:8001/` URL in order to get the list of available Go packages and browse their documentation. If you are using Go for the first time, you will find the Go documentation very handy for learning the parameters and the return values of the functions you want to use—as you progress in your Go journey, you will use the Go documentation for learning the gory details of the functions and variables that you want to use.

Note:  This example below is of the locally hosted doc server page.

#  Go Documentation Server Example



| Name | Synopsis |
| --- | --- |
| [archive](http://localhost:8001/pkg/archive/) |  |
| [tar](http://localhost:8001/pkg/archive/tar/) | Package tar implements access to tar archives. |
| [zip](http://localhost:8001/pkg/archive/zip/) | Package zip provides support for reading and writing ZIP archives. |
| [bufio](http://localhost:8001/pkg/bufio/) | Package bufio implements buffered I/O. It wraps an io.Reader or io.Writer object, creating another object (Reader or Writer) that also implements the interface but provides buffering and some help for textual I/O. |
| [builtin](http://localhost:8001/pkg/builtin/) | Package builtin provides documentation for Go's predeclared identifiers. |
| [bytes](http://localhost:8001/pkg/bytes/) | Package bytes implements functions for the manipulation of byte slices. |
| [compress](http://localhost:8001/pkg/compress/) |  |
| [bzip2](http://localhost:8001/pkg/compress/bzip2/) | Package bzip2 implements bzip2 decompression. |
| [flate](http://localhost:8001/pkg/compress/flate/) | Package flate implements the DEFLATE compressed data format, described in RFC 1951. |
| [gzip](http://localhost:8001/pkg/compress/gzip/) | Package gzip implements reading and writing of gzip format compressed files, as specified in RFC 1952. |
| [lzw](http://localhost:8001/pkg/compress/lzw/) | Package lzw implements the Lempel-Ziv-Welch compressed data format, described in T. A. Welch, “A Technique for High-Performance Data Compression”, Computer, 17(6) (June 1984), pp 8-19. |
| [zlib](http://localhost:8001/pkg/compress/zlib/) | Package zlib implements reading and writing of zlib format compressed data, as specified in RFC 1950. |
| [container](http://localhost:8001/pkg/container/) |  |
| [heap](http://localhost:8001/pkg/container/heap/) | Package heap provides heap operations for any type that implements heap.Interface. |
| [list](http://localhost:8001/pkg/container/list/) | Package list implements a doubly linked list. |
| [ring](http://localhost:8001/pkg/container/ring/) | Package ring implements operations on circular lists. |
| [context](http://localhost:8001/pkg/context/) | Package context defines the Context type, which carries deadlines, cancellation signals, and other request-scoped values across API boundaries and between processes. |
| [crypto](http://localhost:8001/pkg/crypto/) | Package crypto collects common cryptographic constants. |
| [aes](http://localhost:8001/pkg/crypto/aes/) | Package aes implements AES encryption (formerly Rijndael), as defined in U.S. Federal Information Processing Standards Publication 197. |
| [cipher](http://localhost:8001/pkg/crypto/cipher/) | Package cipher implements standard block cipher modes that can be wrapped around low-level block cipher implementations. |
| [des](http://localhost:8001/pkg/crypto/des/) | Package des implements the Data Encryption Standard (DES) and the Triple Data Encryption Algorithm (TDEA) as defined in U.S. Federal Information Processing Standards Publication 46-3. |
| [dsa](http://localhost:8001/pkg/crypto/dsa/) | Package dsa implements the Digital Signature Algorithm, as defined in FIPS 186-3. |
| [ecdsa](http://localhost:8001/pkg/crypto/ecdsa/) | Package ecdsa implements the Elliptic Curve Digital Signature Algorithm, as defined in FIPS 186-4 and SEC 1, Version 2.0. |
| [ed25519](http://localhost:8001/pkg/crypto/ed25519/) | Package ed25519 implements the Ed25519 signature algorithm. |
| [elliptic](http://localhost:8001/pkg/crypto/elliptic/) | Package elliptic implements the standard NIST P-224, P-256, P-384, and P-521 elliptic curves over prime fields. |
| [hmac](http://localhost:8001/pkg/crypto/hmac/) | Package hmac implements the Keyed-Hash Message Authentication Code (HMAC) as defined in U.S. Federal Information Processing Standards Publication 198. |
| [md5](http://localhost:8001/pkg/crypto/md5/) | Package md5 implements the MD5 hash algorithm as defined in RFC 1321. |
| [rand](http://localhost:8001/pkg/crypto/rand/) | Package rand implements a cryptographically secure random number generator. |
| [rc4](http://localhost:8001/pkg/crypto/rc4/) | Package rc4 implements RC4 encryption, as defined in Bruce Schneier's Applied Cryptography. |
| [rsa](http://localhost:8001/pkg/crypto/rsa/) | Package rsa implements RSA encryption as specified in PKCS #1 and RFC 8017. |
| [sha1](http://localhost:8001/pkg/crypto/sha1/) | Package sha1 implements the SHA-1 hash algorithm as defined in RFC 3174. |
| [sha256](http://localhost:8001/pkg/crypto/sha256/) | Package sha256 implements the SHA224 and SHA256 hash algorithms as defined in FIPS 180-4. |
| [sha512](http://localhost:8001/pkg/crypto/sha512/) | Package sha512 implements the SHA-384, SHA-512, SHA-512/224, and SHA-512/256 hash algorithms as defined in FIPS 180-4. |
| [subtle](http://localhost:8001/pkg/crypto/subtle/) | Package subtle implements functions that are often useful in cryptographic code but require careful thought to use correctly. |
| [tls](http://localhost:8001/pkg/crypto/tls/) | Package tls partially implements TLS 1.2, as specified in RFC 5246, and TLS 1.3, as specified in RFC 8446. |
| [x509](http://localhost:8001/pkg/crypto/x509/) | Package x509 parses X.509-encoded keys and certificates. |
| [pkix](http://localhost:8001/pkg/crypto/x509/pkix/) | Package pkix contains shared, low level structures used for ASN.1 parsing and serialization of X.509 certificates, CRL and OCSP. |
| [database](http://localhost:8001/pkg/database/) |  |
| [sql](http://localhost:8001/pkg/database/sql/) | Package sql provides a generic interface around SQL (or SQL-like) databases. |
| [driver](http://localhost:8001/pkg/database/sql/driver/) | Package driver defines interfaces to be implemented by database drivers as used by package sql. |
| [debug](http://localhost:8001/pkg/debug/) |  |
| [buildinfo](http://localhost:8001/pkg/debug/buildinfo/) | Package buildinfo provides access to information embedded in a Go binary about how it was built. |
| [dwarf](http://localhost:8001/pkg/debug/dwarf/) | Package dwarf provides access to DWARF debugging information loaded from executable files, as defined in the DWARF 2.0 Standard at http://dwarfstd.org/doc/dwarf-2.0.0.pdf |
| [elf](http://localhost:8001/pkg/debug/elf/) | Package elf implements access to ELF object files. |
| [gosym](http://localhost:8001/pkg/debug/gosym/) | Package gosym implements access to the Go symbol and line number tables embedded in Go binaries generated by the gc compilers. |
| [macho](http://localhost:8001/pkg/debug/macho/) | Package macho implements access to Mach-O object files. |
| [pe](http://localhost:8001/pkg/debug/pe/) | Package pe implements access to PE (Microsoft Windows Portable Executable) files. |
| [plan9obj](http://localhost:8001/pkg/debug/plan9obj/) | Package plan9obj implements access to Plan 9 a.out object files. |
| [embed](http://localhost:8001/pkg/embed/) | Package embed provides access to files embedded in the running Go program. |
| [encoding](http://localhost:8001/pkg/encoding/) | Package encoding defines interfaces shared by other packages that convert data to and from byte-level and textual representations. |
| [ascii85](http://localhost:8001/pkg/encoding/ascii85/) | Package ascii85 implements the ascii85 data encoding as used in the btoa tool and Adobe's PostScript and PDF document formats. |
| [asn1](http://localhost:8001/pkg/encoding/asn1/) | Package asn1 implements parsing of DER-encoded ASN.1 data structures, as defined in ITU-T Rec X.690. |
| [base32](http://localhost:8001/pkg/encoding/base32/) | Package base32 implements base32 encoding as specified by RFC 4648. |
| [base64](http://localhost:8001/pkg/encoding/base64/) | Package base64 implements base64 encoding as specified by RFC 4648. |
| [binary](http://localhost:8001/pkg/encoding/binary/) | Package binary implements simple translation between numbers and byte sequences and encoding and decoding of varints. |
| [csv](http://localhost:8001/pkg/encoding/csv/) | Package csv reads and writes comma-separated values (CSV) files. |
| [gob](http://localhost:8001/pkg/encoding/gob/) | Package gob manages streams of gobs - binary values exchanged between an Encoder (transmitter) and a Decoder (receiver). |
| [hex](http://localhost:8001/pkg/encoding/hex/) | Package hex implements hexadecimal encoding and decoding. |
| [json](http://localhost:8001/pkg/encoding/json/) | Package json implements encoding and decoding of JSON as defined in RFC 7159. |
| [pem](http://localhost:8001/pkg/encoding/pem/) | Package pem implements the PEM data encoding, which originated in Privacy Enhanced Mail. |
| [xml](http://localhost:8001/pkg/encoding/xml/) | Package xml implements a simple XML 1.0 parser that understands XML name spaces. |
| [errors](http://localhost:8001/pkg/errors/) | Package errors implements functions to manipulate errors. |
| [expvar](http://localhost:8001/pkg/expvar/) | Package expvar provides a standardized interface to public variables, such as operation counters in servers. |
| [flag](http://localhost:8001/pkg/flag/) | Package flag implements command-line flag parsing. |
| [fmt](http://localhost:8001/pkg/fmt/) | Package fmt implements formatted I/O with functions analogous to C's printf and scanf. |
| [go](http://localhost:8001/pkg/go/) |  |
| [ast](http://localhost:8001/pkg/go/ast/) | Package ast declares the types used to represent syntax trees for Go packages. |
| [build](http://localhost:8001/pkg/go/build/) | Package build gathers information about Go packages. |
| [constraint](http://localhost:8001/pkg/go/build/constraint/) | Package constraint implements parsing and evaluation of build constraint lines. |
| [constant](http://localhost:8001/pkg/go/constant/) | Package constant implements Values representing untyped Go constants and their corresponding operations. |
| [doc](http://localhost:8001/pkg/go/doc/) | Package doc extracts source code documentation from a Go AST. |
| [format](http://localhost:8001/pkg/go/format/) | Package format implements standard formatting of Go source. |
| [importer](http://localhost:8001/pkg/go/importer/) | Package importer provides access to export data importers. |
| [parser](http://localhost:8001/pkg/go/parser/) | Package parser implements a parser for Go source files. |
| [printer](http://localhost:8001/pkg/go/printer/) | Package printer implements printing of AST nodes. |
| [scanner](http://localhost:8001/pkg/go/scanner/) | Package scanner implements a scanner for Go source text. |
| [token](http://localhost:8001/pkg/go/token/) | Package token defines constants representing the lexical tokens of the Go programming language and basic operations on tokens (printing, predicates). |
| [types](http://localhost:8001/pkg/go/types/) | Package types declares the data types and implements the algorithms for type-checking of Go packages. |
| [hash](http://localhost:8001/pkg/hash/) | Package hash provides interfaces for hash functions. |
| [adler32](http://localhost:8001/pkg/hash/adler32/) | Package adler32 implements the Adler-32 checksum. |
| [crc32](http://localhost:8001/pkg/hash/crc32/) | Package crc32 implements the 32-bit cyclic redundancy check, or CRC-32, checksum. |
| [crc64](http://localhost:8001/pkg/hash/crc64/) | Package crc64 implements the 64-bit cyclic redundancy check, or CRC-64, checksum. |
| [fnv](http://localhost:8001/pkg/hash/fnv/) | Package fnv implements FNV-1 and FNV-1a, non-cryptographic hash functions created by Glenn Fowler, Landon Curt Noll, and Phong Vo. |
| [maphash](http://localhost:8001/pkg/hash/maphash/) | Package maphash provides hash functions on byte sequences. |
| [html](http://localhost:8001/pkg/html/) | Package html provides functions for escaping and unescaping HTML text. |
| [template](http://localhost:8001/pkg/html/template/) | Package template (html/template) implements data-driven templates for generating HTML output safe against code injection. |
| [image](http://localhost:8001/pkg/image/) | Package image implements a basic 2-D image library. |
| [color](http://localhost:8001/pkg/image/color/) | Package color implements a basic color library. |
| [palette](http://localhost:8001/pkg/image/color/palette/) | Package palette provides standard color palettes. |
| [draw](http://localhost:8001/pkg/image/draw/) | Package draw provides image composition functions. |
| [gif](http://localhost:8001/pkg/image/gif/) | Package gif implements a GIF image decoder and encoder. |
| [jpeg](http://localhost:8001/pkg/image/jpeg/) | Package jpeg implements a JPEG image decoder and encoder. |
| [png](http://localhost:8001/pkg/image/png/) | Package png implements a PNG image decoder and encoder. |
| [index](http://localhost:8001/pkg/index/) |  |
| [suffixarray](http://localhost:8001/pkg/index/suffixarray/) | Package suffixarray implements substring search in logarithmic time using an in-memory suffix array. |
| [io](http://localhost:8001/pkg/io/) | Package io provides basic interfaces to I/O primitives. |
| [fs](http://localhost:8001/pkg/io/fs/) | Package fs defines basic interfaces to a file system. |
| [ioutil](http://localhost:8001/pkg/io/ioutil/) | Package ioutil implements some I/O utility functions. |
| [log](http://localhost:8001/pkg/log/) | Package log implements a simple logging package. |
| [syslog](http://localhost:8001/pkg/log/syslog/) | Package syslog provides a simple interface to the system log service. |
| [math](http://localhost:8001/pkg/math/) | Package math provides basic constants and mathematical functions. |
| [big](http://localhost:8001/pkg/math/big/) | Package big implements arbitrary-precision arithmetic (big numbers). |
| [bits](http://localhost:8001/pkg/math/bits/) | Package bits implements bit counting and manipulation functions for the predeclared unsigned integer types. |
| [cmplx](http://localhost:8001/pkg/math/cmplx/) | Package cmplx provides basic constants and mathematical functions for complex numbers. |
| [rand](http://localhost:8001/pkg/math/rand/) | Package rand implements pseudo-random number generators unsuitable for security-sensitive work. |
| [mime](http://localhost:8001/pkg/mime/) | Package mime implements parts of the MIME spec. |
| [multipart](http://localhost:8001/pkg/mime/multipart/) | Package multipart implements MIME multipart parsing, as defined in RFC 2046. |
| [quotedprintable](http://localhost:8001/pkg/mime/quotedprintable/) | Package quotedprintable implements quoted-printable encoding as specified by RFC 2045. |
| [net](http://localhost:8001/pkg/net/) | Package net provides a portable interface for network I/O, including TCP/IP, UDP, domain name resolution, and Unix domain sockets. |
| [http](http://localhost:8001/pkg/net/http/) | Package http provides HTTP client and server implementations. |
| [cgi](http://localhost:8001/pkg/net/http/cgi/) | Package cgi implements CGI (Common Gateway Interface) as specified in RFC 3875. |
| [cookiejar](http://localhost:8001/pkg/net/http/cookiejar/) | Package cookiejar implements an in-memory RFC 6265-compliant http.CookieJar. |
| [fcgi](http://localhost:8001/pkg/net/http/fcgi/) | Package fcgi implements the FastCGI protocol. |
| [httptest](http://localhost:8001/pkg/net/http/httptest/) | Package httptest provides utilities for HTTP testing. |
| [httptrace](http://localhost:8001/pkg/net/http/httptrace/) | Package httptrace provides mechanisms to trace the events within HTTP client requests. |
| [httputil](http://localhost:8001/pkg/net/http/httputil/) | Package httputil provides HTTP utility functions, complementing the more common ones in the net/http package. |
| [pprof](http://localhost:8001/pkg/net/http/pprof/) | Package pprof serves via its HTTP server runtime profiling data in the format expected by the pprof visualization tool. |
| [mail](http://localhost:8001/pkg/net/mail/) | Package mail implements parsing of mail messages. |
| [netip](http://localhost:8001/pkg/net/netip/) | Package netip defines an IP address type that's a small value type. |
| [rpc](http://localhost:8001/pkg/net/rpc/) | Package rpc provides access to the exported methods of an object across a network or other I/O connection. |
| [jsonrpc](http://localhost:8001/pkg/net/rpc/jsonrpc/) | Package jsonrpc implements a JSON-RPC 1.0 ClientCodec and ServerCodec for the rpc package. |
| [smtp](http://localhost:8001/pkg/net/smtp/) | Package smtp implements the Simple Mail Transfer Protocol as defined in RFC 5321. |
| [textproto](http://localhost:8001/pkg/net/textproto/) | Package textproto implements generic support for text-based request/response protocols in the style of HTTP, NNTP, and SMTP. |
| [url](http://localhost:8001/pkg/net/url/) | Package url parses URLs and implements query escaping. |
| [os](http://localhost:8001/pkg/os/) | Package os provides a platform-independent interface to operating system functionality. |
| [exec](http://localhost:8001/pkg/os/exec/) | Package exec runs external commands. |
| [signal](http://localhost:8001/pkg/os/signal/) | Package signal implements access to incoming signals. |
| [user](http://localhost:8001/pkg/os/user/) | Package user allows user account lookups by name or id. |
| [path](http://localhost:8001/pkg/path/) | Package path implements utility routines for manipulating slash-separated paths. |
| [filepath](http://localhost:8001/pkg/path/filepath/) | Package filepath implements utility routines for manipulating filename paths in a way compatible with the target operating system-defined file paths. |
| [plugin](http://localhost:8001/pkg/plugin/) | Package plugin implements loading and symbol resolution of Go plugins. |
| [reflect](http://localhost:8001/pkg/reflect/) | Package reflect implements run-time reflection, allowing a program to manipulate objects with arbitrary types. |
| [regexp](http://localhost:8001/pkg/regexp/) | Package regexp implements regular expression search. |
| [syntax](http://localhost:8001/pkg/regexp/syntax/) | Package syntax parses regular expressions into parse trees and compiles parse trees into programs. |
| [runtime](http://localhost:8001/pkg/runtime/) | Package runtime contains operations that interact with Go's runtime system, such as functions to control goroutines. |
| [asan](http://localhost:8001/pkg/runtime/asan/) |  |
| [cgo](http://localhost:8001/pkg/runtime/cgo/) | Package cgo contains runtime support for code generated by the cgo tool. |
| [debug](http://localhost:8001/pkg/runtime/debug/) | Package debug contains facilities for programs to debug themselves while they are running. |
| [metrics](http://localhost:8001/pkg/runtime/metrics/) | Package metrics provides a stable interface to access implementation-defined metrics exported by the Go runtime. |
| [msan](http://localhost:8001/pkg/runtime/msan/) |  |
| [pprof](http://localhost:8001/pkg/runtime/pprof/) | Package pprof writes runtime profiling data in the format expected by the pprof visualization tool. |
| [race](http://localhost:8001/pkg/runtime/race/) | Package race implements data race detection logic. |
| [trace](http://localhost:8001/pkg/runtime/trace/) | Package trace contains facilities for programs to generate traces for the Go execution tracer. |
| [sort](http://localhost:8001/pkg/sort/) | Package sort provides primitives for sorting slices and user-defined collections. |
| [strconv](http://localhost:8001/pkg/strconv/) | Package strconv implements conversions to and from string representations of basic data types. |
| [strings](http://localhost:8001/pkg/strings/) | Package strings implements simple functions to manipulate UTF-8 encoded strings. |
| [sync](http://localhost:8001/pkg/sync/) | Package sync provides basic synchronization primitives such as mutual exclusion locks. |
| [atomic](http://localhost:8001/pkg/sync/atomic/) | Package atomic provides low-level atomic memory primitives useful for implementing synchronization algorithms. |
| [syscall](http://localhost:8001/pkg/syscall/) | Package syscall contains an interface to the low-level operating system primitives. |
| [js](http://localhost:8001/pkg/syscall/js/) | Package js gives access to the WebAssembly host environment when using the js/wasm architecture. |
| [testing](http://localhost:8001/pkg/testing/) | Package testing provides support for automated testing of Go packages. |
| [fstest](http://localhost:8001/pkg/testing/fstest/) | Package fstest implements support for testing implementations and users of file systems. |
| [iotest](http://localhost:8001/pkg/testing/iotest/) | Package iotest implements Readers and Writers useful mainly for testing. |
| [quick](http://localhost:8001/pkg/testing/quick/) | Package quick implements utility functions to help with black box testing. |
| [text](http://localhost:8001/pkg/text/) |  |
| [scanner](http://localhost:8001/pkg/text/scanner/) | Package scanner provides a scanner and tokenizer for UTF-8-encoded text. |
| [tabwriter](http://localhost:8001/pkg/text/tabwriter/) | Package tabwriter implements a write filter (tabwriter.Writer) that translates tabbed columns in input into properly aligned text. |
| [template](http://localhost:8001/pkg/text/template/) | Package template implements data-driven templates for generating textual output. |
| [parse](http://localhost:8001/pkg/text/template/parse/) | Package parse builds parse trees for templates as defined by text/template and html/template. |
| [time](http://localhost:8001/pkg/time/) | Package time provides functionality for measuring and displaying time. |
| [tzdata](http://localhost:8001/pkg/time/tzdata/) | Package tzdata provides an embedded copy of the timezone database. |
| [unicode](http://localhost:8001/pkg/unicode/) | Package unicode provides data and functions to test some properties of Unicode code points. |
| [utf16](http://localhost:8001/pkg/unicode/utf16/) | Package utf16 implements encoding and decoding of UTF-16 sequences. |
| [utf8](http://localhost:8001/pkg/unicode/utf8/) | Package utf8 implements functions and constants to support text encoded in UTF-8. |
| [unsafe](http://localhost:8001/pkg/unsafe/) | Package unsafe contains operations that step around the type safety of Go programs. |
