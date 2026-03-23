import std/[net, strutils]

proc newUnixSocket*(): Socket =
  newSocket AF_UNIX, SOCK_STREAM, IPPROTO_IP

proc recvAll*(socket: Socket): string =
  var res: string
  while socket.recv(res, 8192) > 0:
    result.add res

proc cut*(str: string, chr: char): (string, string) =
  let i = str.find chr
  if i >= 0:
    (str[0 .. i - 1], str[i + 1 .. ^1])
  else:
    (str, "")
