package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMustReadConfigFile(t *testing.T) {
	MustReadConfigFile("../../test/config.yaml")
}

func TestWriteConfigFile(t *testing.T) {
	assert.NoError(t, WriteConfigFile(Config{}, "../../test/config-empty.yaml"))
}
