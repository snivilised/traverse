// Code generated by "stringer -type=SkipTraversal -linecomment -trimprefix=Skip -output skip-traversal-en-auto.go"; DO NOT EDIT.

package enums

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SkipNoneTraversal-0]
	_ = x[SkipDirTraversal-1]
	_ = x[SkipAllTraversal-2]
}

const _SkipTraversal_name = "skip-noneskip-dirskip-all"

var _SkipTraversal_index = [...]uint8{0, 9, 17, 25}

func (i SkipTraversal) String() string {
	if i >= SkipTraversal(len(_SkipTraversal_index)-1) {
		return "SkipTraversal(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _SkipTraversal_name[_SkipTraversal_index[i]:_SkipTraversal_index[i+1]]
}