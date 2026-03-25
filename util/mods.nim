import std/[bitops, strutils]

type Mask* = distinct byte

type Info* = tuple
  name: string
  keys: seq[string]

const MODS = [
  (@["SHIFT"], @["shift_l", "shift_r"]),
  (@["CAPS"], @["caps_lock"]),
  (@["CTRL", "CONTROL"], @["control_l", "control_r"]),
  (@["ALT"], @["alt_l", "alt_r"]),
  (@["MOD2"], @[]),
  (@["MOD3"], @[]),
  (@["SUPER", "WIN", "LOGO", "MOD4"], @["super_l", "super_r"]),
  (@["MOD5"], @[]),
]

const Empty* = 0.Mask

proc toMask*(mods: string): Mask =
  for i, (names, _) in MODS:
    for name in names:
      if mods.contains name:
        result.byte.setBit i

proc `$`*(mask: Mask): string =
  var mods: seq[string]
  for i, (names, _) in MODS:
    if mask.byte.testBit i:
      mods.add names[0]
  mods.join "_"

iterator info*(mask: Mask): Info =
  for i, (names, keys) in MODS:
    if mask.byte.testBit i:
      yield (names[0], keys)

proc includes*(parent, child: Mask): bool =
  child.byte == parent.byte.bitand child.byte

proc without*(parent, child: Mask): Mask =
  Mask(parent.byte.clearMasked child.byte)

proc `==`*(a, b: Mask): bool {.borrow.}
proc `<`*(a, b: Mask): bool {.borrow.}
