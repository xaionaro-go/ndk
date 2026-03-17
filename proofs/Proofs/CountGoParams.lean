/-
  CountGoParams.lean — Proofs about the countGoParams function.

  Key properties:
  1. func() → 0
  2. func(T) → 1
  3. func(T, U) → 2
  4. No parens → 0
-/
import Proofs.Spec

/-! ## Concrete verification -/

/-- Empty signature yields 0. -/
theorem countGoParams_empty : countGoParams "" = 0 := by
  native_decide

/-- func() yields 0 parameters. -/
theorem countGoParams_no_params : countGoParams "func()" = 0 := by
  native_decide

/-- func(int) yields 1 parameter. -/
theorem countGoParams_one_param : countGoParams "func(int)" = 1 := by
  native_decide

/-- func(int, string) yields 2 parameters. -/
theorem countGoParams_two_params : countGoParams "func(int, string)" = 2 := by
  native_decide

/-- func(int, string, bool) yields 3 parameters. -/
theorem countGoParams_three_params : countGoParams "func(int, string, bool)" = 3 := by
  native_decide

/-- No parentheses yields 0. -/
theorem countGoParams_no_parens : countGoParams "hello" = 0 := by
  native_decide

/-- func( ) with only whitespace yields 0. -/
theorem countGoParams_whitespace : countGoParams "func( )" = 0 := by
  native_decide
