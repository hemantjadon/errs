package errs_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hemantjadon/errs"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("empty message", func(t *testing.T) {
		t.Parallel()

		msg := ""
		err := errs.New(msg)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})

	t.Run("simple message", func(t *testing.T) {
		t.Parallel()

		msg := "error occurred"
		err := errs.New(msg)

		if !strings.Contains(err.Error(), msg) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), msg)
		}
	})

	t.Run("message with fields", func(t *testing.T) {
		t.Parallel()

		msg := "error occurred"
		field1, value1, field2, value2 := "field1", "value1", "field2", "value2"
		err := errs.New(msg, errs.F(field1, value1), errs.F(field2, value2))

		if !strings.Contains(err.Error(), msg) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), msg)
		}
		if !strings.Contains(err.Error(), field1) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), field1)
		}
		if !strings.Contains(err.Error(), value1) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), value1)
		}
		if !strings.Contains(err.Error(), field2) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), field2)
		}
		if !strings.Contains(err.Error(), value2) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), value2)
		}
	})

	t.Run("different field types", func(t *testing.T) {
		t.Parallel()

		msg := "error occurred"
		field1, value1 := "field1", "value1"
		field2, value2 := "field2", 1234
		field3, value3 := "field3", 12.34
		field4, value4 := "field4", true
		tt1 := time.Now()
		field5, value5 := "field5", tt1
		tt2 := time.Now().Add(1 * time.Hour)
		field6, value6 := "field6", &tt2
		field7 := "field7"

		err := errs.New(msg, errs.F(field1, value1), errs.F(field2, value2), errs.F(field3, value3), errs.F(field4, value4), errs.F(field5, value5), errs.F(field6, value6), errs.F(field7, nil))
		if !strings.Contains(err.Error(), msg) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), msg)
		}
		if !strings.Contains(err.Error(), field1) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), field1)
		}
		if !strings.Contains(err.Error(), value1) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), value1)
		}
		if !strings.Contains(err.Error(), field2) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), field2)
		}
		if !strings.Contains(err.Error(), fmt.Sprintf("%v", value2)) {
			t.Fatalf("Error(): got = '%s', want contains = '%v'", err.Error(), value2)
		}
		if !strings.Contains(err.Error(), field3) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), field3)
		}
		if !strings.Contains(err.Error(), fmt.Sprintf("%v", value3)) {
			t.Fatalf("Error(): got = '%s', want contains = '%v'", err.Error(), value3)
		}
		if !strings.Contains(err.Error(), field4) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), field4)
		}
		if !strings.Contains(err.Error(), fmt.Sprintf("%v", value4)) {
			t.Fatalf("Error(): got = '%s', want contains = '%v'", err.Error(), value4)
		}
		if !strings.Contains(err.Error(), field5) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), field5)
		}
		if !strings.Contains(err.Error(), tt1.Format(time.RFC3339)) {
			t.Fatalf("Error(): got = '%s', want contains = '%v'", err.Error(), tt1.Format(time.RFC3339))
		}
		if !strings.Contains(err.Error(), field6) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), field6)
		}
		if !strings.Contains(err.Error(), tt2.Format(time.RFC3339)) {
			t.Fatalf("Error(): got = '%s', want contains = '%v'", err.Error(), tt2.Format(time.RFC3339))
		}
		if !strings.Contains(err.Error(), field7) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), field7)
		}
		if !strings.Contains(err.Error(), fmt.Sprintf("%v", nil)) {
			t.Fatalf("Error(): got = '%s', want contains = '%v'", err.Error(), nil)
		}
	})

	t.Run("FieldsError", func(t *testing.T) {
		t.Parallel()

		msg := "error occurred"
		field1, value1, field2, value2 := "field1", "value1", "field2", "value2"
		err := errs.New(msg, errs.F(field1, value1), errs.F(field2, value2))

		ferr, ok := err.(errs.FieldsError)
		if !ok {
			t.Fatalf("got type = '%T', want = 'FieldsError'", err)
		}

		fields := ferr.Fields()
		if len(fields) != 2 {
			t.Fatalf("Fields(): got = '%d', want = '2'", len(fields))
		}
		if fields[0].Key() != field1 {
			t.Fatalf("fields[0].Key(): got = '%s', want = '%s'", fields[0].Key(), field1)
		}
		if fields[0].Value() != value1 {
			t.Fatalf("fields[0].Value(): got = '%s', want = '%s'", fields[0].Value(), value1)
		}
		if fields[1].Key() != field2 {
			t.Fatalf("fields[1].Key(): got = '%s', want = '%s'", fields[1].Key(), field2)
		}
		if fields[1].Value() != value2 {
			t.Fatalf("fields[1].Value(): got = '%s', want = '%s'", fields[1].Value(), value2)
		}
	})

	t.Run("LocationError", func(t *testing.T) {
		t.Parallel()

		msg := "error occurred"
		err := errs.New(msg)

		lerr, ok := err.(errs.LocationError)
		if !ok {
			t.Fatalf("got type = '%T', want = 'LocationError'", err)
		}

		fn, file, line := lerr.Location()
		if len(fn) == 0 {
			t.Fatalf("got fn = '%s', want = '%s'", fn, "non-empty")
		}
		if len(file) == 0 {
			t.Fatalf("got file = '%s', want = '%s'", file, "non-empty")
		}
		if line == 0 {
			t.Fatalf("got line = '%d', want = '%s'", line, "non-zero")
		}
	})
}

func TestWrap(t *testing.T) {
	t.Parallel()

	baseField, baseValue := "base_field", "base_value"
	baseErr := errs.New("base error", errs.F(baseField, baseValue))

	t.Run("nil error", func(t *testing.T) {
		t.Parallel()

		msg := "error occurred"
		err := errs.Wrap(nil, msg)

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})

	t.Run("empty message", func(t *testing.T) {
		t.Parallel()

		msg := ""
		err := errs.Wrap(baseErr, msg)
		if err == nil {
			t.Fatalf("expected non-nil, got %v", err)
		}
		if !strings.Contains(err.Error(), baseErr.Error()) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), baseErr.Error())
		}
		if !strings.Contains(err.Error(), baseField) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), baseField)
		}
		if !strings.Contains(err.Error(), baseValue) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), baseValue)
		}
	})

	t.Run("simple message", func(t *testing.T) {
		t.Parallel()

		msg := "error occurred"
		err := errs.Wrap(baseErr, msg)

		if !strings.Contains(err.Error(), msg) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), msg)
		}
		if !strings.Contains(err.Error(), baseErr.Error()) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), baseErr.Error())
		}
		if !strings.Contains(err.Error(), baseField) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), baseField)
		}
		if !strings.Contains(err.Error(), baseValue) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), baseValue)
		}
		if !strings.Contains(err.Error(), msg) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), msg)
		}
	})

	t.Run("message with fields", func(t *testing.T) {
		t.Parallel()

		msg := "error occurred"
		field1, value1, field2, value2 := "field1", "value1", "field2", "value2"
		err := errs.Wrap(baseErr, msg, errs.F(field1, value1), errs.F(field2, value2))

		if !strings.Contains(err.Error(), baseErr.Error()) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), baseErr.Error())
		}
		if !strings.Contains(err.Error(), baseField) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), baseField)
		}
		if !strings.Contains(err.Error(), baseValue) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), baseValue)
		}
		if !strings.Contains(err.Error(), msg) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), msg)
		}
		if !strings.Contains(err.Error(), field1) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), field1)
		}
		if !strings.Contains(err.Error(), value1) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), value1)
		}
		if !strings.Contains(err.Error(), field2) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), field2)
		}
		if !strings.Contains(err.Error(), value2) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), value2)
		}
	})

	t.Run("FieldsError", func(t *testing.T) {
		t.Parallel()

		msg := "error occurred"
		field1, value1, field2, value2 := "field1", "value1", "field2", "value2"
		err := errs.Wrap(baseErr, msg, errs.F(field1, value1), errs.F(field2, value2))

		ferr, ok := err.(errs.FieldsError)
		if !ok {
			t.Fatalf("got type = '%T', want = 'FieldsError'", err)
		}

		fields := ferr.Fields()
		if len(fields) != 2 {
			t.Fatalf("Fields(): got = '%d', want = '2'", len(fields))
		}
		if fields[0].Key() != field1 {
			t.Fatalf("fields[0].Key(): got = '%s', want = '%s'", fields[0].Key(), field1)
		}
		if fields[0].Value() != value1 {
			t.Fatalf("fields[0].Value(): got = '%s', want = '%s'", fields[0].Value(), value1)
		}
		if fields[1].Key() != field2 {
			t.Fatalf("fields[1].Key(): got = '%s', want = '%s'", fields[1].Key(), field2)
		}
		if fields[1].Value() != value2 {
			t.Fatalf("fields[1].Value(): got = '%s', want = '%s'", fields[1].Value(), value2)
		}
	})

	t.Run("LocationError", func(t *testing.T) {
		t.Parallel()

		msg := "error occurred"
		err := errs.Wrap(baseErr, msg)

		lerr, ok := err.(errs.LocationError)
		if !ok {
			t.Fatalf("got type = '%T', want = 'LocationError'", err)
		}

		fn, file, line := lerr.Location()
		if len(fn) == 0 {
			t.Fatalf("got fn = '%s', want = '%s'", fn, "non-empty")
		}
		if len(file) == 0 {
			t.Fatalf("got file = '%s', want = '%s'", file, "non-empty")
		}
		if line == 0 {
			t.Fatalf("got line = '%d', want = '%s'", line, "non-zero")
		}
	})

	t.Run("ChainError", func(t *testing.T) {
		t.Parallel()

		msg1 := "error one"
		msg2 := "error two"
		msg3 := "error three"
		msg4 := "error four"

		err1 := errs.New(msg1)
		err2 := errs.Wrap(err1, msg2)
		err3 := errs.Wrap(err2, msg3)
		err4 := errs.Wrap(err3, msg4)

		msg := "error occurred"
		err := errs.Wrap(err4, msg)

		cerr, ok := err.(errs.ChainError)
		if !ok {
			t.Fatalf("got type = '%T', want = 'ChainError'", err)
		}

		chain := cerr.Chain()
		if len(chain) != 5 {
			t.Fatalf("len(chain): got = %d, want = %d", len(chain), 5)
		}
		if chain[0].Error() != msg {
			t.Fatalf("chain[0].Error(): got = '%s', want = '%s'", chain[0].Error(), msg)
		}
		if chain[1].Error() != msg4 {
			t.Fatalf("chain[1].Error(): got = '%s', want = '%s'", chain[1].Error(), msg4)
		}
		if chain[2].Error() != msg3 {
			t.Fatalf("chain[2].Error(): got = '%s', want = '%s'", chain[2].Error(), msg3)
		}
		if chain[3].Error() != msg2 {
			t.Fatalf("chain[3].Error(): got = '%s', want = '%s'", chain[3].Error(), msg2)
		}
		if chain[4].Error() != msg1 {
			t.Fatalf("chain[4].Error(): got = '%s', want = '%s'", chain[4].Error(), msg1)
		}
	})

	t.Run("Unwrap", func(t *testing.T) {
		t.Parallel()

		msg := "error occurred"
		err := errs.Wrap(baseErr, msg)

		if !errors.Is(err, baseErr) {
			t.Fatalf("should wrap base error")
		}
	})
}

func TestBox(t *testing.T) {
	t.Parallel()

	baseField, baseValue := "base_field", "base_value"
	baseErr := errs.New("base error", errs.F(baseField, baseValue))

	t.Run("nil error", func(t *testing.T) {
		t.Parallel()

		msg := "error occurred"
		err := errs.Box(nil, msg)

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})

	t.Run("empty message", func(t *testing.T) {
		t.Parallel()

		msg := ""
		err := errs.Box(baseErr, msg)
		if err == nil {
			t.Fatalf("expected non-nil, got %v", err)
		}
		if !strings.Contains(err.Error(), baseErr.Error()) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), baseErr.Error())
		}
		if !strings.Contains(err.Error(), baseField) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), baseField)
		}
		if !strings.Contains(err.Error(), baseValue) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), baseValue)
		}
	})

	t.Run("simple message", func(t *testing.T) {
		t.Parallel()

		msg := "error occurred"
		err := errs.Box(baseErr, msg)

		if !strings.Contains(err.Error(), msg) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), msg)
		}
		if !strings.Contains(err.Error(), baseErr.Error()) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), baseErr.Error())
		}
		if !strings.Contains(err.Error(), baseField) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), baseField)
		}
		if !strings.Contains(err.Error(), baseValue) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), baseValue)
		}
		if !strings.Contains(err.Error(), msg) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), msg)
		}
	})

	t.Run("message with fields", func(t *testing.T) {
		t.Parallel()

		msg := "error occurred"
		field1, value1, field2, value2 := "field1", "value1", "field2", "value2"
		err := errs.Box(baseErr, msg, errs.F(field1, value1), errs.F(field2, value2))

		if !strings.Contains(err.Error(), baseErr.Error()) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), baseErr.Error())
		}
		if !strings.Contains(err.Error(), baseField) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), baseField)
		}
		if !strings.Contains(err.Error(), baseValue) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), baseValue)
		}
		if !strings.Contains(err.Error(), msg) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), msg)
		}
		if !strings.Contains(err.Error(), field1) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), field1)
		}
		if !strings.Contains(err.Error(), value1) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), value1)
		}
		if !strings.Contains(err.Error(), field2) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), field2)
		}
		if !strings.Contains(err.Error(), value2) {
			t.Fatalf("Error(): got = '%s', want contains = '%s'", err.Error(), value2)
		}
	})

	t.Run("FieldsError", func(t *testing.T) {
		t.Parallel()

		msg := "error occurred"
		field1, value1, field2, value2 := "field1", "value1", "field2", "value2"
		err := errs.Box(baseErr, msg, errs.F(field1, value1), errs.F(field2, value2))

		ferr, ok := err.(errs.FieldsError)
		if !ok {
			t.Fatalf("got type = '%T', want = 'FieldsError'", err)
		}

		fields := ferr.Fields()
		if len(fields) != 2 {
			t.Fatalf("Fields(): got = '%d', want = '2'", len(fields))
		}
		if fields[0].Key() != field1 {
			t.Fatalf("fields[0].Key(): got = '%s', want = '%s'", fields[0].Key(), field1)
		}
		if fields[0].Value() != value1 {
			t.Fatalf("fields[0].Value(): got = '%s', want = '%s'", fields[0].Value(), value1)
		}
		if fields[1].Key() != field2 {
			t.Fatalf("fields[1].Key(): got = '%s', want = '%s'", fields[1].Key(), field2)
		}
		if fields[1].Value() != value2 {
			t.Fatalf("fields[1].Value(): got = '%s', want = '%s'", fields[1].Value(), value2)
		}
	})

	t.Run("LocationError", func(t *testing.T) {
		t.Parallel()

		msg := "error occurred"
		err := errs.Box(baseErr, msg)

		lerr, ok := err.(errs.LocationError)
		if !ok {
			t.Fatalf("got type = '%T', want = 'LocationError'", err)
		}

		fn, file, line := lerr.Location()
		if len(fn) == 0 {
			t.Fatalf("got fn = '%s', want = '%s'", fn, "non-empty")
		}
		if len(file) == 0 {
			t.Fatalf("got file = '%s', want = '%s'", file, "non-empty")
		}
		if line == 0 {
			t.Fatalf("got line = '%d', want = '%s'", line, "non-zero")
		}
	})

	t.Run("ChainError", func(t *testing.T) {
		t.Parallel()

		msg1 := "error one"
		msg2 := "error two"
		msg3 := "error three"
		msg4 := "error four"

		err1 := errs.New(msg1)
		err2 := errs.Box(err1, msg2)
		err3 := errs.Box(err2, msg3)
		err4 := errs.Box(err3, msg4)

		msg := "error occurred"
		err := errs.Box(err4, msg)

		cerr, ok := err.(errs.ChainError)
		if !ok {
			t.Fatalf("got type = '%T', want = 'ChainError'", err)
		}

		chain := cerr.Chain()
		if len(chain) != 5 {
			t.Fatalf("len(chain): got = %d, want = %d", len(chain), 5)
		}
		if chain[0].Error() != msg {
			t.Fatalf("chain[0].Error(): got = '%s', want = '%s'", chain[0].Error(), msg)
		}
		if chain[1].Error() != msg4 {
			t.Fatalf("chain[1].Error(): got = '%s', want = '%s'", chain[1].Error(), msg4)
		}
		if chain[2].Error() != msg3 {
			t.Fatalf("chain[2].Error(): got = '%s', want = '%s'", chain[2].Error(), msg3)
		}
		if chain[3].Error() != msg2 {
			t.Fatalf("chain[3].Error(): got = '%s', want = '%s'", chain[3].Error(), msg2)
		}
		if chain[4].Error() != msg1 {
			t.Fatalf("chain[4].Error(): got = '%s', want = '%s'", chain[4].Error(), msg1)
		}
	})

	t.Run("Unwrap", func(t *testing.T) {
		t.Parallel()

		msg := "error occurred"
		err := errs.Box(baseErr, msg)

		if errors.Is(err, baseErr) {
			t.Fatalf("should not wrap base error")
		}
	})
}
