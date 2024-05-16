package kernel

// type Navigator interface {
// 	Navigate() (core.TraverseResult, error)
// }

// type NavigatorFunc func() (core.TraverseResult, error)

// func (fn NavigatorFunc) Navigate() (core.TraverseResult, error) {
// 	return fn()
// }

type navigationResult struct {
	err error
}

func (r *navigationResult) Error() error {
	return r.err
}
