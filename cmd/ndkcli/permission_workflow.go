package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xaionaro-go/ndk/permission"
)

var permissionCmd = &cobra.Command{
	Use:   "permission",
	Short: "permission NDK module",
}

var permissionCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check whether a permission is granted for a given PID/UID",
	RunE: func(cmd *cobra.Command, args []string) (_err error) {
		name, _ := cmd.Flags().GetString("name")
		pid, _ := cmd.Flags().GetInt32("pid")
		uid, _ := cmd.Flags().GetInt32("uid")

		var outResult int32
		ret := permission.CheckPermission(
			name,
			permission.Pid_t(pid),
			permission.Uid_t(uid),
			&outResult,
		)

		fmt.Printf("permission: %s\n", name)
		fmt.Printf("pid:        %d\n", pid)
		fmt.Printf("uid:        %d\n", uid)
		fmt.Printf("return:     %d\n", ret)

		switch permission.Result(outResult) {
		case permission.Granted:
			fmt.Println("result:     granted")
		case permission.Denied:
			fmt.Println("result:     denied")
		default:
			fmt.Printf("result:     unknown (%d)\n", outResult)
		}

		return nil
	},
}

func init() {
	permissionCheckCmd.Flags().String("name", "", "permission name (e.g. android.permission.CAMERA)")
	_ = permissionCheckCmd.MarkFlagRequired("name")
	permissionCheckCmd.Flags().Int32("pid", 0, "process ID to check")
	permissionCheckCmd.Flags().Int32("uid", 0, "user ID to check")

	permissionCmd.AddCommand(permissionCheckCmd)
	rootCmd.AddCommand(permissionCmd)
}
