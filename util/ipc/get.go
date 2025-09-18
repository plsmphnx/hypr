package ipc

import "encoding/json"

type Get map[string]any

func (g Get) Call() error {
	cmd := make([]string, 0, len(g))
	arg := make([]any, 0, len(g))
	for c, a := range g {
		cmd = append(cmd, "j/"+c)
		arg = append(arg, a)
	}

	res, err := Call(cmd...)
	if err != nil {
		return err
	}

	for i, r := range res {
		if err := json.Unmarshal(r, arg[i]); err != nil {
			return err
		}
	}

	return nil
}
