// Code generated by "stringer --type Direction"; DO NOT EDIT.

package gridutils

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[DIRECTION_UP-0]
	_ = x[DIRECTION_RIGHT-1]
	_ = x[DIRECTION_DOWN-2]
	_ = x[DIRECTION_LEFT-3]
}

const _Direction_name = "DIRECTION_UPDIRECTION_RIGHTDIRECTION_DOWNDIRECTION_LEFT"

var _Direction_index = [...]uint8{0, 12, 27, 41, 55}

func (i Direction) String() string {
	if i < 0 || i >= Direction(len(_Direction_index)-1) {
		return "Direction(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Direction_name[_Direction_index[i]:_Direction_index[i+1]]
}