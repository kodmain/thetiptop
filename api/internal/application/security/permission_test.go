package security_test

import (
	"testing"

	"github.com/kodmain/thetiptop/api/internal/application/security"
	"github.com/kodmain/thetiptop/api/internal/domain/user/entities"
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

func CustomRule(p *security.UserAccess, r entities.Entity) bool {
	// Example custom rule: allow if user has a specific role or if the entity owner ID ends with "xyz"
	return p.IsGranted("special_role") || r.GetOwnerID() == "owner-xyz"
}

func TestIsGranted(t *testing.T) {
	tests := []struct {
		name      string
		userRoles []string
		roles     []string
		expected  bool
	}{
		{"Role granted", []string{"admin"}, []string{"admin"}, true},
		{"Role not granted", []string{"user"}, []string{"admin"}, false},
		{"Multiple roles granted", []string{"user", "admin"}, []string{"admin"}, true},
		{"Multiple roles not granted", []string{"user", "guest"}, []string{"admin"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &security.UserAccess{Role: tt.userRoles}
			result := p.IsGranted(tt.roles...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCanRead(t *testing.T) {
	tests := []struct {
		name       string
		userAccess *security.UserAccess
		entity     entities.Entity
		expected   bool
	}{
		{"Admin access", &security.UserAccess{Role: []string{"admin"}}, &MockEntity{}, true},
		{"Public entity", &security.UserAccess{}, &MockEntity{Public: true}, true},
		{"Owner access", &security.UserAccess{CredentialID: "owner-id"}, &MockEntity{OwnerID: "owner-id"}, true},
		{"Access denied", &security.UserAccess{}, &MockEntity{OwnerID: "other-id"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.userAccess.CanRead(tt.entity)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCanCreate(t *testing.T) {
	tests := []struct {
		name       string
		userAccess *security.UserAccess
		entity     entities.Entity
		expected   bool
	}{
		{"Admin access", &security.UserAccess{Role: []string{"admin"}}, &MockEntity{}, true},
		{"Owner access", &security.UserAccess{CredentialID: "owner-id"}, &MockEntity{OwnerID: "owner-id"}, true},
		{"Access denied", &security.UserAccess{}, &MockEntity{OwnerID: "other-id"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.userAccess.CanCreate(tt.entity)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCanUpdate(t *testing.T) {
	tests := []struct {
		name       string
		userAccess *security.UserAccess
		entity     entities.Entity
		expected   bool
	}{
		{"Admin access", &security.UserAccess{Role: []string{"admin"}}, &MockEntity{}, true},
		{"Owner access", &security.UserAccess{CredentialID: "owner-id"}, &MockEntity{OwnerID: "owner-id"}, true},
		{"Access denied", &security.UserAccess{}, &MockEntity{OwnerID: "other-id"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.userAccess.CanUpdate(tt.entity)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCanDelete(t *testing.T) {
	tests := []struct {
		name       string
		userAccess *security.UserAccess
		entity     entities.Entity
		expected   bool
	}{
		{"Admin access", &security.UserAccess{Role: []string{"admin"}}, &MockEntity{}, true},
		{"Owner access", &security.UserAccess{CredentialID: "owner-id"}, &MockEntity{OwnerID: "owner-id"}, true},
		{"Access denied", &security.UserAccess{}, &MockEntity{OwnerID: "other-id"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.userAccess.CanDelete(tt.entity)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewUserAccess(t *testing.T) {
	token := &jwt.Token{ID: "test-id"}
	p := security.NewUserAccess(token)
	assert.Equal(t, "test-id", p.CredentialID)
}

func TestCanRead_WithRule(t *testing.T) {
	tests := []struct {
		name       string
		userAccess *security.UserAccess
		entity     entities.Entity
		rules      []security.Rule
		expected   bool
	}{
		{"Admin access", &security.UserAccess{Role: []string{"admin"}}, &MockEntity{}, nil, true},
		{"Public entity", &security.UserAccess{}, &MockEntity{Public: true}, nil, true},
		{"Owner access", &security.UserAccess{CredentialID: "owner-id"}, &MockEntity{OwnerID: "owner-id"}, nil, true},
		{"Access denied", &security.UserAccess{}, &MockEntity{OwnerID: "other-id"}, nil, false},
		{"Custom rule granted", &security.UserAccess{Role: []string{"user"}}, &MockEntity{OwnerID: "owner-xyz"}, []security.Rule{CustomRule}, true},
		{"Custom rule denied", &security.UserAccess{Role: []string{"user"}}, &MockEntity{OwnerID: "owner-abc"}, []security.Rule{CustomRule}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.userAccess.CanRead(tt.entity, tt.rules...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCanCreate_WithRule(t *testing.T) {
	tests := []struct {
		name       string
		userAccess *security.UserAccess
		entity     entities.Entity
		rules      []security.Rule
		expected   bool
	}{
		{"Admin access", &security.UserAccess{Role: []string{"admin"}}, &MockEntity{}, nil, true},
		{"Owner access", &security.UserAccess{CredentialID: "owner-id"}, &MockEntity{OwnerID: "owner-id"}, nil, true},
		{"Access denied", &security.UserAccess{}, &MockEntity{OwnerID: "other-id"}, nil, false},
		{"Custom rule granted", &security.UserAccess{Role: []string{"special_role"}}, &MockEntity{OwnerID: "owner-xyz"}, []security.Rule{CustomRule}, true},
		{"Custom rule denied", &security.UserAccess{Role: []string{"user"}}, &MockEntity{OwnerID: "owner-abc"}, []security.Rule{CustomRule}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.userAccess.CanCreate(tt.entity, tt.rules...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCanUpdate_WithRule(t *testing.T) {
	tests := []struct {
		name       string
		userAccess *security.UserAccess
		entity     entities.Entity
		rules      []security.Rule
		expected   bool
	}{
		{"Admin access", &security.UserAccess{Role: []string{"admin"}}, &MockEntity{}, nil, true},
		{"Owner access", &security.UserAccess{CredentialID: "owner-id"}, &MockEntity{OwnerID: "owner-id"}, nil, true},
		{"Access denied", &security.UserAccess{}, &MockEntity{OwnerID: "other-id"}, nil, false},
		{"Custom rule granted", &security.UserAccess{Role: []string{"special_role"}}, &MockEntity{OwnerID: "owner-xyz"}, []security.Rule{CustomRule}, true},
		{"Custom rule denied", &security.UserAccess{Role: []string{"user"}}, &MockEntity{OwnerID: "owner-abc"}, []security.Rule{CustomRule}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.userAccess.CanUpdate(tt.entity, tt.rules...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCanDelete_WithRule(t *testing.T) {
	tests := []struct {
		name       string
		userAccess *security.UserAccess
		entity     entities.Entity
		rules      []security.Rule
		expected   bool
	}{
		{"Admin access", &security.UserAccess{Role: []string{"admin"}}, &MockEntity{}, nil, true},
		{"Owner access", &security.UserAccess{CredentialID: "owner-id"}, &MockEntity{OwnerID: "owner-id"}, nil, true},
		{"Access denied", &security.UserAccess{}, &MockEntity{OwnerID: "other-id"}, nil, false},
		{"Custom rule granted", &security.UserAccess{Role: []string{"special_role"}}, &MockEntity{OwnerID: "owner-xyz"}, []security.Rule{CustomRule}, true},
		{"Custom rule denied", &security.UserAccess{Role: []string{"user"}}, &MockEntity{OwnerID: "owner-abc"}, []security.Rule{CustomRule}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.userAccess.CanDelete(tt.entity, tt.rules...)
			assert.Equal(t, tt.expected, result)
		})
	}
}
