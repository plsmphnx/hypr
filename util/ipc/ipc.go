package ipc

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

type IPC string

func New() IPC {
	return IPC(fmt.Sprintf("%s/hypr/%s/.socket.sock",
		os.Getenv("XDG_RUNTIME_DIR"),
		os.Getenv("HYPRLAND_INSTANCE_SIGNATURE"),
	))
}

func (ipc IPC) Call(cmds ...string) ([][]byte, error) {
	conn, err := net.DialUnix("unix", nil, &net.UnixAddr{Name: string(ipc)})
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
