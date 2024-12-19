package towel

import (
	"log/slog"
	"strings"
)

type TowelCollection struct {
	towelAtoms []string
}

func NewTowelCollection(towelAtoms []string) TowelCollection {
	return TowelCollection{
		towelAtoms: towelAtoms,
	}
}

