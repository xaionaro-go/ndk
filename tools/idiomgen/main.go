// Command idiomgen generates idiomatic Go packages from specs and overlays.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/xaionaro-go/ndk/tools/pkg/idiomgen"
)

func main() {
	var (
		specPath    = flag.String("spec", "", "path to generated spec YAML")
		overlayPath = flag.String("overlay", "", "path to overlay YAML (optional)")
		tmplDir     = flag.String("templates", "templates", "path to template directory")
		outDir      = flag.String("out", "", "output directory for generated Go package")
		capiDir     = flag.String("capi-dir", "", "output directory for generated bridge files in capi package (optional)")
	)
	flag.Parse()

	if *specPath == "" || *outDir == "" {
		fmt.Fprintln(os.Stderr, "usage: idiomgen -spec <file.yaml> -out <dir> [-overlay <file.yaml>] [-templates <dir>] [-capi-dir <dir>]")
		os.Exit(1)
	}

	spec, err := idiomgen.LoadSpec(*specPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load spec: %v\n", err)
		os.Exit(1)
	}

	overlay, err := idiomgen.LoadOverlayOrEmpty(*overlayPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load overlay: %v\n", err)
		os.Exit(1)
	}

	if err := idiomgen.Generate(spec, overlay, *tmplDir, *outDir, *capiDir); err != nil {
		fmt.Fprintf(os.Stderr, "generate: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("idiomgen: generated %s -> %s\n", spec.Module, *outDir)
}
