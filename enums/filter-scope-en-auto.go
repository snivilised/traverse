// Code generated by "stringer -type=FilterScope -linecomment -trimprefix=Scope -output filter-scope-en-auto.go"; DO NOT EDIT.

package enums

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ScopeUndefined-0]
	_ = x[ScopeRoot-1]
	_ = x[ScopeTop-2]
	_ = x[ScopeLeaf-4]
	_ = x[ScopeIntermediate-8]
	_ = x[ScopeFile-16]
	_ = x[ScopeFolder-32]
	_ = x[ScopeCustom-64]
}

const (
	_FilterScope_name_0 = "undefined-scoperoot-scopetop-scope"
	_FilterScope_name_1 = "leaf-scope"
	_FilterScope_name_2 = "intermediate-scope"
	_FilterScope_name_3 = "file-scope"
	_FilterScope_name_4 = "folder-scope"
	_FilterScope_name_5 = "custom-scope"
)

var (
	_FilterScope_index_0 = [...]uint8{0, 15, 25, 34}
)

func (i FilterScope) String() string {
	switch {
	case i <= 2:
		return _FilterScope_name_0[_FilterScope_index_0[i]:_FilterScope_index_0[i+1]]
	case i == 4:
		return _FilterScope_name_1
	case i == 8:
		return _FilterScope_name_2
	case i == 16:
		return _FilterScope_name_3
	case i == 32:
		return _FilterScope_name_4
	case i == 64:
		return _FilterScope_name_5
	default:
		return "FilterScope(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}