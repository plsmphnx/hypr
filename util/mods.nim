import std/[bitops, strutils]

type Mask* = distinct uint8

type Info* = tuple[name: seq[string], keys: seq[string]]

const MODS = [
  (name: @["SHIFT"], keys: @["shift_l", "shift_r"]),
  (name: @["CAPS"], keys: @["caps_lock"]),
  (name: @["CTRL", "CONTROL"], keys: @["control_l", "control_r"]),
  (name: @["ALT"], keys: @["alt_l", "alt_r"]),
  (name: @["MOD2"], keys: @[]),
  (name: @["MOD3"], keys: @[]),
  (name: @["SUPER", "WIN", "LOGO", "MOD4"], keys: @["super_l", "super_r"]),
  (name: @["MOD5"], keys: @[]),
]

const Empty* = 0.Mask

proc parse*(mods: string): Mask =
  for i, m in MODS:
    for n in m.name:
      if mods.contains n:
        result.uint8.setBit i

proc `$`*(mask: Mask): string =
  var mods: seq[string]
  for i, m in MODS:
    if mask.uint8.testBit i:
      mods.add m.name[0]
  mods.join "_"

proc info*(mask: Mask): seq[Info] =
  for i, m in MODS:
    if mask.uint8.testBit i:
      result.add m

proc includes*(parent, child: Mask): bool =
  child.uint8 == parent.uint8.bitand child.uint8

proc without*(parent, child: Mask): Mask =
  Mask(parent.uint8.clearMasked child.uint8)

proc `==`*(a, b: Mask): bool {.borrow.}
proc `<`*(a, b: Mask): bool {.borrow.}
