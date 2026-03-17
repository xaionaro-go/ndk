/-
  ConstructorDestructor.lean — Proofs about constructor/destructor classification.

  Key properties:
  1. Disjointness: a function cannot be both constructor and destructor for the same type
  2. Detection correctness: matching suffixes are detected
  3. Non-matching names are rejected

  The disjointness property is critical for the merge logic: functions are assigned
  to exactly one role.
-/
import Proofs.Spec

set_option maxRecDepth 1024

/-! ## String append injectivity -/

private theorem string_append_left_cancel {a b c : String} (h : a ++ b = a ++ c) : b = c := by
  have hd : a.data ++ b.data = a.data ++ c.data := by
    rw [← String.data_append, ← String.data_append, h]
  have := List.append_cancel_left hd
  exact congrArg String.mk this

/-! ## Disjointness: constructor suffixes are never destructor suffixes -/

private theorem ne_of_append_ne {t a b : String} (hab : a ≠ b) : t ++ a ≠ t ++ b :=
  fun h => absurd (string_append_left_cancel h) hab

theorem create_not_destructor (t : String) :
    isDestructorFunc (t ++ "_create") t = false := by
  simp only [isDestructorFunc, destructorSuffixes, List.any_cons, List.any_nil, Bool.or_false,
    Bool.or_eq_false_iff, beq_eq_false_iff_ne, ne_eq]
  exact ⟨ne_of_append_ne (by decide), ne_of_append_ne (by decide),
         ne_of_append_ne (by decide), ne_of_append_ne (by decide),
         ne_of_append_ne (by decide)⟩

theorem new_not_destructor (t : String) :
    isDestructorFunc (t ++ "_new") t = false := by
  simp only [isDestructorFunc, destructorSuffixes, List.any_cons, List.any_nil, Bool.or_false,
    Bool.or_eq_false_iff, beq_eq_false_iff_ne, ne_eq]
  exact ⟨ne_of_append_ne (by decide), ne_of_append_ne (by decide),
         ne_of_append_ne (by decide), ne_of_append_ne (by decide),
         ne_of_append_ne (by decide)⟩

theorem delete_not_constructor (t : String) :
    isConstructorFunc (t ++ "_delete") t = false := by
  simp only [isConstructorFunc, constructorSuffixes, List.any_cons, List.any_nil, Bool.or_false,
    Bool.or_eq_false_iff, beq_eq_false_iff_ne, ne_eq]
  exact ⟨ne_of_append_ne (by decide), ne_of_append_ne (by decide)⟩

theorem free_not_constructor (t : String) :
    isConstructorFunc (t ++ "_free") t = false := by
  simp only [isConstructorFunc, constructorSuffixes, List.any_cons, List.any_nil, Bool.or_false,
    Bool.or_eq_false_iff, beq_eq_false_iff_ne, ne_eq]
  exact ⟨ne_of_append_ne (by decide), ne_of_append_ne (by decide)⟩

/-! ## Detection correctness -/

theorem constructor_create_detected (t : String) :
    isConstructorFunc (t ++ "_create") t = true := by
  simp [isConstructorFunc, constructorSuffixes]

theorem constructor_new_detected (t : String) :
    isConstructorFunc (t ++ "_new") t = true := by
  simp [isConstructorFunc, constructorSuffixes]

theorem destructor_delete_detected (t : String) :
    isDestructorFunc (t ++ "_delete") t = true := by
  simp [isDestructorFunc, destructorSuffixes]

theorem destructor_release_detected (t : String) :
    isDestructorFunc (t ++ "_release") t = true := by
  simp [isDestructorFunc, destructorSuffixes]

/-! ## Concrete test vectors -/

theorem ctor_ASensor_create :
    isConstructorFunc "ASensor_create" "ASensor" = true := by native_decide
theorem dtor_ASensor_delete :
    isDestructorFunc "ASensor_delete" "ASensor" = true := by native_decide
theorem not_ctor_ASensor_delete :
    isConstructorFunc "ASensor_delete" "ASensor" = false := by native_decide
theorem not_dtor_ASensor_create :
    isDestructorFunc "ASensor_create" "ASensor" = false := by native_decide
theorem not_ctor_unrelated :
    isConstructorFunc "ASensor_getName" "ASensor" = false := by native_decide
theorem not_dtor_unrelated :
    isDestructorFunc "ASensor_getName" "ASensor" = false := by native_decide
