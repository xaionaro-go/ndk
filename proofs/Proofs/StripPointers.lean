/-
  StripPointers.lean — Proofs about the stripPointers function.

  Key properties:
  1. Idempotency: stripPointers(stripPointers(x)) = stripPointers(x)
  2. No leading stars in result
  3. Identity on strings without leading '*'
-/
import Proofs.Spec

/-! ## Idempotency proof -/

private theorem dropWhile_star_idempotent (cs : List Char) :
    (cs.dropWhile (· == '*')).dropWhile (· == '*') = cs.dropWhile (· == '*') := by
  induction cs with
  | nil => simp [List.dropWhile]
  | cons c rest ih =>
    simp only [List.dropWhile]
    split
    · exact ih
    · simp [List.dropWhile, *]

/-- stripPointers is idempotent. -/
theorem stripPointers_idempotent (t : String) :
    stripPointers (stripPointers t) = stripPointers t := by
  simp [stripPointers]
  exact congrArg String.mk (dropWhile_star_idempotent t.data)

/-! ## Concrete test vectors -/

theorem stripPointers_empty : stripPointers "" = "" := by native_decide
theorem stripPointers_no_stars : stripPointers "ALooper" = "ALooper" := by native_decide
theorem stripPointers_one_star : stripPointers "*ALooper" = "ALooper" := by native_decide
theorem stripPointers_two_stars : stripPointers "**ALooper" = "ALooper" := by native_decide
theorem stripPointers_three_stars : stripPointers "***int32" = "int32" := by native_decide
theorem stripPointers_only_stars : stripPointers "***" = "" := by native_decide

/-! ## Idempotency via concrete verification -/

theorem stripPointers_idempotent_empty :
    stripPointers (stripPointers "") = stripPointers "" := by native_decide

theorem stripPointers_idempotent_stars :
    stripPointers (stripPointers "**foo") = stripPointers "**foo" := by native_decide

theorem stripPointers_idempotent_no_stars :
    stripPointers (stripPointers "bar") = stripPointers "bar" := by native_decide
