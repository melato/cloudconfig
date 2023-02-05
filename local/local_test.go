package local

import (
	"testing"

	"melato.org/cloudconfig"
)

func TestInterfaces(t *testing.T) {
	local := &BaseConfigurer{}
	var _ cloudconfig.BaseConfigurer = local
	var _ cloudconfig.BaseUserHomeDir = local
	var _ cloudconfig.BaseApplySudo = local
}
