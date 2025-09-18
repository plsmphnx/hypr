package mod

import "strings"

type (
	Mask byte

	Info struct {
		Name []string
		Mask
		Keys []string
	}
)

var modifiers = []*Info{{
	Name: []string{"SHIFT"},
	Mask: 0b00000001,
	Keys: []string{"shift_l", "shift_r"},
}, {
	Name: []string{"CAPS"},
	Mask: 0b00000010,
	Keys: []string{"caps_lock"},
}, {
	Name: []string{"CTRL", "CONTROL"},
	Mask: 0b00000100,
	Keys: []string{"control_l", "control_r"},
}, {
	Name: []string{"ALT"},
	Mask: 0b00001000,
	Keys: []string{"alt_l", "alt_r"},
}, {
	Name: []string{"MOD2"},
	Mask: 0b00010000,
	Keys: []string{},
}, {
	Name: []string{"MOD3"},
	Mask: 0b00100000,
	Keys: []string{},
}, {
	Name: []string{"SUPER", "WIN", "LOGO", "MOD4"},
	Mask: 0b01000000,
	Keys: []string{"super_l", "super_r"},
}, {
	Name: []string{"MOD5"},
	Mask: 0b10000000,
	Keys: []string{},
}}

func Parse(mods string) Mask {
	var mask Mask
	for _, mod := range modifiers {
		for _, name := range mod.Name {
			if strings.Contains(mods, name) {
				mask |= mod.Mask
			}
		}
	}
	return mask
}

func (m Mask) String() string {
	var mods []string
	for _, mod := range modifiers {
		if m&mod.Mask != 0 {
			mods = append(mods, mod.Name[0])
		}
	}
	return strings.Join(mods, "_")
}

func (m Mask) Info() []*Info {
	var mods []*Info
	for _, mod := range modifiers {
		if m&mod.Mask != 0 {
			mods = append(mods, mod)
		}
	}
	return mods
}
