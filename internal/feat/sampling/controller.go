package sampling

import (
	"io/fs"

	nef "github.com/snivilised/nefilim"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type controller struct {
	o      *pref.SamplingOptions
	filter core.ChildTraverseFilter
}

func (p *controller) Role() enums.Role {
	return enums.RoleSampler
}

func (p *controller) Next(_ core.Servant, _ types.Inspection) (bool, error) {
	return true, nil
}

func (p *controller) sample(result []fs.DirEntry, _ error,
	_ fs.ReadDirFS, _ string,
) ([]fs.DirEntry, error) {
	files, folders := nef.Separate(result)

	return union(&readResult{
		files:   files,
		folders: folders,
		o:       p.o,
	}), nil
}

type (
	samplerFunc func(n uint, entries []fs.DirEntry) []fs.DirEntry
)

type readResult struct {
	files   []fs.DirEntry
	folders []fs.DirEntry
	o       *pref.SamplingOptions
}

func union(r *readResult) []fs.DirEntry {
	noOfFiles := lo.Ternary(r.o.NoOf.Files == 0,
		uint(len(r.files)), r.o.NoOf.Files,
	)

	both := lo.Ternary(
		r.o.InReverse, last, first,
	)(noOfFiles, r.files)

	noOfFolders := lo.Ternary(r.o.NoOf.Folders == 0,
		uint(len(r.folders)), r.o.NoOf.Folders,
	)
	both = append(both, lo.Ternary(
		r.o.InReverse, last, first,
	)(noOfFolders, r.folders)...)

	return both
}

func first(n uint, entries []fs.DirEntry) []fs.DirEntry {
	return entries[:(min(n, uint(len(entries))))]
}

func last(n uint, entries []fs.DirEntry) []fs.DirEntry {
	return entries[uint(len(entries))-(min(n, uint(len(entries)))):]
}
