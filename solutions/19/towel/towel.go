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

// Given the (remaining) pattern to create, returns an array of atoms used to create the pattern
// or nil if the pattern is impossible
func (towel TowelCollection) isPatternValidRecursive(remainingPattern string) []string {
	if len(remainingPattern) == 0 {
		return make([]string, 0)
	}

	for _, atom := range towel.towelAtoms {
		if patternLessAtom, hasPrefix := strings.CutPrefix(remainingPattern, atom); hasPrefix {
			// slog.Debug("found prefix", "remaining pattern", remainingPattern, "prefix atom", atom, "pattern less atom", patternLessAtom)
			if constructingAtoms := towel.isPatternValidRecursive(patternLessAtom); constructingAtoms != nil {
				constructingAtoms = append(constructingAtoms, atom)
				return constructingAtoms
			}
		}
	}
	return nil
}
func (towel TowelCollection) IsPatternValid(pattern string) bool {
	constructingAtoms := towel.isPatternValidRecursive(pattern)
	if constructingAtoms != nil {
		slog.Debug("constructed pattern successfully", "pattern", pattern, "atoms", constructingAtoms)
		return true
	} else {
		return false
	}
}

