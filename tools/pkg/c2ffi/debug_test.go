package c2ffi

import (
	"fmt"
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
		SourcePackage: "github.com/xaionaro-go/ndk/capi/logging",
	})
	if err != nil {
		t.Fatal(err)
	}

	for k, vals := range spec.Enums {
		fmt.Printf("Enum %s:\n", k)
		for _, v := range vals {
			fmt.Printf("  %s = %d\n", v.Name, v.Value)
		}
	}
}
