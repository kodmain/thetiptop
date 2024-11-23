package security_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kodmain/thetiptop/api/internal/application/security"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/jwt"
	"github.com/stretchr/testify/assert"
)

type MockEntityPublic struct {
	OwnerID string
	Public  bool
}

func (e *MockEntityPublic) GetOwnerID() string {
	return e.OwnerID
}

func (e *MockEntityPublic) IsPublic() bool {
	return e.Public
}

func CustomRule(p *security.UserAccess, args ...any) bool {
	return p.CredentialID != ""
}

func CustomRuleDeny(p *security.UserAccess, args ...any) bool {
	return false
}

func CustomRuleGrant(p *security.UserAccess, args ...any) bool {
	return true
}

func TestGetCredentialID(t *testing.T) {
	p := &security.UserAccess{CredentialID: "test-id"}
	assert.Equal(t, aws.String("test-id"), p.GetCredentialID())

	p = &security.UserAccess{}
	assert.Nil(t, p.GetCredentialID())
}

func TestIsAuthenticated(t *testing.T) {
	p := &security.UserAccess{CredentialID: "test-id"}
	assert.True(t, p.IsAuthenticated())

	p = &security.UserAccess{}
	assert.False(t, p.IsAuthenticated())
}

func TestIsGrantedByRoles(t *testing.T) {
	tests := []struct {
		name     string
		userRole security.Role
		roles    []security.Role
		expected bool
	}{
		{"Role granted", "admin", []security.Role{"admin"}, true},
		{"Role not granted", "user", []security.Role{"admin"}, false},
		{"Multiple roles granted", "admin", []security.Role{"user", "admin"}, true},
		{"Multiple roles not granted", "guest", []security.Role{"admin"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &security.UserAccess{Role: tt.userRole}
			result := p.IsGrantedByRoles(tt.roles...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCanRead(t *testing.T) {
	tests := []struct {
		name       string
		userAccess *security.UserAccess
		entity     database.Entity
		rule       security.Rule
		expected   bool
	}{
		{"Public entity", &security.UserAccess{}, &MockEntityPublic{Public: true}, CustomRule, true},
		{"Owner access", &security.UserAccess{CredentialID: "owner-id"}, &MockEntityPublic{OwnerID: "owner-id"}, CustomRule, true},
		{"Access denied", &security.UserAccess{}, &MockEntityPublic{OwnerID: ""}, CustomRule, false},
		{"Custom rule granted", &security.UserAccess{}, &MockEntityPublic{}, CustomRuleGrant, true},
		{"Custom rule denied", &security.UserAccess{}, &MockEntityPublic{}, CustomRuleDeny, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.userAccess.CanRead(tt.entity, tt.rule)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCanCreate(t *testing.T) {
	tests := []struct {
		name       string
		userAccess *security.UserAccess
		entity     database.Entity
		rule       security.Rule
		expected   bool
	}{
		{"Owner access", &security.UserAccess{CredentialID: "owner-id"}, &MockEntityPublic{OwnerID: "owner-id"}, CustomRule, true},
		{"Access denied", &security.UserAccess{}, &MockEntityPublic{OwnerID: "other-id"}, CustomRule, false},
		{"Custom rule granted", &security.UserAccess{}, &MockEntityPublic{}, CustomRuleGrant, true},
		{"Custom rule denied", &security.UserAccess{}, &MockEntityPublic{}, CustomRuleDeny, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.userAccess.CanCreate(tt.entity, tt.rule)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCanUpdate(t *testing.T) {
	tests := []struct {
		name       string
		userAccess *security.UserAccess
		entity     database.Entity
		rule       security.Rule
		expected   bool
	}{
		{"Owner access", &security.UserAccess{CredentialID: "owner-id"}, &MockEntityPublic{OwnerID: "owner-id"}, CustomRule, true},
		{"Access denied", &security.UserAccess{}, &MockEntityPublic{OwnerID: "other-id"}, CustomRule, false},
		{"Custom rule granted", &security.UserAccess{}, &MockEntityPublic{}, CustomRuleGrant, true},
		{"Custom rule denied", &security.UserAccess{}, &MockEntityPublic{}, CustomRuleDeny, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.userAccess.CanUpdate(tt.entity, tt.rule)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCanDelete(t *testing.T) {
	tests := []struct {
		name       string
		userAccess *security.UserAccess
		entity     database.Entity
		rule       security.Rule
		expected   bool
	}{
		{"Owner access", &security.UserAccess{CredentialID: "owner-id"}, &MockEntityPublic{OwnerID: "owner-id"}, CustomRule, true},
		{"Access denied", &security.UserAccess{}, &MockEntityPublic{OwnerID: "other-id"}, CustomRule, false},
		{"Custom rule granted", &security.UserAccess{CredentialID: "owner-xyz"}, &MockEntityPublic{OwnerID: "owner-xyz"}, CustomRule, true},
		{"Custom rule denied", &security.UserAccess{}, &MockEntityPublic{OwnerID: "owner-abc"}, CustomRule, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.userAccess.CanDelete(tt.entity, tt.rule)
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
	assert.Equal(t, security.ROLE_ADMIN, p.Role)
}

func TestNewUserAccess_NoRole(t *testing.T) {
	token := &jwt.Token{ID: "test-id"}
	p := security.NewUserAccess(token)
	assert.Equal(t, "test-id", p.CredentialID)
	assert.Equal(t, security.ROLE_ANONYMOUS, p.Role)
}

func TestNewUserAccess_NoToken(t *testing.T) {
	p := security.NewUserAccess(nil)
	assert.Equal(t, "", p.CredentialID)
	assert.Equal(t, security.ROLE_ANONYMOUS, p.Role)
}

func TestNewUserAccess_InvalidToken(t *testing.T) {
	p := security.NewUserAccess("invalid")
	assert.Equal(t, "", p.CredentialID)
	assert.Equal(t, security.ROLE_ANONYMOUS, p.Role)
}

func TestNewUserAccess_InvalidRole(t *testing.T) {
	token := &jwt.Token{
		ID:   "test-id",
		Data: map[string]interface{}{"role": 123},
	}
	p := security.NewUserAccess(token)
	assert.Equal(t, "test-id", p.CredentialID)
	assert.Equal(t, security.ROLE_ANONYMOUS, p.Role)
}

func TestNewUserAccess_IsGrantedByRules(t *testing.T) {
	p := &security.UserAccess{CredentialID: "test-id"}
	assert.True(t, p.IsGrantedByRules(CustomRule))
}

func TestNewUserAccess_IsGrantedByRules_NoRules(t *testing.T) {
	p := &security.UserAccess{CredentialID: "test-id"}
	assert.False(t, p.IsGrantedByRules())
}

func TestNewUserAccess_IsGrantedByRules_Denied(t *testing.T) {
	p := &security.UserAccess{}
	assert.False(t, p.IsGrantedByRules(CustomRule))
}
