import std/[algorithm, cmdline, strutils, sugar]
import ../util/ipc

type Workspace = object
  id: int
  monitorID: int
  windows: int

type Filter = enum
  None
  Used
  Free

var disp: seq[string]
var prev: bool
var free: Filter

for arg in commandLineParams():
  case arg
  of "next":
    prev = false
  of "prev":
    prev = true
  of "used":
    free = Used
  of "free":
    free = Free
  else:
    disp.add:
      if not arg.contains Whitespace:
        arg & " ^"
      else:
        arg

if disp.len == 0:
  disp.add "workspace ^"

ipc:
  let activeworkspace: Workspace
  var workspaces: seq[Workspace]
workspaces.sort (a, b) => a.id - b.id

var active: int
var monitor: seq[int]
for i, this in workspaces:
  if this.monitorID == activeworkspace.monitorID:
    if this.id == activeworkspace.id:
      active = monitor.len
    monitor.add i

let id =
  if prev:
    if active == 0 or free == Free:
      var i = monitor[0]
      var id = workspaces[i].id
      if free != Used and workspaces[i].windows != 0:
        while i >= 0 and workspaces[i].id == id:
          i -= 1
          id -= 1
      id
    else:
      workspaces[monitor[active - 1]].id
  else:
    if active == monitor.high or free == Free:
      var i = monitor[monitor.high]
      var id = workspaces[i].id
      if free != Used and workspaces[i].windows != 0:
        while i < workspaces.len and workspaces[i].id == id:
          i += 1
          id += 1
      id
    else:
      workspaces[monitor[active + 1]].id
let ws = $id.clamp(1, int32.high)

var cmd: Ipc
for d in disp:
  cmd.dispatch:
    d.replace "^", ws
cmd.call
