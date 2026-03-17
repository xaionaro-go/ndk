/-
  ResolveType.lean — Proofs about the resolveType function.

  Key properties:
  1. Identity for empty typeMap
  2. Direct lookup takes precedence
  3. Prefix preservation
-/
import Proofs.Spec

/-! ## Main theorems -/

/-- resolveType returns the input unchanged when the typeMap is empty. -/
theorem resolveType_empty_map (specType : String) :
    resolveType specType [] = specType := by
  simp [resolveType, resolveType.go]

/-- Direct lookup takes precedence over prefix-based lookup. -/
theorem resolveType_direct_precedence (specType goType : String)
    (typeMap : List (String × String))
    (h : typeMap.lookup specType = some goType) :
    resolveType specType typeMap = goType := by
  simp [resolveType, h]

/-! ## Concrete test vectors -/

-- Direct lookup.
theorem resolveType_direct : resolveType "foo" [("foo", "bar")] = "bar" := by native_decide
-- Pointer prefix preservation.
theorem resolveType_star : resolveType "*foo" [("foo", "Bar")] = "*Bar" := by native_decide
-- Double pointer prefix preservation.
theorem resolveType_dstar : resolveType "**foo" [("foo", "Bar")] = "**Bar" := by native_decide
-- Slice prefix preservation.
theorem resolveType_slice : resolveType "[]foo" [("foo", "Bar")] = "[]Bar" := by native_decide
-- Unknown type returns unchanged.
theorem resolveType_unknown : resolveType "unknown" [("foo", "Bar")] = "unknown" := by
  native_decide
-- Direct lookup takes precedence over prefix.
theorem resolveType_direct_over_prefix :
    resolveType "*foo" [("*foo", "Direct"), ("foo", "Prefix")] = "Direct" := by native_decide

/-! ## Fixed-point property for concrete cases -/

-- Resolving an already-resolved type is identity (when result is not in map).
theorem resolveType_fixedpoint_simple :
    resolveType (resolveType "foo" [("foo", "Bar")]) [("foo", "Bar")] = "Bar" := by native_decide

-- Resolving an unknown type twice is identity.
theorem resolveType_fixedpoint_unknown :
    resolveType (resolveType "unknown" [("foo", "Bar")]) [("foo", "Bar")] = "unknown" := by
  native_decide
