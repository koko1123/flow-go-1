package ptrie

import "github.com/koko1123/flow-go-1/ledger"

type ErrMissingPath struct {
	Paths []ledger.Path
}

func (e ErrMissingPath) Error() string {
	str := "paths are missing: \n"
	for _, k := range e.Paths {
		str += "\t" + k.String() + "\n"
	}
	return str
}
