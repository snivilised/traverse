package tv

import (
	"io/fs"

	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/resume"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type extent interface {
	using() *pref.Using
	was() *pref.Was
	plugin(types.Mediator) types.Plugin
	options(...pref.Option) (*pref.Options, error)
	navFS() fs.FS
	resFS() fs.FS
	complete() bool
}

type fileSystems struct {
	nas fs.FS
	res fs.FS
}

type baseExtent struct {
	fsys fileSystems
}

func (ex *baseExtent) navFS() fs.FS {
	return ex.fsys.nas
}

func (ex *baseExtent) resFS() fs.FS {
	return ex.fsys.nas
}

type primeExtent struct {
	baseExtent
	u *pref.Using
}

func (ex *primeExtent) using() *pref.Using {
	return ex.u
}

func (ex *primeExtent) was() *pref.Was {
	return nil
}

func (ex *primeExtent) plugin(types.Mediator) types.Plugin {
	return nil
}

func (ex *primeExtent) options(settings ...pref.Option) (*pref.Options, error) {
	return pref.Get(settings...)
}

func (ex *primeExtent) complete() bool {
	return true
}

type resumeExtent struct {
	baseExtent
	w      *pref.Was
	loaded *pref.LoadInfo
}

func (ex *resumeExtent) using() *pref.Using {
	return &ex.w.Using
}

func (ex *resumeExtent) was() *pref.Was {
	return ex.w
}

func (ex *resumeExtent) plugin(mediator types.Mediator) types.Plugin {
	return &resume.Plugin{
		BasePlugin: kernel.BasePlugin{
			Mediator: mediator,
		},
	}
}

func (ex *resumeExtent) options(settings ...pref.Option) (*pref.Options, error) {
	loaded, err := resume.Load(ex.fsys.res, ex.w.From, settings...)
	ex.loaded = loaded

	// get the resume point from the resume persistence file
	// then set up hibernation with this defined as a hibernation
	// filter.
	//
	return loaded.O, err
}

func (ex *resumeExtent) complete() bool {
	// "NOT-IMPL: resumeExtent.complete -> the strategy knows this"
	// ===> send to plugin?
	//
	return true
}
