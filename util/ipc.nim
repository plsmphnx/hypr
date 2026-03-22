import std/[envvars, json, macros, net, paths, strutils]
import ./util

type Ipc* = distinct seq[string]
type Err* = object of CatchableError

proc call(cmds: varargs[string]): seq[string] =
  if cmds.len == 0:
    return

  let sock = newUnixSocket:
    "XDG_RUNTIME_DIR".getEnv.Path / "hypr".Path /
      "HYPRLAND_INSTANCE_SIGNATURE".getEnv.Path / ".socket.sock".Path
  defer:
    sock.close

  sock.send:
    if cmds.len > 1:
      "[[BATCH]]" & cmds.join ";"
    else:
      cmds[0]
  sock.recvAll.split "\n\n\n"

proc add(ipc: var Ipc, t, n: string, args: varargs[string]) =
  seq[string](ipc).add:
    if args.len == 0:
      t & ' ' & n.strip
    else:
      t & ' ' & n.strip & ' ' & args.join ","

proc keyword*(ipc: var Ipc, keyword: string, args: varargs[string]) =
  ipc.add "/keyword", keyword, args

proc dispatch*(ipc: var Ipc, dispatcher: string, args: varargs[string]) =
  ipc.add "/dispatch", dispatcher, args

proc call*(ipc: Ipc) =
  let cmds = seq[string](ipc)
  var err: string
  for i, r in cmds.call:
    if r != "ok":
      err &= '\n' & cmds[i] & '\n' & r & '\n'

  if err != "":
    raise Err.newException err

macro ipc*(calls: untyped): untyped =
  result = newStmtList()
  var names = newNimNode nnkBracket

  calls.expectKind nnkStmtList
  for section in calls:
    section.expectKind {nnkLetSection, nnkVarSection}
    for ident in section:
      ident[2].expectKind nnkEmpty
      names.add newLit "j/" & $ident[0]

  let res = genSym()
  result.add quote do:
    let `res` = `names`.call

  var i: int
  for section in calls:
    for ident in section:
      let t = ident[1]
      ident[2] = quote:
        `res`[`i`].parseJson.to `t`
      i += 1
    result.add section
