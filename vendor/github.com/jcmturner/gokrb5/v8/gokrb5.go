/*
Package gokrb5 provides a Kerberos 5 implementation for Go.

This is a pure Go implementation and does not have dependencies on native libraries.

Feature include:

Server Side

HTTP handler wrapper implements SPNEGO Kerberos authentication.

HTTP handler wrapper decodes Microsoft AD PAC authorization data.

Client Side

Client that can authenticate to an SPNEGO Kerberos authenticated web service.

Ability to change client's password.

General

Kerberos libraries for custom integration.

Parsing Keytab files.

Parsing krb5.conf files.
*/
package gokrb5
