package pkg

import (
	"testing"

	"github.com/hyahm/golog"
)

type compareStruct struct {
	one     string
	two     string
	sep     string
	expect  bool
	compare string // "gt" | "ge" |"lt" |"le"
}

func TestCompareVersion(t *testing.T) {
	defer golog.Sync()
	// cs := []compareStruct{
	// 	{
	// 		one:     "v1.5.12",
	// 		two:     "v1.6.1",
	// 		sep:     ".",
	// 		compare: "le",
	// 		expect:  true,
	// 	},
	// 	{
	// 		one:     "v1.6.12",
	// 		two:     "v1.6.5",
	// 		sep:     ".",
	// 		compare: "gt",
	// 		expect:  true,
	// 	},
	// 	{
	// 		one:     "v1.60.1",
	// 		two:     "v1.6.5",
	// 		sep:     ".",
	// 		compare: "gt",
	// 		expect:  true,
	// 	},
	// 	{
	// 		one:     "v5.0.1",
	// 		two:     "v1.6.5",
	// 		sep:     ".",
	// 		compare: "gt",
	// 		expect:  true,
	// 	},
	// 	{
	// 		one:     "v1.6.5.3",
	// 		two:     "v1.6.5",
	// 		sep:     ".",
	// 		compare: "le",
	// 		expect:  false,
	// 	},
	// 	{
	// 		one:     "v1.6.5.3",
	// 		two:     "v1.6.5",
	// 		sep:     ".",
	// 		compare: "gt",
	// 		expect:  true,
	// 	},
	// 	{
	// 		one:     "v1.6.5.3",
	// 		two:     "v1.6.5",
	// 		sep:     ".",
	// 		compare: "ge",
	// 		expect:  true,
	// 	},
	// 	{
	// 		one:     "v1.6.5.3",
	// 		two:     "v1.6.5",
	// 		sep:     ".",
	// 		compare: "lt",
	// 		expect:  false,
	// 	},
	// 	{
	// 		one:     "v1.6.4.3",
	// 		two:     "v1.6.5",
	// 		sep:     ".",
	// 		compare: "lt",
	// 		expect:  true,
	// 	},
	// }
	// for _, v := range cs {
	// 	switch v.compare {
	// 	case "gt":
	// 		if gt(v.one, v.two, v.sep) != v.expect {
	// 			t.Errorf("%#v\n", v)
	// 		}
	// 	case "le":

	// 		if Le(v.one, v.two, v.sep) != v.expect {
	// 			t.Log(le(v.one, v.two, v.sep))
	// 			t.Errorf("%#v\n", v)
	// 		}
	// 	}
	// }
}
