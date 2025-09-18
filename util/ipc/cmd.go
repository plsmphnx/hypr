package ipc

import "strings"

type (
	Cmd []string
	Err []struct{ Cmd, Err string }
)

func (c *Cmd) Keyword(keyword string, args ...string) {
	c.add("/keyword", keyword, args...)
}

func (c *Cmd) Dispatch(dispatcher string, args ...string) {
	c.add("/dispatch", dispatcher, args...)
}

func (c *Cmd) add(t, n string, a ...string) {
	t += " " + strings.TrimSpace(n)
	if len(a) > 0 {
		t += " " + strings.Join(a, ",")
	}
	*c = append(*c, t)
}

func (cmd Cmd) Call() error {
	res, err := Call(cmd...)
	if err != nil {
		return err
	}

	var e Err
	for i, r := range res {
		sr := string(r)
		if sr != "ok" {
			e = append(e, struct{ Cmd, Err string }{cmd[i], sr})
		}
	}
	if len(e) > 0 {
		return e
	}

	return nil
}

func (err Err) Error() string {
	str := make([]string, len(err))
	for i, e := range err {
		str[i] = e.Cmd + "\n" + e.Err + "\n"
	}
	return strings.Join(str, "\n")
}
