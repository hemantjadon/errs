package errs

import (
	"fmt"
	"runtime"
	"strings"
)

// FieldsError defines an error interface with an extra Fields method to get
// context fields associated with the error.
//
// Fields associated with an error must not be used for any logical deductions,
// they should only be used to make the error more dynamic.
//
// Fields of an error, like the Error string, should not be covered under api
// guarantees.
type FieldsError interface {
	error
	Fields() []Field
}

// LocationError defines an error interface with extra location method to get
// location source-code location at which error was created.
//
// Location associated with an error must not be used for any logical
// deductions, they should only be used for reporting purpose.
//
// Location of an error, like Error string, should not be covered under api
// guarantees.
type LocationError interface {
	error
	Location() (fn string, file string, line int)
}

// ChainError defines an error interface with an extra Chain method to get
// chain of errors resulting to the error.
//
// The chain of an error must not be used for any logical deductions, it should
// only be used to make the errors more visible.
//
// The Chain of an error, like the Error string, should not be covered under api
// guarantees.
type ChainError interface {
	error
	Chain() []error
}

// New creates a new error with the given message.
//
// If empty message is given, then nil error is returned.
func New(message string, fields ...Field) error {
	if len(message) == 0 {
		return nil
	}
	fdm := fundamental{msg: message, fields: fields, loc: getLocation(1)}
	return &fdm
}

// Wrap creates a new error with the given message wrapping the given error.
// When error is wrapped it can be unwrapped to get the underlying errors.
//
// If the given error is nil, then nil error is returned.
func Wrap(err error, message string, fields ...Field) error {
	if err == nil {
		return nil
	}
	fdm := fundamental{msg: message, fields: fields, loc: getLocation(1)}
	var chn []error
	chn = append(chn, &fdm)
	switch e := err.(type) {
	case ChainError:
		chn = append(chn, e.Chain()...)
	default:
		chn = append(chn, e)
	}
	wrp := wrapping{fundamental: &fdm, err: err, wrapped: true, chain: chn}
	return &wrp
}

// Box creates a new error with the given message boxing the given error.
// When error is boxed it cannot be unwrapped to get the underlying errors.
//
// If the given error is nil, then nil error is returned.
func Box(err error, message string, fields ...Field) error {
	if err == nil {
		return nil
	}
	fdm := fundamental{msg: message, fields: fields, loc: getLocation(1)}
	var chn []error
	chn = append(chn, &fdm)
	switch e := err.(type) {
	case ChainError:
		chn = append(chn, e.Chain()...)
	default:
		chn = append(chn, e)
	}
	wrp := wrapping{fundamental: &fdm, err: err, wrapped: false, chain: chn}
	return &wrp
}

type fundamental struct {
	msg    string
	fields []Field
	loc    location
}

func (f fundamental) Error() string {
	if len(f.msg) == 0 {
		return ""
	}
	if len(f.fields) == 0 {
		return f.msg
	}
	return fmt.Sprintf("%s (%s)", f.msg, f.fieldsString())
}

func (f fundamental) fieldsString() string {
	var sb strings.Builder
	for idx, field := range f.fields {
		sb.WriteString(field.Key())
		sb.WriteString("=")
		sb.WriteString(valueString(field))
		if idx < len(f.fields)-1 {
			sb.WriteString(" ")
		}
	}
	return sb.String()
}

// Fields gives the fields associated with the error.
func (f fundamental) Fields() []Field {
	fields := make([]Field, 0, len(f.fields))
	fields = append(fields, f.fields...)
	return fields
}

// Location gives function name, file name and line number of the location
// where error was created.
func (f fundamental) Location() (fn string, file string, line int) {
	return f.loc.Function, f.loc.File, f.loc.Line
}

type wrapping struct {
	*fundamental
	chain   []error
	err     error
	wrapped bool
}

// Chain gives the chain of errors associated with the error.
func (w wrapping) Chain() []error {
	stk := make([]error, 0, len(w.chain))
	stk = append(stk, w.chain...)
	return stk
}

func (w wrapping) Error() string {
	fes := w.fundamental.Error()
	if len(fes) == 0 && w.err == nil {
		return ""
	}
	if len(fes) != 0 && w.err == nil {
		return fes
	}
	if len(fes) == 0 && w.err != nil {
		return w.err.Error()
	}
	return fmt.Sprintf("%s: %s", fes, w.err.Error())
}

// Unwrap unwraps the error giving the underlying error. If error does not wrap
// any error then nil is returned.
func (w wrapping) Unwrap() error {
	if !w.wrapped {
		return nil
	}
	return w.err
}

type location struct {
	Function string
	File     string
	Line     int
}

func getLocation(skip int) location {
	pc, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return location{}
	}
	fn := runtime.FuncForPC(pc)
	return location{
		Function: fn.Name(),
		File:     file,
		Line:     line,
	}
}
