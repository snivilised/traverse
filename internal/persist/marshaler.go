package persist

import (
	ejson "encoding/json"
	"io/fs"

	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/internal/enclave"
	"github.com/snivilised/agenor/internal/opts/json"
	"github.com/snivilised/agenor/pref"
	nef "github.com/snivilised/nefilim"
)

type (
	// TamperFunc provides a way for unit tests to modify the JSON before
	// it is un-marshaled. The unit tests marshal a default JSON object
	// instance, so a TamperFunc is used to allow modification of that
	// default. Typically a single test will focus on a single field,
	// so that the TamperFunc is expected to only update 1 of the members at a
	// time.
	TamperFunc func(result *MarshalResult)

	MarshalRequest struct {
		Active *core.ActiveState
		O      *pref.Options
		Path   string
		Perm   fs.FileMode
		FS     nef.WriteFileFS
	}

	MarshalResult struct {
		Active *core.ActiveState
		JO     *json.Options
	}

	UnmarshalRequest struct {
		Restore *enclave.RestoreState
	}

	UnmarshalResult struct {
		Active *core.ActiveState
		JO     *json.Options
		O      *pref.Options
	}

	Comparison struct {
		JO *json.Options
		O  *pref.Options
	}
)

func Marshal(request *MarshalRequest) (*MarshalResult, error) {
	jo := ToJSON(request.O)
	result := &MarshalResult{
		JO:     jo,
		Active: request.Active.Clone(),
	}

	data, err := ejson.MarshalIndent(
		result,
		JSONMarshalNoPrefix, JSONMarshal2SpacesIndent,
	)

	if err != nil {
		return nil, err
	}

	if err := (&Comparison{
		O:  request.O,
		JO: jo,
	}).Equals(); err != nil {
		return result, err
	}

	return result, request.FS.WriteFile(request.Path, data, request.Perm)
}

func Unmarshal(request *UnmarshalRequest,
	tampers ...TamperFunc,
) (*UnmarshalResult, error) {
	bytes, err := request.Restore.FS.ReadFile(request.Restore.Path)

	if err != nil {
		return nil, err
	}

	var (
		mr MarshalResult
	)

	if err := ejson.Unmarshal(bytes, &mr); err != nil {
		return nil, err
	}

	for _, fn := range tampers {
		fn(&mr)
	}

	result := UnmarshalResult{
		O:      FromJSON(mr.JO),
		Active: mr.Active,
		JO:     mr.JO,
	}

	return &result, (&Comparison{
		O:  result.O,
		JO: result.JO,
	}).Equals()
}
