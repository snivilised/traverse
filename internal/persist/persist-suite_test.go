package persist_test

import (
	_ "embed"
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/traverse/internal/persist"
	"github.com/snivilised/traverse/pref"
)

func TestPersist(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Persist Suite")
}

const (
	NoOverwrite = false
	source      = "json/unmarshal"
	destination = "json/marshal"
	tempFile    = "test-state-marshal.TEMP.json"
	restoreFile = "test-restore.DEFAULT.json"
	home        = "/home"
	foo         = "foo"
	flac        = "*.flac"
	bar         = "*.bar"
)

var (
	//go:embed data/test-restore.DEFAULT.json
	content []byte
)

type (
	persistTE struct {
		given string
	}

	tampered struct {
		o      *pref.Options
		result *persist.MarshalResult
	}

	checkerFunc func(entry *checkerTE, err error) error

	checkerTE struct {
		field string
		// checker ensures the resultant error is reporting the correct field.
		// If failure, then an error is returned, indicating name of the actual
		// field reported, vs the expected.
		checker checkerFunc
	}

	marshalTE struct {
		persistTE
		*checkerTE
		// option defines a single option to be defined for the unit test. When
		// a test case wants to test an optional option in pref.Options (ie it
		// is a pointer), then that test case will not define this option. Instead
		// it will define the tweak function to contain the corresponding member
		// on the json instance, such that the pref.member is nil and json.member
		// is not til, thereby triggering an unequal error.
		option func() pref.Option

		// tweak allows a test case to change json.Options to provoke unequal error
		tweak persist.TamperFunc
	}

	wrongUnequalError struct {
		field string
		err   error
	}
)

func (e wrongUnequalError) Error() string {
	return fmt.Sprintf("wrong unequal error (%v) type for field: %q",
		e.err, e.field,
	)
}
