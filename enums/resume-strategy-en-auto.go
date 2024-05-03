// Code generated by "stringer -type=ResumeStrategy -linecomment -trimprefix=ResumeStrategy -output resume-strategy-en-auto.go"; DO NOT EDIT.

package enums

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ResumeStrategyUndefined-0]
	_ = x[ResumeStrategySpawn-1]
	_ = x[ResumeStrategyFastward-2]
}

const _ResumeStrategy_name = "undefined-resume-strategyspawn-resume-strategyfastward-resume-strategy"

var _ResumeStrategy_index = [...]uint8{0, 25, 46, 70}

func (i ResumeStrategy) String() string {
	if i >= ResumeStrategy(len(_ResumeStrategy_index)-1) {
		return "ResumeStrategy(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ResumeStrategy_name[_ResumeStrategy_index[i]:_ResumeStrategy_index[i+1]]
}