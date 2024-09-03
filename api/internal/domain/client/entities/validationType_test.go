package entities_test

import (
	"encoding/json"
	"testing"

	"github.com/kodmain/thetiptop/api/internal/domain/client/entities"
	"github.com/stretchr/testify/assert"
)

func TestValidationType(t *testing.T) {

	vt := entities.MailValidation

	assert.Equal(t, "mail", vt.String())

	by, err := json.Marshal(vt)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`"mail"`), by)

	var vt2 entities.ValidationType
	err = json.Unmarshal(by, &vt2)
	assert.NoError(t, err)

}
