// Command capigen generates raw CGo binding packages from spec YAML and
// manifest configuration.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/AndroidGoLab/ndk/tools/pkg/capigen"
	"github.com/AndroidGoLab/ndk/tools/pkg/specmodel"
	"gopkg.in/yaml.v3"
)

func main() {
	var (
		specPath     = flag.String("spec", "", "path to spec YAML file (required)")
		manifestPath = flag.String("manifest", "", "path to capi manifest YAML (required)")
		overlayPath  = flag.String("overlay", "", "path to overlay YAML (optional, for api_levels)")
		outDir       = flag.String("out", "", "output directory for generated package (required)")
	)
	flag.Parse()

	if *specPath == "" || *manifestPath == "" || *outDir == "" {
		fmt.Fprintln(os.Stderr, "usage: capigen -spec <spec.yaml> -manifest <manifest.yaml> -out <dir> [-overlay <overlay.yaml>]")
		os.Exit(1)
	}

	spec, err := readSpec(*specPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read spec: %v\n", err)
		os.Exit(1)
	}

	manifest, err := readManifest(*manifestPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read manifest: %v\n", err)
		os.Exit(1)
	}

	var apiLevels map[string]int
	if *overlayPath != "" {
		apiLevels, err = readAPILevels(*overlayPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read overlay: %v\n", err)
			os.Exit(1)
		}
	}

	if err := capigen.GeneratePackage(spec, manifest, *outDir, apiLevels); err != nil {
		fmt.Fprintf(os.Stderr, "generate: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("capigen: wrote %s/ (%d types, %d enums, %d functions, %d callbacks, %d structs)\n",
		*outDir, len(spec.Types), len(spec.Enums), len(spec.Functions), len(spec.Callbacks), len(spec.Structs))
}

func readSpec(path string) (*specmodel.Spec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var spec specmodel.Spec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, err
	}

	return &spec, nil
}

func readManifest(path string) (*capigen.Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var m capigen.Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

// readAPILevels reads only the api_levels map from an overlay YAML file.
func readAPILevels(path string) (map[string]int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var overlay struct {
		APILevels map[string]int `yaml:"api_levels"`
	}
	if err := yaml.Unmarshal(data, &overlay); err != nil {
		return nil, err
	}

	return overlay.APILevels, nil
}
