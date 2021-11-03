# gokrb5

It is recommended to use the latest version: [![Version](https://img.shields.io/github/release/jcmturner/gokrb5.svg)](https://github.com/jcmturner/gokrb5/releases)

Development will be focused on the latest major version. New features will only be targeted at this version.

| Versions | Dependency Management | Import Path | Usage | Godoc | Go Report Card |
|----------|-----------------------|-------------|-------|-------|----------------|
| [![v8](https://github.com/jcmturner/gokrb5/workflows/v8/badge.svg)](https://github.com/jcmturner/gokrb5/actions?query=workflow%3Av8) | Go modules | import "github.com/jcmturner/gokrb5/v8/{sub-package}" | [![Usage](https://img.shields.io/badge/v8-usage-blue)](https://github.com/jcmturner/gokrb5/blob/master/v8/USAGE.md) | [![GoDoc](https://img.shields.io/badge/godoc-reference-blue)](https://pkg.go.dev/github.com/jcmturner/gokrb5/v8) | [![Go Report Card](https://goreportcard.com/badge/github.com/jcmturner/gokrb5/v8)](https://goreportcard.com/report/github.com/jcmturner/gokrb5/v8) |
| [![v7](https://github.com/jcmturner/gokrb5/workflows/v7/badge.svg)](https://github.com/jcmturner/gokrb5/actions?query=workflow%3Av7) | gopkg.in | import "gopkg.in/jcmturner/gokrb5.v7/{sub-package}" | [![Usage](https://img.shields.io/badge/v7-usage-blue)](https://github.com/jcmturner/gokrb5/blob/master/USAGE.md) | [![GoDoc](https://img.shields.io/badge/godoc-reference-blue)](https://pkg.go.dev/github.com/jcmturner/gokrb5@v7.5.0+incompatible) | [![Go Report Card](https://goreportcard.com/badge/gopkg.in/jcmturner/gokrb5.v7)](https://goreportcard.com/report/gopkg.in/jcmturner/gokrb5.v7) |


#### Go Version Support
![Go version](https://img.shields.io/badge/Go-1.15-brightgreen.svg)
![Go version](https://img.shields.io/badge/Go-1.14-brightgreen.svg)

gokrb5 may work with other versions of Go but they are not formally tested.
It has been reported that gokrb5 also works with the [gollvm](https://go.googlesource.com/gollvm/) compiler but this is not formally tested.

## Features
* **Pure Go** - no dependency on external libraries 
* No platform specific code
* Server Side
  * HTTP handler wrapper implements SPNEGO Kerberos authentication
  * HTTP handler wrapper decodes Microsoft AD PAC authorization data
* Client Side
  * Client that can authenticate to an SPNEGO Kerberos authenticated web service
  * Ability to change client's password
* General
  * Kerberos libraries for custom integration
  * Parsing Keytab files
  * Parsing krb5.conf files
  * Parsing client credentials cache files such as `/tmp/krb5cc_$(id -u $(whoami))`

#### Implemented Encryption & Checksum Types

| Implementation | Encryption ID | Checksum ID | RFC |
|-------|-------------|------------|------|
| des3-cbc-sha1-kd | 16 | 12 | 3961 |
| aes128-cts-hmac-sha1-96 | 17 | 15 | 3962 |
| aes256-cts-hmac-sha1-96 | 18 | 16 | 3962 |
| aes128-cts-hmac-sha256-128 | 19 | 19 | 8009 |
| aes256-cts-hmac-sha384-192 | 20 | 20 | 8009 |
| rc4-hmac | 23 | -138 | 4757 |


The following is working/tested:
* Tested against MIT KDC (1.6.3 is the oldest version tested against) and Microsoft Active Directory (Windows 2008 R2)
* Tested against a KDC that supports PA-FX-FAST.
* Tested against users that have pre-authentication required using PA-ENC-TIMESTAMP.
* Microsoft PAC Authorization Data is processed and exposed in the HTTP request context. Available if Microsoft Active Directory is used as the KDC.

## Contributing
If you are interested in contributing to gokrb5, great! Please read the [contribution guidelines](https://github.com/jcmturner/gokrb5/blob/master/CONTRIBUTING.md).

---

## References
* [RFC 3244 Microsoft Windows 2000 Kerberos Change Password and Set Password Protocols](https://tools.ietf.org/html/rfc3244)
* [RFC 4120 The Kerberos Network Authentication Service (V5)](https://tools.ietf.org/html/rfc4120)
* [RFC 3961 Encryption and Checksum Specifications for Kerberos 5](https://tools.ietf.org/html/rfc3961)
* [RFC 3962 Advanced Encryption Standard (AES) Encryption for Kerberos 5](https://tools.ietf.org/html/rfc3962)
* [RFC 4121 The Kerberos Version 5 GSS-API Mechanism](https://tools.ietf.org/html/rfc4121)
* [RFC 4178 The Simple and Protected Generic Security Service Application Program Interface (GSS-API) Negotiation Mechanism](https://tools.ietf.org/html/rfc4178.html)
* [RFC 4559 SPNEGO-based Kerberos and NTLM HTTP Authentication in Microsoft Windows](https://tools.ietf.org/html/rfc4559.html)
* [RFC 4757 The RC4-HMAC Kerberos Encryption Types Used by Microsoft Windows](https://tools.ietf.org/html/rfc4757)
* [RFC 6806 Kerberos Principal Name Canonicalization and Cross-Realm Referrals](https://tools.ietf.org/html/rfc6806.html)
* [RFC 6113 A Generalized Framework for Kerberos Pre-Authentication](https://tools.ietf.org/html/rfc6113.html)
* [RFC 8009 AES Encryption with HMAC-SHA2 for Kerberos 5](https://tools.ietf.org/html/rfc8009)
* [IANA Assigned Kerberos Numbers](http://www.iana.org/assignments/kerberos-parameters/kerberos-parameters.xhtml)
* [HTTP-Based Cross-Platform Authentication by Using the Negotiate Protocol - Part 1](https://msdn.microsoft.com/en-us/library/ms995329.aspx)
* [HTTP-Based Cross-Platform Authentication by Using the Negotiate Protocol - Part 2](https://msdn.microsoft.com/en-us/library/ms995330.aspx)
* [Microsoft PAC Validation](https://blogs.msdn.microsoft.com/openspecification/2009/04/24/understanding-microsoft-kerberos-pac-validation/)
* [Microsoft Kerberos Protocol Extensions](https://msdn.microsoft.com/en-us/library/cc233855.aspx)
* [Windows Data Types](https://msdn.microsoft.com/en-us/library/cc230273.aspx)

### Useful Links
* https://en.wikipedia.org/wiki/Ciphertext_stealing#CBC_ciphertext_stealing

## Thanks
* Greg Hudson from the MIT Consortium for Kerberos and Internet Trust for providing useful advice.

## Contributing
Thank you for your interest in contributing to gokrb5 please read the 
[contribution guide](https://github.com/jcmturner/gokrb5/blob/master/CONTRIBUTING.md) as it should help you get started.

## Known Issues
| Issue | Worked around? | References |
|-------|-------------|------------|
| The Go standard library's encoding/asn1 package cannot unmarshal into slice of asn1.RawValue | Yes | https://github.com/golang/go/issues/17321 |
| The Go standard library's encoding/asn1 package cannot marshal into a GeneralString | Yes - using https://github.com/jcmturner/gofork/tree/master/encoding/asn1 | https://github.com/golang/go/issues/18832 |
| The Go standard library's encoding/asn1 package cannot marshal into slice of strings and pass stringtype parameter tags to members | Yes - using https://github.com/jcmturner/gofork/tree/master/encoding/asn1 | https://github.com/golang/go/issues/18834 |
| The Go standard library's encoding/asn1 package cannot marshal with application tags | Yes | |
| The Go standard library's x/crypto/pbkdf2.Key function uses the int type for iteraction count limiting meaning the 4294967296 count specified in https://tools.ietf.org/html/rfc3962 section 4 cannot be met on 32bit systems | Yes - using https://github.com/jcmturner/gofork/tree/master/x/crypto/pbkdf2 | https://go-review.googlesource.com/c/crypto/+/85535 |
