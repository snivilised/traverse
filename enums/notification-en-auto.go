// Code generated by "stringer -type=Notification -linecomment -trimprefix=Notification -output notification-en-auto.go"; DO NOT EDIT.

package enums

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[NotificationUndefined-0]
}

const _Notification_name = "undefined-notification"

var _Notification_index = [...]uint8{0, 22}

func (i Notification) String() string {
	if i >= Notification(len(_Notification_index)-1) {
		return "Notification(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Notification_name[_Notification_index[i]:_Notification_index[i+1]]
}
