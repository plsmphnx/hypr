package main

import (
	"math"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/plsmphnx/hypr/util/ipc"
)

type Workspace struct {
	ID        int `json:"id"`
	MonitorID int `json:"monitorID"`
	Windows   int `json:"windows"`
	i         int
}

func main() {
	var disp []string
	var prev bool
	var free int8

	for _, arg := range os.Args[1:] {
		switch arg {
		case "next":
			prev = false
		case "prev":
			prev = true
		case "used":
			free = -1
		case "free":
			free = 1
		default:
			if strings.IndexByte(arg, ' ') < 0 {
				disp = append(disp, arg+" ^")
			} else {
				disp = append(disp, arg)
			}
		}
	}

	if len(disp) == 0 {
		disp = []string{"workspace ^"}
	}

	var active Workspace
	var workspaces []Workspace
	check(ipc.Get{"activeworkspace": &active, "workspaces": &workspaces}.Call())
	slices.SortFunc(workspaces, func(a, b Workspace) int { return a.ID - b.ID })

	var monitor []Workspace
	for i, ws := range workspaces {
		ws.i = i
		if ws.MonitorID == active.MonitorID {
			if ws.ID == active.ID {
				active.i = len(monitor)
			}
			monitor = append(monitor, ws)
		}
	}

	id := max(1, min(math.MaxInt32, func() int {
		if prev {
			if active.i == 0 || free > 0 {
				tgt := monitor[0]
				if free < 0 || tgt.Windows == 0 {
					return tgt.ID
				}
				for tgt.i >= 0 && workspaces[tgt.i].ID == tgt.ID {
					tgt.i--
					tgt.ID--
				}
				return tgt.ID
			}
			return monitor[active.i-1].ID
		} else {
			if active.i == len(monitor)-1 || free > 0 {
				tgt := monitor[len(monitor)-1]
				if free < 0 || tgt.Windows == 0 {
					return tgt.ID
				}
				for tgt.i < len(workspaces) && workspaces[tgt.i].ID == tgt.ID {
					tgt.i++
					tgt.ID++
				}
				return tgt.ID
			}
			return monitor[active.i+1].ID
		}
	}()))

	var cmd ipc.Cmd
	for _, d := range disp {
		cmd.Dispatch(strings.ReplaceAll(d, "^", strconv.Itoa(id)))
	}
	check(cmd.Call())
}

func check(e error) {
	if e != nil {
		os.Stderr.WriteString(e.Error())
		os.Exit(1)
	}
}
