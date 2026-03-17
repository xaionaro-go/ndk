/-
  Spec.lean — Formal specifications of NDK code generator algorithms.

  Models the core pure functions from the Go codebase as Lean 4 definitions.
  These are the reference implementations that the proofs verify.
-/

/-! ## String utilities -/

/-- Strip all leading underscores from a string. -/
def stripLeadingUnderscores (s : String) : String :=
  ⟨s.data.dropWhile (· == '_')⟩

/-- Capitalize the first character if it is ASCII lowercase. -/
def capitalizeFirst (s : String) : String :=
  match s.data with
  | [] => s
  | c :: rest =>
    if c.val ≥ 'a'.val ∧ c.val ≤ 'z'.val then
      ⟨⟨c.val - 'a'.val + 'A'.val, sorry⟩ :: rest⟩
    else s

/-- ExportName: strip leading underscores and capitalize.
    Models `capigen.ExportName` from generate.go:23-37. -/
def exportName (name : String) : String :=
  let stripped := stripLeadingUnderscores name
  if stripped.isEmpty then name
  else capitalizeFirst stripped

/-- Strip all leading '*' from a type string.
    Models `capigen.stripPointers` from generate.go:1331-1336. -/
def stripPointers (t : String) : String :=
  ⟨t.data.dropWhile (· == '*')⟩

/-- Check if a string is a Go scalar type.
    Models `capigen.isScalarGoType` from typeconv.go:155-162. -/
def isScalarGoType (t : String) : Bool :=
  t ∈ ["int8", "uint8", "int16", "uint16", "int32", "uint32",
        "int64", "uint64", "float32", "float64", "bool", "int", "uint"]

/-- Check if a string is a fixed-size array type like "[16]float32".
    Models `capigen.isFixedArrayType` from typeconv.go:165-167. -/
def isFixedArrayType (t : String) : Bool :=
  t.length > 2 ∧ t.get 0 == '[' ∧ t.get ⟨1⟩ != ']'

/-! ## Common prefix -/

/-- Compute the longest common prefix of two strings.
    Models `c2ffi.commonPrefix` from tospec.go:579-590. -/
def commonPrefix (a b : String) : String :=
  ⟨go a.data b.data⟩
where
  go : List Char → List Char → List Char
    | c₁ :: r₁, c₂ :: r₂ => if c₁ == c₂ then c₁ :: go r₁ r₂ else []
    | _, _ => []

/-! ## Constructor / destructor classification -/

/-- Destructor function suffixes. -/
def destructorSuffixes : List String :=
  ["_delete", "_free", "_destroy", "_release", "_close"]

/-- Constructor function suffixes. -/
def constructorSuffixes : List String :=
  ["_create", "_new"]

/-- Check if funcName is a destructor for specTypeName.
    Models `idiomgen.isDestructorFunc` from merge.go:95-102. -/
def isDestructorFunc (funcName specTypeName : String) : Bool :=
  destructorSuffixes.any (fun suffix => funcName == specTypeName ++ suffix)

/-- Check if funcName is a constructor for specTypeName.
    Models `idiomgen.isConstructorFunc` from merge.go:126-133. -/
def isConstructorFunc (funcName specTypeName : String) : Bool :=
  constructorSuffixes.any (fun suffix => funcName == specTypeName ++ suffix)

/-! ## toSignedInt64 -/

/-- The maximum value of a signed 32-bit integer. -/
def maxInt32 : Nat := 2147483647

/-- The maximum value of an unsigned 32-bit integer. -/
def maxUint32 : Nat := 4294967295

/-- Convert an unsigned uint64 enum value to signed int64.
    Models `c2ffi.toSignedInt64` from tospec.go:470-475.
    c2ffi outputs all enum values as unsigned, so -1 becomes 4294967295. -/
def toSignedInt64 (v : UInt64) : Int64 :=
  if v.toNat > maxInt32 ∧ v.toNat ≤ maxUint32 then
    -- Reinterpret as signed 32-bit, then extend to 64-bit.
    Int64.ofInt (Int32.ofNat v.toNat).toInt
  else
    Int64.ofNat v.toNat

/-! ## countGoParams -/

/-- Count commas in a string. -/
def countCommas (s : String) : Nat :=
  s.data.filter (· == ',') |>.length

/-- Count parameters in a Go function signature string.
    Models `idiomgen.countGoParams` from merge.go:1502-1514. -/
def countGoParams (sig : String) : Nat :=
  match sig.data.dropWhile (· != '(') with
  | [] => 0
  | _ :: rest =>
    let inner := rest.takeWhile (· != ')')
    let trimmed := inner.dropWhile Char.isWhitespace
    if trimmed.isEmpty then 0
    else countCommas ⟨inner⟩ + 1

/-! ## resolveType -/

/-- A simple model of the typeMap-based type resolution.
    Models `idiomgen.resolveType` from merge.go:1468-1482. -/
def resolveType (specType : String) (typeMap : List (String × String)) : String :=
  -- Direct lookup.
  match typeMap.lookup specType with
  | some goType => goType
  | none =>
    -- Try stripping prefixes.
    let prefixes := ["[]", "**", "*"]
    go prefixes
where
  go : List String → String
    | [] => specType
    | pfx :: rest =>
      if specType.startsWith pfx then
        let base := specType.drop pfx.length
        match typeMap.lookup base with
        | some goType => pfx ++ goType
        | none => go rest
      else go rest

/-! ## stripAndTitle -/

/-- Capitalize the first character of a string (for toTitleCase). -/
private def capFirst (s : String) : String := capitalizeFirst s

/-- Convert UPPER_SNAKE_CASE to TitleCase. -/
def toTitleCase (s : String) : String :=
  let lower := s.toLower
  let parts := lower.splitOn "_"
  let titled := parts.filterMap fun p =>
    if p.isEmpty then none
    else some (capFirst p)
  String.join titled

/-- Strip a given constant string from the beginning of name and TitleCase the remainder.
    Returns the original name if the result would be empty.
    Models `idiomgen.stripAndTitle` from merge.go:1107-1116. -/
def stripAndTitle (name pfx : String) : String :=
  if pfx.isEmpty then name
  else
    let s := if name.startsWith pfx then name.drop pfx.length else name
    if s.isEmpty then name
    else toTitleCase s
