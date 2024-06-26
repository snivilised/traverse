package i18n

import (
	"github.com/snivilised/extendio/i18n"
)

type UsingConfigFileTemplData struct {
	traverseTemplData
	ConfigFileName string
}

func (td UsingConfigFileTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "using-config-file",
		Description: "Message to indicate which config is being used",
		Other:       "Using config file: {{.ConfigFileName}}",
	}
}
