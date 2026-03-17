import Proofs

/-!
  Differential test oracle.

  Reads commands from stdin, runs the Lean spec functions, prints results.
  Protocol: one command per line, tab-separated fields.

    exportName\t<name>
    stripPointers\t<type>
    commonPrefix\t<a>\t<b>
    isConstructorFunc\t<funcName>\t<specTypeName>
    isDestructorFunc\t<funcName>\t<specTypeName>
    countGoParams\t<sig>
    isScalarGoType\t<type>
    resolveType\t<specType>\t<key1>=<val1>,<key2>=<val2>,...
    toSignedInt64\t<value>

  Output: one result per line.
-/

/-- Parse a comma-separated list of key=value pairs into an assoc list. -/
def parseTypeMap (s : String) : List (String × String) :=
  if s.isEmpty then []
  else
    let pairs := s.splitOn ","
    pairs.filterMap fun pair =>
      match pair.splitOn "=" with
      | [k, v] => some (k, v)
      | _ => none

def processLine (line : String) : IO Unit := do
  let fields := line.splitOn "\t"
  match fields with
  | ["exportName", name] =>
    IO.println (exportName name)
  | ["stripPointers", t] =>
    IO.println (stripPointers t)
  | ["commonPrefix", a, b] =>
    IO.println (commonPrefix a b)
  | ["isConstructorFunc", funcName, specTypeName] =>
    IO.println (toString (isConstructorFunc funcName specTypeName))
  | ["isDestructorFunc", funcName, specTypeName] =>
    IO.println (toString (isDestructorFunc funcName specTypeName))
  | ["countGoParams", sig] =>
    IO.println (toString (countGoParams sig))
  | ["isScalarGoType", t] =>
    IO.println (toString (isScalarGoType t))
  | ["resolveType", specType, mapStr] =>
    let typeMap := parseTypeMap mapStr
    IO.println (resolveType specType typeMap)
  | ["toSignedInt64", valStr] =>
    match valStr.toNat? with
    | some n => IO.println (toString (toSignedInt64 (UInt64.ofNat n)))
    | none => IO.println "ERROR: invalid uint64"
  | _ =>
    IO.eprintln s!"Unknown command: {line}"

def main : IO Unit := do
  let stdin ← IO.getStdin
  let mut done := false
  while !done do
    let line ← stdin.getLine
    if line.isEmpty then
      done := true
    else
      -- Only strip trailing newline/carriage-return, not tabs (which delimit fields).
      let trimmed := line.dropRightWhile (fun c => c == '\n' || c == '\r')
      if !trimmed.isEmpty then
        processLine trimmed
