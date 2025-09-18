package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/plsmphnx/hypr/util/flag"
	"github.com/plsmphnx/hypr/util/ipc"
	"github.com/plsmphnx/hypr/util/mod"
)

type (
	Target struct {
		flag.Flags
		Key string `json:"key"`
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

	var cmd ipc.Cmd
	ipc := ipc.New()

	var binds []*Bind
	check(json.Unmarshal(must(ipc.Call("j/binds"))[0], &binds))
	for _, bind := range binds {
		if submap, ok := submaps[bind.Modmask]; ok && bind.Submap == "" {
			submap.Binds = append(submap.Binds, bind)
		}
	}

	Submaps(&cmd, submaps)
	check(ipc.Exec(cmd))
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
	for _, bind := range binds {
		if bind.Flags.Mouse {
			c.Keyword(bind.String(), mods, bind.Key, bind.Arg)
		} else {
			c.Keyword(bind.String(), mods, bind.Key, bind.Dispatcher, bind.Arg)
		}
	}

	dedup := make(map[Target]struct{}, len(binds))
	for _, bind := range binds {
		if _, ok := dedup[bind.Target]; !ok {
			dedup[bind.Target] = struct{}{}

			flags := bind.Flags
			if flags.Mouse {
				flags.Release = true
				flags.Mouse = false
			}

			c.Keyword(flags.String(), mods, bind.Key, "submap", "reset")
		}
	}
}

func check(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}

func must[T any](t T, e error) T {
	check(e)
	return t
}
