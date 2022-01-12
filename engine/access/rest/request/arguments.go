package request

import (
	"fmt"
)

const maxArgumentsLength = 100

type Arguments [][]byte

func (a *Arguments) Parse(raw []string) error {
	args := make([][]byte, len(raw))
	for i, a := range raw {
		arg, err := fromBase64(a)
		if err != nil {
			return fmt.Errorf("invalid argument encoding: %v", err)
		}
		args[i] = arg
	}

	if len(args) > maxArgumentsLength {
		return fmt.Errorf("too many arguments. Maximum arguments allowed: %d", maxAllowedScriptArguments)
	}

	*a = args
	return nil
}

func (a Arguments) Flow() [][]byte {
	return a
}
