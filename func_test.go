package argmapper

import (
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"
)

func init() {
	hclog.L().SetLevel(hclog.Trace)
}

func TestFunc(t *testing.T) {
	cases := []struct {
		Name     string
		Callback interface{}
		Args     []Arg
		Out      []interface{}
		Err      string
	}{
		{
			"basic matching",
			func(in struct {
				A, B int
			}) int {
				return in.A + in.B
			},
			[]Arg{
				Named("a", 12),
				Named("b", 24),
			},
			[]interface{}{
				36,
			},
			"",
		},

		{
			"missing argument",
			func(in struct {
				A, B int
			}) int {
				return in.A + in.B
			},
			[]Arg{
				Named("a", 12),
			},
			nil,
			"argument cannot",
		},

		{
			"unexported field ignored",
			func(in struct {
				A int
				b int
			}) int {
				return in.A
			},
			[]Arg{Named("a", 12)},
			[]interface{}{12},
			"",
		},

		{
			"renamed field",
			func(in struct {
				A int `argmapper:"C"`
				B int
			}) int {
				return in.A + in.B
			},
			[]Arg{
				Named("b", 24),
				Named("c", 12),
			},
			[]interface{}{
				36,
			},
			"",
		},

		{
			"direct named converter",
			func(in struct {
				A string
			}) string {
				return in.A + "!"
			},
			[]Arg{
				Named("a", 12),
				WithConvFunc(func(s struct {
					A int
				}) struct {
					A string
				} {
					return struct{ A string }{strconv.Itoa(s.A)}
				}),
			},
			[]interface{}{"12!"},
			"",
		},

		{
			"generic type converter",
			func(in struct {
				A string
				B int
			}) string {
				return strings.Repeat(in.A, in.B)
			},
			[]Arg{
				Named("a", 12),
				Named("b", 2),
				WithConvFunc(func(s struct {
					C string
				}) struct {
					A string
				} {
					return struct {
						A string
					}{"FOO"}
				}),
				WithConvFunc(func(s struct {
					C bool
				}) struct {
					A string
				} {
					return struct {
						A string
					}{"FOO"}
				}),
				WithConvFunc(func(s struct {
					B int `argmapper:",wildcard"`
				}) struct {
					B string `argmapper:",wildcard"`
				} {
					return struct {
						B string `argmapper:",wildcard"`
					}{strconv.Itoa(s.B)}
				}),
			},
			[]interface{}{"1212"},
			"",
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)

			f, err := NewFunc(tt.Callback)
			require.NoError(err)
			result := f.Call(tt.Args...)

			// If we expect an error, check that
			if tt.Err == "" {
				require.NoError(result.Err())
			} else {
				require.Error(result.Err())
				require.Contains(result.Err().Error(), tt.Err)
			}

			// Verify outputs
			require.Equal(len(tt.Out), result.Len())
			for i, out := range tt.Out {
				require.Equal(out, result.Out(i))
			}
		})
	}
}
