package local

import (
	"testing"

	"melato.org/cloudinit"
)

func TestInterfaces(t *testing.T) {
	local := &BaseConfigurer{}
	var _ cloudinit.BaseConfigurer = local
	var _ cloudinit.BaseUserHomeDir = local
	var _ cloudinit.BaseApplySudo = local
}
