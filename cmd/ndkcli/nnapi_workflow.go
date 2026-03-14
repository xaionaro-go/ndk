package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xaionaro-go/ndk/nnapi"
)

var nnapiProbeCmd = &cobra.Command{
	Use:   "probe",
	Short: "Probe NNAPI availability by creating and finishing an empty model",
	RunE: func(cmd *cobra.Command, args []string) (_err error) {
		model, err := nnapi.NewModel()
		if err != nil {
			fmt.Printf("NNAPI not available: %v\n", err)
			return nil
		}
		defer model.Close()

		fmt.Println("NNAPI model created successfully")

		if err := model.Finish(); err != nil {
			fmt.Printf("model.Finish (empty model): %v\n", err)
			return nil
		}
		fmt.Println("model.Finish completed successfully")

		return nil
	},
}

func init() {
	nnapiCmd.AddCommand(nnapiProbeCmd)
}
