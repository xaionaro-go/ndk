package c2ffi

import (
	"os"
	"testing"
)

func TestEnumIDCollision(t *testing.T) {
	data, err := os.ReadFile("/tmp/logging_pp_c2ffi.json")
	if err != nil {
		t.Skip("no test data")
	}
	spec, err := Convert(data, ConvertOptions{
		Module:        "logging",
		SourcePackage: "github.com/AndroidGoLab/ndk/capi/logging",
	})
	if err != nil {
		t.Fatal(err)
	}

	for k, vals := range spec.Enums {
		t.Logf("Enum %s:", k)
		for _, v := range vals {
			t.Logf("  %s = %d", v.Name, v.Value)
		}
	}
}
