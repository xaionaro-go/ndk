// Simulates c-for-go output for Android Permission Manager (android/permission_manager.h).
// This file is parsed at AST level only; it does not compile.
package permission

// Integer typedefs.
type Permission_result_t int32

// Result codes.
const (
	PERMISSION_MANAGER_PERMISSION_GRANTED Permission_result_t = 0
	PERMISSION_MANAGER_PERMISSION_DENIED  Permission_result_t = -1
)

// --- Permission functions ---
func APermissionManager_checkPermission(permission *byte, pid int32, uid uint32, outResult *int32) int32 { return 0 }
