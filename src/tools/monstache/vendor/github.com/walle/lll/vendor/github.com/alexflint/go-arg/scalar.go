package arg

import (
	"encoding"
	"fmt"
	"net"
	"net/mail"
	"reflect"
	"strconv"
	"time"
)

// The reflected form of some special types
var (
	textUnmarshalerType = reflect.TypeOf([]encoding.TextUnmarshaler{}).Elem()
	durationType        = reflect.TypeOf(time.Duration(0))
	mailAddressType     = reflect.TypeOf(mail.Address{})
	ipType              = reflect.TypeOf(net.IP{})
	macType             = reflect.TypeOf(net.HardwareAddr{})
)

// isScalar returns true if the type can be parsed from a single string
func isScalar(t reflect.Type) (scalar, boolean bool) {
	// If it implements encoding.TextUnmarshaler then use that
	if t.Implements(textUnmarshalerType) {
		// scalar=YES, boolean=NO
		return true, false
	}

	// If we have a pointer then dereference it
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Check for other special types
	switch t {
	case durationType, mailAddressType, ipType, macType:
		// scalar=YES, boolean=NO
		return true, false
	}

	// Fall back to checking the kind
	switch t.Kind() {
	case reflect.Bool:
		// scalar=YES, boolean=YES
		return true, true
	case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		// scalar=YES, boolean=NO
		return true, false
	}
	// scalar=NO, boolean=NO
	return false, false
}

// set a value from a string
func setScalar(v reflect.Value, s string) error {
	if !v.CanSet() {
		return fmt.Errorf("field is not exported")
	}

	// If we have a nil pointer then allocate a new object
	if v.Kind() == reflect.Ptr && v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}

	// Get the object as an interface
	scalar := v.Interface()

	// If it implements encoding.TextUnmarshaler then use that
	if scalar, ok := scalar.(encoding.TextUnmarshaler); ok {
		return scalar.UnmarshalText([]byte(s))
	}

	// If we have a pointer then dereference it
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Switch on concrete type
	switch scalar.(type) {
	case time.Duration:
		duration, err := time.ParseDuration(s)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(duration))
		return nil
	case mail.Address:
		addr, err := mail.ParseAddress(s)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(*addr))
		return nil
	case net.IP:
		ip := net.ParseIP(s)
		if ip == nil {
			return fmt.Errorf(`invalid IP address: "%s"`, s)
		}
		v.Set(reflect.ValueOf(ip))
		return nil
	case net.HardwareAddr:
		ip, err := net.ParseMAC(s)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(ip))
		return nil
	}

	// Switch on kind so that we can handle derived types
	switch v.Kind() {
	case reflect.String:
		v.SetString(s)
	case reflect.Bool:
		x, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		v.SetBool(x)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x, err := strconv.ParseInt(s, 10, v.Type().Bits())
		if err != nil {
			return err
		}
		v.SetInt(x)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		x, err := strconv.ParseUint(s, 10, v.Type().Bits())
		if err != nil {
			return err
		}
		v.SetUint(x)
	case reflect.Float32, reflect.Float64:
		x, err := strconv.ParseFloat(s, v.Type().Bits())
		if err != nil {
			return err
		}
		v.SetFloat(x)
	default:
		return fmt.Errorf("cannot parse argument into %s", v.Type().String())
	}
	return nil
}
