package filter

import (
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/internal/enclave"
	"github.com/snivilised/agenor/internal/third/lo"
	"github.com/snivilised/agenor/pref"
)

type (
	scheme interface {
		create() error
		init(pi *enclave.PluginInit, crate *enclave.Crate)
		next(servant core.Servant, inspection enclave.Inspection) (bool, error)
	}
)

type common struct {
	o     *pref.Options
	crate *enclave.Crate
}

func (f *common) init(_ *enclave.PluginInit, crate *enclave.Crate) {
	f.crate = crate
}

func newScheme(o *pref.Options) scheme {
	c := common{o: o}
	primary, nanny := binary(o)

	if primary != nil && nanny != nil {
		return &hybridScheme{
			common:  c,
			primary: primary,
			nanny:   nanny,
		}
	}

	if primary != nil {
		return primary
	}

	return nanny
}

func binary(o *pref.Options) (primary, nanny scheme) {
	c := common{o: o}

	primary = unary(c)
	nanny = lo.TernaryF(o.Filter.IsChildFilteringActive(),
		func() scheme {
			return &nannyScheme{
				common: c,
			}
		},
		func() scheme {
			return nil
		},
	)

	if nanny == nil {
		return primary, nil
	}

	return primary, nanny
}

func unary(c common) scheme {
	if c.o.Filter.IsCustomFilteringActive() {
		return &customScheme{
			common: c,
		}
	}

	if c.o.Filter.IsNodeFilteringActive() {
		return &nativeScheme{
			common: c,
		}
	}

	if c.o.Filter.IsSampleFilteringActive() {
		return &samplerScheme{
			common: c,
		}
	}

	return nil // only nanny is active
}
