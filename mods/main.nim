import std/[algorithm, cmdline, sets, strutils, tables]
import ../util/[ipc, mods, util]

type Target = object of RootObj
  locked: bool
  release: bool
  longPress: bool
  repeat: bool
  non_consuming: bool
  mouse: bool
  key: string

type Bind = object of Target
  modmask: Mask
  submap: string
  dispatcher: string
  arg: string

type Submap = ref object
  alias: string
  binds: seq[Bind]

proc keyword(tgt: Target): string =
  result = "bind"
  if tgt.locked:
    result &= "l"
  if tgt.release:
    result &= "r"
  if tgt.longPress:
    result &= "o"
  if tgt.repeat:
    result &= "e"
  if tgt.non_consuming:
    result &= "n"
  if tgt.mouse:
    result &= "m"

proc enter(cmd: var Ipc, mask: Mask, sub: string) =
  let mods = $mask
  for info in mask.info:
    for key in info.keys:
      cmd.keyword "bindr", mods, key, "submap", sub

proc exit(cmd: var Ipc, mask: Mask) =
  for info in mask.info:
    for key in info.keys:
      cmd.keyword "bindr", info.name[0], key, "submap", "reset"

proc keys(cmd: var Ipc, mask: Mask, binds: seq[Bind]) =
  let mods = $mask
  var reset = initHashSet[Target] binds.len

  for b in binds:
    if b.mouse:
      cmd.keyword b.keyword, mods, b.key, b.arg
    else:
      cmd.keyword b.keyword, mods, b.key, b.dispatcher, b.arg
    if mask == Empty or not b.repeat:
      reset.incl b

  for tgt in reset:
    var cpy = tgt
    if cpy.mouse:
      cpy.release = true
      cpy.mouse = false
    cmd.keyword cpy.keyword, mods, cpy.key, "submap", "reset"

proc submaps(cmd: var Ipc, subs: Table[Mask, Submap]) =
  var order = newSeqOfCap[Mask] subs.len
  for mask, _ in subs:
    order.add mask
  order.sort

  for i, mask in order:
    let sub = subs[mask]

    cmd.enter mask, sub.alias
    cmd.keyword "submap", sub.alias
    cmd.exit mask

    cmd.keys Empty, sub.binds

    for next in order[i + 1 ..^ 1]:
      if next.includes mask:
        let child = subs[next]
        let diff = next.without mask

        cmd.enter diff, child.alias
        cmd.keys diff, child.binds

    cmd.keyword "bindrn", "", "catchall", "submap", "reset"
    cmd.keyword "submap", "reset"

let args = commandLineParams()
var subs = initTable[Mask, Submap] args.len
for arg in args:
  let (mods, alias) = arg.cut '='
  subs[mods.parse] = Submap(alias: if alias != "": alias else: mods.strip)

ipc:
  let binds: seq[Bind]
for b in binds:
  if b.submap == "" and subs.hasKey b.modmask:
    subs[b.modmask].binds.add b

var cmd: Ipc
cmd.submaps subs
cmd.call
