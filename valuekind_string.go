// Code generated by "stringer -type=ValueKind"; DO NOT EDIT.

package argmapper

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ValueInvalid-0]
	_ = x[ValueNamed-1]
	_ = x[ValueTyped-2]
}

const _ValueKind_name = "ValueInvalidValueNamedValueTyped"

var _ValueKind_index = [...]uint8{0, 12, 22, 32}

func (i ValueKind) String() string {
	if i >= ValueKind(len(_ValueKind_index)-1) {
		return "ValueKind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ValueKind_name[_ValueKind_index[i]:_ValueKind_index[i+1]]
}
