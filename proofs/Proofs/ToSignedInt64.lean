/-
  ToSignedInt64.lean — Proofs about the toSignedInt64 function.

  Key properties:
  1. Values ≤ MaxInt32 are preserved as-is
  2. Values in (MaxInt32, MaxUint32] are wrapped as signed int32
  3. Values > MaxUint32 are preserved as uint64→int64
-/
import Proofs.Spec

/-! ## Concrete verification -/

/-- Zero maps to zero. -/
theorem toSignedInt64_zero : toSignedInt64 0 = 0 := by native_decide

/-- Small positive value preserved. -/
theorem toSignedInt64_positive : toSignedInt64 42 = 42 := by native_decide

/-- MaxInt32 (2147483647) is preserved. -/
theorem toSignedInt64_max_int32 : toSignedInt64 2147483647 = 2147483647 := by native_decide

/-- MaxInt32 + 1 (2147483648 = 0x80000000) wraps to -2147483648. -/
theorem toSignedInt64_min_int32_unsigned : toSignedInt64 2147483648 = -2147483648 := by
  native_decide

/-- MaxUint32 (4294967295) wraps to -1 (c2ffi represents enum -1 as 4294967295). -/
theorem toSignedInt64_minus_one : toSignedInt64 4294967295 = -1 := by native_decide

/-- MaxUint32 + 1 (4294967296) is NOT in the uint32 wrap range, stays as-is. -/
theorem toSignedInt64_above_uint32 : toSignedInt64 4294967296 = 4294967296 := by native_decide
