<p align="center">
  <a href="https://pkg.go.dev/github.com/rymdport/go-qrcode" title="Go API Reference" rel="nofollow"><img src="https://img.shields.io/badge/go-documentation-blue.svg?style=flat" alt="Go API Reference"></a>
  [![Go Report Card](https://goreportcard.com/badge/github.com/rymdport/go-qrcode)](https://goreportcard.com/report/github.com/rymdport/go-qrcode)
  [![Analysis](https://github.com/rymdport/go-qrcode/actions/workflows/analysis.yml/badge.svg)](https://github.com/rymdport/go-qrcode/actions/workflows/analysis.yml)
  [![Tests](https://github.com/rymdport/go-qrcode/actions/workflows/tests.yml/badge.svg)](https://github.com/rymdport/go-qrcode/actions/workflows/tests.yml)
</p>

# go-qrcode

<img src='https://skip.org/img/nyancat-youtube-qr.png' align='right'>

Package qrcode implements a QR Code encoder.

A QR Code is a matrix (two-dimensional) barcode. Arbitrary content may be encoded, with URLs being a popular choice :)

Each QR Code contains error recovery information to aid reading damaged or obscured codes. There are four levels of error recovery: Low, medium, high and highest. QR Codes with a higher recovery level are more robust to damage, at the cost of being physically larger.

## Usage

```go
import qrcode "github.com/skip2/go-qrcode"
```

- **Create a 256x256 PNG image:**

```go
var png []byte
png, err := qrcode.Encode("https://example.org", qrcode.Medium, 256)
```

- **Create a 256x256 PNG image and write to a file:**

```go
err := qrcode.WriteFile("https://example.org", qrcode.Medium, 256, "qr.png")
```

- **Create a 256x256 PNG image with custom colors and write to file:**
```go
err := qrcode.WriteColorFile("https://example.org", qrcode.Medium, 256, color.Black, color.White, "qr.png")
```

All examples use the `qrcode.Medium` error Recovery Level and create a fixed 256x256px size QR Code. The last function creates a white on black instead of black on white QR Code.

## Maximum capacity
The maximum capacity of a QR Code varies according to the content encoded and the error recovery level. The maximum capacity is 2,953 bytes, 4,296 alphanumeric characters, 7,089 numeric digits, or a combination of these.

## Borderless QR Codes

To aid QR Code reading software, QR codes have a built in whitespace border.
If you know what you're doing, and don't want a border, you can set the `.DisableBorder` field on the `QRCode` type.
```go
code, err := qrcode.New("https://example.org", qrcode.Medium)
code.DisableBorder = true
```

## Links

- [http://en.wikipedia.org/wiki/QR_code](http://en.wikipedia.org/wiki/QR_code)
- [ISO/IEC 18004:2006](http://www.iso.org/iso/catalogue_detail.htm?csnumber=43655) - Main QR Code specification (approx CHF 198,00)<br>
