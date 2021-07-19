package cerrors_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/1debit/cerrors"
)

func TestUnwrapMessage(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		result string
	}{
		{
			name:   "at the end with a colon",
			err:    fmt.Errorf("one: %w", fmt.Errorf("two: %w", errors.New("three"))),
			result: "one",
		},
		{
			name:   "at the end without colon",
			err:    fmt.Errorf("one %w", fmt.Errorf("two: %w", errors.New("three"))),
			result: "one",
		},
		{
			name:   "in the middle",
			err:    fmt.Errorf("foo %w bar", fmt.Errorf("two: %w", errors.New("three"))),
			result: "foo * bar",
		},
		{
			name:   "chained without wrapping",
			err:    fmt.Errorf("one: %v", fmt.Errorf("two: %w", errors.New("three"))),
			result: "one: two: three",
		},
		{
			name:   "not chained",
			err:    errors.New("foo bar"),
			result: "foo bar",
		},
		{
			name:   "empty",
			err:    errors.New(""),
			result: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.result, cerrors.UnwrapMessage(test.err))
		})
	}
}

func TestCerrorFormat(t *testing.T) {
	var (
		location         = `    \S*error_test.go:\d+`
		locationIndented = `        \S*error_test.go:\d+`
		function         = `    github.com/1debit/cerrors_test.TestCerrorFormat.*`
	)

	tests := []struct {
		name   string
		chain  []func(error) error
		verb   string
		root   error
		result []string
	}{
		{
			name:   "single error, short",
			root:   cerrors.New("foo"),
			verb:   "%v",
			result: []string{"foo"},
		},
		{
			name: "multiple levels wrapped, short",
			root: cerrors.New("one"),
			chain: []func(error) error{
				func(e error) error { return cerrors.Wrap(e, "two") },
				func(e error) error { return cerrors.Wrap(e, "three") },
			},
			verb:   "%v",
			result: []string{"three: two: one"},
		},
		{
			name: "including non-cerror, short",
			root: cerrors.New("one"),
			chain: []func(error) error{
				func(e error) error { return fmt.Errorf("two: %w", e) },
				func(e error) error { return cerrors.Wrap(e, "three") },
			},
			verb:   "%v",
			result: []string{"three: two: one"},
		},
		{
			name: "single error, detailed",
			root: cerrors.New("foo"),
			verb: "%+v",
			result: []string{
				"foo",
				location,
			},
		},
		{
			name: "multiple levels wrapped, detailed",
			root: cerrors.New("one"),
			chain: []func(error) error{
				func(e error) error { return cerrors.Wrap(e, "two") },
				func(e error) error { return cerrors.Wrap(e, "three") },
			},
			verb: "%+v",
			result: []string{
				"three",
				location,
				"  - two",
				location,
				"  - one",
				location,
			},
		},
		{
			name: "including non-cerror, detailed",
			root: cerrors.New("one"),
			chain: []func(error) error{
				func(e error) error { return fmt.Errorf("two: %w", e) },
				func(e error) error { return cerrors.Wrap(e, "three") },
			},
			verb: "%+v",
			result: []string{
				"three",
				location,
				"  - two",
				"  - one",
				location,
			},
		},
		{
			name: "including non-cerror, detailed with functions",
			root: cerrors.New("one"),
			chain: []func(error) error{
				func(e error) error { return fmt.Errorf("two: %w", e) },
				func(e error) error { return cerrors.Wrap(e, "three") },
			},
			verb: "%#v",
			result: []string{
				"three",
				function,
				locationIndented,
				"  - two",
				"  - one",
				function,
				locationIndented,
			},
		},
		{
			name:   "invalid verb",
			root:   cerrors.New("one"),
			verb:   "%d",
			result: []string{`%!d\(cerror\)`},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.root
			for _, f := range test.chain {
				err = f(err)
			}

			result := strings.Split(strings.TrimSpace(fmt.Sprintf(test.verb, err)), "\n")
			require.Len(t, result, len(test.result))
			for idx, line := range test.result {
				assert.Regexp(t, "^"+line+"$", result[idx], "at line %d", idx)
			}
		})
	}
}
