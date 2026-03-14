package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xaionaro-go/ndk/nnapi"
)

var nnapiCmd = &cobra.Command{
	Use:   "nnapi",
	Short: "Neural Networks API operations",
}

var nnapiInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show NNAPI capability information",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("NNAPI surface:")
		fmt.Println("  - Model:       create/configure neural network models (NewModel, AddOperation, Finish)")
		fmt.Println("  - Compilation: compile a model for execution (Finish, SetPreference, SetPriority, SetTimeout)")
		fmt.Println("  - Execution:   run inference (Compute, StartCompute, SetMeasureTiming, GetDuration)")
		fmt.Println("  - Burst:       reusable execution object for repeated inferences")
		fmt.Println("  - Event:       async completion handle (Wait)")
		fmt.Println()

		model, err := nnapi.NewModel()
		if err != nil {
			fmt.Printf("NNAPI availability: NOT available (%v)\n", err)
			return nil
		}
		defer model.Close()
		fmt.Println("NNAPI availability: OK (model created successfully)")
		return nil
	},
}

func init() {
	nnapiCmd.AddCommand(nnapiInfoCmd)
	rootCmd.AddCommand(nnapiCmd)
}
