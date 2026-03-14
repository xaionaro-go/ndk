package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xaionaro-go/ndk/permission"
)

var permissionCmd = &cobra.Command{
	Use:   "permission",
	Short: "Permission NDK operations",
}

var (
	permCheckName string
	permCheckPID  int32
	permCheckUID  uint32
)

var permissionCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check an Android permission for a given pid/uid",
	RunE: func(cmd *cobra.Command, args []string) error {
		var outResult int32
		rc := permission.CheckPermission(
			permCheckName,
			permission.Pid_t(permCheckPID),
			permission.Uid_t(permCheckUID),
			&outResult,
		)
		fmt.Printf("Permission: %s\n", permCheckName)
		fmt.Printf("PID:        %d\n", permCheckPID)
		fmt.Printf("UID:        %d\n", permCheckUID)
		fmt.Printf("Return:     %d\n", rc)
		switch permission.Result(outResult) {
		case permission.Granted:
			fmt.Println("Result:     GRANTED")
		case permission.Denied:
			fmt.Println("Result:     DENIED")
		default:
			fmt.Printf("Result:     UNKNOWN (%d)\n", outResult)
		}
		return nil
	},
}

func init() {
	permissionCheckCmd.Flags().StringVar(&permCheckName, "name", "", "permission name (e.g. android.permission.CAMERA)")
	permissionCheckCmd.Flags().Int32Var(&permCheckPID, "pid", 0, "process ID")
	permissionCheckCmd.Flags().Uint32Var(&permCheckUID, "uid", 0, "user ID")
	_ = permissionCheckCmd.MarkFlagRequired("name")
	_ = permissionCheckCmd.MarkFlagRequired("pid")
	_ = permissionCheckCmd.MarkFlagRequired("uid")

	permissionCmd.AddCommand(permissionCheckCmd)
	rootCmd.AddCommand(permissionCmd)
}
