package flag

type Flags struct {
	Locked       bool `json:"locked"`
	Release      bool `json:"release"`
	LongPress    bool `json:"longPress"`
	Repeat       bool `json:"repeat"`
	NonConsuming bool `json:"non_consuming"`
	Mouse        bool `json:"mouse"`
}

func (f *Flags) String() string {
	str := "bind"
	if f.Locked {
		str += "l"
	}
	if f.Release {
		str += "r"
	}
	if f.LongPress {
		str += "o"
	}
	if f.Repeat {
		str += "e"
	}
	if f.NonConsuming {
		str += "n"
	}
	if f.Mouse {
		str += "m"
	}
	return str
}
