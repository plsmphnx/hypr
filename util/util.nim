import std/[net, paths, strutils]

proc newUnixSocket*(path: Path): Socket =
  let socket = newSocket(AF_UNIX, SOCK_STREAM, IPPROTO_IP)
  socket.connectUnix $path
  socket

const BUFFER = 8192

proc recvAll*(socket: Socket): string =
  var res: string
  var read = BUFFER
  while read == BUFFER:
    read = socket.recv(res, BUFFER)
    result.add res

proc cut*(str: string, chr: char): (string, string) =
  let i = str.find chr
  if i >= 0:
    (str[0 .. i - 1], str[i + 1 .. ^1])
  else:
    (str, "")
