package actions_on_google

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnum(t *testing.T) {
	assert := assert.New(t)

	data := IVT_DATETIME
	assert.Equal(data.Intent(), "actions.intent.DATETIME")
}
