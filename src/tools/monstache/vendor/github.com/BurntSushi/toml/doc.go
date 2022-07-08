/*
Package toml implements decoding and encoding of TOML files.

This pakcage supports TOML v1.0.0, as listed on https://toml.io

There is also support for delaying decoding with the Primitive type, and
querying the set of keys in a TOML document with the MetaData type.

The sub-command github.com/BurntSushi/toml/cmd/tomlv can be used to verify
whether a file is a valid TOML document. It can also be used to print the type
of each key in a TOML document.
*/
package toml
