// Code generated by "stringer -type=Subscription -linecomment -trimprefix=Subscribe -output subscription-en-auto.go"; DO NOT EDIT.

package enums

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SubscribeUndefined-0]
	_ = x[SubscribeFiles-1]
	_ = x[SubscribeDirectories-2]
	_ = x[SubscribeDirectoriesWithFiles-3]
	_ = x[SubscribeUniversal-4]
}

const _Subscription_name = "Undefinedsubscribe-filessubscribe-directoriessubscribe-directories-with-filessubscribe-to-everything"

var _Subscription_index = [...]uint8{0, 9, 24, 45, 77, 100}

func (i Subscription) String() string {
	if i >= Subscription(len(_Subscription_index)-1) {
		return "Subscription(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Subscription_name[_Subscription_index[i]:_Subscription_index[i+1]]
}
