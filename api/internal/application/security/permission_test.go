package security_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kodmain/thetiptop/api/internal/application/security"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/jwt"
	"github.com/stretchr/testify/assert"
)

type MockEntity struct {
	OwnerID string
	Public  bool
}

func (e *MockEntity) GetOwnerID() string {
	return e.OwnerID
}

func (e *MockEntity) IsPublic() bool {
	return e.Public
}

func CustomRule(p *security.UserAccess, r database.Entity) bool {
	// Example custom rule: allow if the entity owner ID ends with "xyz"
	return r.GetOwnerID() == "owner-xyz"
}

func TestGetCredentialID(t *testing.T) {
	p := &security.UserAccess{CredentialID: "test-id"}
	assert.Equal(t, aws.String("test-id"), p.GetCredentialID())

	p = &security.UserAccess{}
	assert.Nil(t, p.GetCredentialID())
}

func TestIsGranted(t *testing.T) {
	tests := []struct {
		name     string
		userRole string
		roles    []string
		expected bool
	}{
		{"Role granted", "admin", []string{"admin"}, true},
		{"Role not granted", "user", []string{"admin"}, false},
		{"Multiple roles granted", "admin", []string{"user", "admin"}, true},
		{"Multiple roles not granted", "guest", []string{"admin"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &security.UserAccess{Role: tt.userRole}
			result := p.IsGranted(tt.roles...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCanRead(t *testing.T) {
	tests := []struct {
		name       string
		userAccess *security.UserAccess
		entity     database.Entity
		expected   bool
	}{
		{"Public entity", &security.UserAccess{}, &MockEntity{Public: true}, true},
		{"Owner access", &security.UserAccess{CredentialID: "owner-id"}, &MockEntity{OwnerID: "owner-id"}, true},
		{"Access denied", &security.UserAccess{}, &MockEntity{OwnerID: "other-id"}, false},
		{"Custom rule granted", &security.UserAccess{}, &MockEntity{OwnerID: "owner-xyz"}, true},
		{"Custom rule denied", &security.UserAccess{}, &MockEntity{OwnerID: "owner-abc"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.userAccess.CanRead(tt.entity, CustomRule)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCanCreate(t *testing.T) {
	tests := []struct {
		name       string
		userAccess *security.UserAccess
		entity     database.Entity
		expected   bool
	}{
		{"Owner access", &security.UserAccess{CredentialID: "owner-id"}, &MockEntity{OwnerID: "owner-id"}, true},
		{"Access denied", &security.UserAccess{}, &MockEntity{OwnerID: "other-id"}, false},
		{"Custom rule granted", &security.UserAccess{}, &MockEntity{OwnerID: "owner-xyz"}, true},
		{"Custom rule denied", &security.UserAccess{}, &MockEntity{OwnerID: "owner-abc"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.userAccess.CanCreate(tt.entity, CustomRule)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCanUpdate(t *testing.T) {
	tests := []struct {
		name       string
		userAccess *security.UserAccess
		entity     database.Entity
		expected   bool
	}{
		{"Owner access", &security.UserAccess{CredentialID: "owner-id"}, &MockEntity{OwnerID: "owner-id"}, true},
		{"Access denied", &security.UserAccess{}, &MockEntity{OwnerID: "other-id"}, false},
		{"Custom rule granted", &security.UserAccess{}, &MockEntity{OwnerID: "owner-xyz"}, true},
		{"Custom rule denied", &security.UserAccess{}, &MockEntity{OwnerID: "owner-abc"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.userAccess.CanUpdate(tt.entity, CustomRule)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCanDelete(t *testing.T) {
	tests := []struct {
		name       string
		userAccess *security.UserAccess
		entity     database.Entity
		expected   bool
	}{
		{"Owner access", &security.UserAccess{CredentialID: "owner-id"}, &MockEntity{OwnerID: "owner-id"}, true},
		{"Access denied", &security.UserAccess{}, &MockEntity{OwnerID: "other-id"}, false},
		{"Custom rule granted", &security.UserAccess{}, &MockEntity{OwnerID: "owner-xyz"}, true},
		{"Custom rule denied", &security.UserAccess{}, &MockEntity{OwnerID: "owner-abc"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.userAccess.CanDelete(tt.entity, CustomRule)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewUserAccess(t *testing.T) {
	token := &jwt.Token{
		ID:   "test-id",
		Data: map[string]interface{}{"role": "admin"},
	}
	p := security.NewUserAccess(token)
	assert.Equal(t, "test-id", p.CredentialID)
	assert.Equal(t, "admin", p.Role)
}

func TestNewUserAccess_NoRole(t *testing.T) {
	token := &jwt.Token{ID: "test-id"}
	p := security.NewUserAccess(token)
	assert.Equal(t, "test-id", p.CredentialID)
	assert.Equal(t, "", p.Role)
}
