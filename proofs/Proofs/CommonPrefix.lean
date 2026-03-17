/-
  CommonPrefix.lean — Proofs about the commonPrefix function.

  Key properties verified:
  1. Symmetry: commonPrefix a b = commonPrefix b a  (structural proof)
  2. Self-identity: commonPrefix a a = a  (structural proof)
  3. Empty absorption: commonPrefix "" b = ""  (structural proof)
  4. Concrete test vectors  (native_decide)
-/
import Proofs.Spec

/-! ## Helper lemma: go is symmetric -/

private theorem go_symmetric (as bs : List Char) :
    commonPrefix.go as bs = commonPrefix.go bs as := by
  induction as generalizing bs with
  | nil =>
    cases bs with
    | nil => rfl
    | cons _ _ => simp [commonPrefix.go]
  | cons a as' ih =>
    cases bs with
    | nil => simp [commonPrefix.go]
    | cons b bs' =>
      simp only [commonPrefix.go]
      by_cases hab : (a == b) = true
      · have hba : (b == a) = true := by
          simp [BEq.beq, beq_iff_eq] at hab ⊢; exact hab.symm
        simp [hab, hba]
        have : a = b := by simpa [BEq.beq, beq_iff_eq] using hab
        subst this; exact ⟨rfl, ih bs'⟩
      · have hba : ¬((b == a) = true) := by
          simp [BEq.beq, beq_iff_eq] at hab ⊢; exact fun h => hab h.symm
        simp [hab, hba]

private theorem go_self (cs : List Char) :
    commonPrefix.go cs cs = cs := by
  induction cs with
  | nil => simp [commonPrefix.go]
  | cons c rest ih => simp [commonPrefix.go, ih]

/-! ## Main theorems -/

theorem commonPrefix_symmetric (a b : String) :
    commonPrefix a b = commonPrefix b a := by
  simp [commonPrefix]
  exact congrArg String.mk (go_symmetric a.data b.data)

theorem commonPrefix_self (s : String) :
    commonPrefix s s = s := by
  simp [commonPrefix]
  exact congrArg String.mk (go_self s.data)

theorem commonPrefix_empty_left (s : String) :
    commonPrefix "" s = "" := by
  simp [commonPrefix, commonPrefix.go]

theorem commonPrefix_empty_right (s : String) :
    commonPrefix s "" = "" := by
  rw [commonPrefix_symmetric]; exact commonPrefix_empty_left s

/-! ## Concrete test vectors -/

theorem commonPrefix_abc_abd : commonPrefix "abc" "abd" = "ab" := by native_decide
theorem commonPrefix_hello_help : commonPrefix "hello" "help" = "hel" := by native_decide
theorem commonPrefix_disjoint : commonPrefix "abc" "xyz" = "" := by native_decide
theorem commonPrefix_same : commonPrefix "test" "test" = "test" := by native_decide
theorem commonPrefix_prefix_of : commonPrefix "ab" "abcde" = "ab" := by native_decide
