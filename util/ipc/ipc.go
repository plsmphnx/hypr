package ipc

import (
	"bytes"
	"io"
	"net"
	"os"
	"path"
	"strings"
)

func Call(cmds ...string) ([][]byte, error) {
	if len(cmds) == 0 {
		return nil, nil
	}

	conn, err := net.DialUnix("unix", nil, &net.UnixAddr{Name: path.Join(
		os.Getenv("XDG_RUNTIME_DIR"), "hypr",
		os.Getenv("HYPRLAND_INSTANCE_SIGNATURE"), ".socket.sock",
	)})
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	cmd := cmds[0]
	if len(cmds) > 1 {
		cmd = "[[BATCH]]" + strings.Join(cmds, ";")
	}

	io.WriteString(conn, cmd)
	res, err := io.ReadAll(conn)
	if err != nil {
		return nil, err
	}

	return bytes.Split(res, []byte{'\n', '\n', '\n'}), nil
}
