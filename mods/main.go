package main

import (
	"os"
	"slices"
	"strings"

	"github.com/plsmphnx/hypr/util/ipc"
	"github.com/plsmphnx/hypr/util/mod"
)

type (
	Target struct {
		Locked       bool   `json:"locked"`
		Release      bool   `json:"release"`
		LongPress    bool   `json:"longPress"`
		Repeat       bool   `json:"repeat"`
		NonConsuming bool   `json:"non_consuming"`
		Mouse        bool   `json:"mouse"`
		Key          string `json:"key"`
	}

	Bind struct {
		Target
		Modmask    mod.Mask `json:"modmask"`
		Submap     string   `json:"submap"`
		Dispatcher string   `json:"dispatcher"`
		Arg        string   `json:"arg"`
	}

	Submap struct {
		Alias string
		Binds []*Bind
	}
)

func main() {
	submaps := make(map[mod.Mask]*Submap, len(os.Args)-1)
	for _, arg := range os.Args[1:] {
		submap := &Submap{}
		mods, alias, _ := strings.Cut(arg, "=")
		if alias != "" {
			submap.Alias = alias
		} else {
			submap.Alias = strings.TrimSpace(mods)
		}
		submaps[mod.Parse(mods)] = submap
	}

	var binds []*Bind
	check(ipc.Get{"binds": &binds}.Call())
	for _, bind := range binds {
		if submap, ok := submaps[bind.Modmask]; ok && bind.Submap == "" {
			submap.Binds = append(submap.Binds, bind)
		}
	}

	var cmd ipc.Cmd
	Submaps(&cmd, submaps)
	check(cmd.Call())
}

func Submaps(c *ipc.Cmd, submaps map[mod.Mask]*Submap) {
	order := make([]mod.Mask, 0, len(submaps))
	for mask := range submaps {
		order = append(order, mask)
	}
	slices.Sort(order)

	for i, mask := range order {
		submap := submaps[mask]

		Enter(c, mask, submap.Alias)
		c.Keyword("submap", submap.Alias)
		Exit(c, mask)

		Binds(c, 0, submap.Binds)

		for _, next := range order[i+1:] {
			if mask&next == mask {
				child := submaps[next]
				diff := next &^ mask

				Enter(c, diff, child.Alias)
				Binds(c, diff, child.Binds)
			}
		}

		c.Keyword("bindrn", "", "catchall", "submap", "reset")
		c.Keyword("submap", "reset")
	}
}

func Enter(c *ipc.Cmd, mask mod.Mask, submap string) {
	mods := mask.String()
	for _, i := range mask.Info() {
		for _, key := range i.Keys {
			c.Keyword("bindr", mods, key, "submap", submap)
		}
	}
}

func Exit(c *ipc.Cmd, mask mod.Mask) {
	for _, i := range mask.Info() {
		for _, key := range i.Keys {
			c.Keyword("bindr", i.Name[0], key, "submap", "reset")
		}
	}
}

func Binds(c *ipc.Cmd, mask mod.Mask, binds []*Bind) {
	mods := mask.String()
	reset := make(map[Target]struct{}, len(binds))

	for _, bind := range binds {
		if bind.Mouse {
			c.Keyword(bind.Keyword(), mods, bind.Key, bind.Arg)
		} else {
			c.Keyword(bind.Keyword(), mods, bind.Key, bind.Dispatcher, bind.Arg)
		}
		if mask == 0 || !bind.Repeat {
			reset[bind.Target] = struct{}{}
		}
	}

	for t := range reset {
		if t.Mouse {
			t.Release = true
			t.Mouse = false
		}
		c.Keyword(t.Keyword(), mods, t.Key, "submap", "reset")
	}
}

func (t *Target) Keyword() string {
	str := "bind"
	if t.Locked {
		str += "l"
	}
	if t.Release {
		str += "r"
	}
	if t.LongPress {
		str += "o"
	}
	if t.Repeat {
		str += "e"
	}
	if t.NonConsuming {
		str += "n"
	}
	if t.Mouse {
		str += "m"
	}
	return str
}

func check(e error) {
	if e != nil {
		os.Stderr.WriteString(e.Error())
		os.Exit(1)
	}
}
