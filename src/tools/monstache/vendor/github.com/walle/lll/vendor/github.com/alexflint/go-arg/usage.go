package arg

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

// the width of the left column
const colWidth = 25

// Fail prints usage information to stderr and exits with non-zero status
func (p *Parser) Fail(msg string) {
	p.WriteUsage(os.Stderr)
	fmt.Fprintln(os.Stderr, "error:", msg)
	os.Exit(-1)
}

// WriteUsage writes usage information to the given writer
func (p *Parser) WriteUsage(w io.Writer) {
	var positionals, options []*spec
	for _, spec := range p.spec {
		if spec.positional {
			positionals = append(positionals, spec)
		} else {
			options = append(options, spec)
		}
	}

	fmt.Fprintf(w, "usage: %s", p.config.Program)

	// write the option component of the usage message
	for _, spec := range options {
		// prefix with a space
		fmt.Fprint(w, " ")
		if !spec.required {
			fmt.Fprint(w, "[")
		}
		fmt.Fprint(w, synopsis(spec, "--"+spec.long))
		if !spec.required {
			fmt.Fprint(w, "]")
		}
	}

	// write the positional component of the usage message
	for _, spec := range positionals {
		// prefix with a space
		fmt.Fprint(w, " ")
		up := strings.ToUpper(spec.long)
		if spec.multiple {
			fmt.Fprintf(w, "[%s [%s ...]]", up, up)
		} else {
			fmt.Fprint(w, up)
		}
	}
	fmt.Fprint(w, "\n")
}

// WriteHelp writes the usage string followed by the full help string for each option
func (p *Parser) WriteHelp(w io.Writer) {
	var positionals, options []*spec
	for _, spec := range p.spec {
		if spec.positional {
			positionals = append(positionals, spec)
		} else {
			options = append(options, spec)
		}
	}

	p.WriteUsage(w)

	// write the list of positionals
	if len(positionals) > 0 {
		fmt.Fprint(w, "\npositional arguments:\n")
		for _, spec := range positionals {
			left := "  " + spec.long
			fmt.Fprint(w, left)
			if spec.help != "" {
				if len(left)+2 < colWidth {
					fmt.Fprint(w, strings.Repeat(" ", colWidth-len(left)))
				} else {
					fmt.Fprint(w, "\n"+strings.Repeat(" ", colWidth))
				}
				fmt.Fprint(w, spec.help)
			}
			fmt.Fprint(w, "\n")
		}
	}

	// write the list of options
	fmt.Fprint(w, "\noptions:\n")
	for _, spec := range options {
		printOption(w, spec)
	}

	// write the list of built in options
	printOption(w, &spec{boolean: true, long: "help", short: "h", help: "display this help and exit"})
}

func printOption(w io.Writer, spec *spec) {
	left := "  " + synopsis(spec, "--"+spec.long)
	if spec.short != "" {
		left += ", " + synopsis(spec, "-"+spec.short)
	}
	fmt.Fprint(w, left)
	if spec.help != "" {
		if len(left)+2 < colWidth {
			fmt.Fprint(w, strings.Repeat(" ", colWidth-len(left)))
		} else {
			fmt.Fprint(w, "\n"+strings.Repeat(" ", colWidth))
		}
		fmt.Fprint(w, spec.help)
	}
	// If spec.dest is not the zero value then a default value has been added.
	v := spec.dest
	if v.IsValid() {
		z := reflect.Zero(v.Type())
		if (v.Type().Comparable() && z.Type().Comparable() && v.Interface() != z.Interface()) || v.Kind() == reflect.Slice && !v.IsNil() {
			fmt.Fprintf(w, " [default: %v]", v)
		}
	}
	fmt.Fprint(w, "\n")
}

func synopsis(spec *spec, form string) string {
	if spec.boolean {
		return form
	}
	return form + " " + strings.ToUpper(spec.long)
}
