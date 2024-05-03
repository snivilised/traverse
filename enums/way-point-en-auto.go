// Code generated by "stringer -type=WayPoint -linecomment -trimprefix=WayPoint -output way-point-en-auto.go"; DO NOT EDIT.

package enums

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[WayPointUndefined-0]
	_ = x[WayPointWake-1]
	_ = x[WayPointSleep-2]
	_ = x[WayPointToggle-3]
}

const _WayPoint_name = "undefined-way-pointwake-from-hibernationsleeptoggle-way-point"

var _WayPoint_index = [...]uint8{0, 19, 40, 45, 61}

func (i WayPoint) String() string {
	if i >= WayPoint(len(_WayPoint_index)-1) {
		return "WayPoint(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _WayPoint_name[_WayPoint_index[i]:_WayPoint_index[i+1]]
}