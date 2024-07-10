package mail_test

import (
	"testing"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	GOOD_EMAIL = "user1@example.com"
	GOOD_PASS  = "ValidP@ssw0rd1"
)

/*
func TestMail(t *testing.T) {
	config.Load(aws.String("../../../../config.test.yml"))

	tpl := template.NewTemplate("signup")
	text, html, err := tpl.Inject(template.Data{
		"AppName": "ThéTipTop",
		"Url":     "https://thetiptop.com",
	})

	assert.NoError(t, err)

	m := &mail.Mail{
		To:      []string{GOOD_EMAIL},
		Subject: "Welcome to The Tip Top",
		Text:    text,
		Html:    html,
	}

	assert.True(t, m.IsValid())

	assert.NotNil(t, m)
	msg, to, err := m.Prepare()
	assert.NotNil(t, msg)
	assert.NotNil(t, to)
	assert.NoError(t, err)
}
*/

// MockService est un mock de ServiceInterface
type MockService struct {
	mock.Mock
}

func (m *MockService) From() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockService) Expeditor() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockService) Send(mail *mail.Mail) error {
	args := m.Called(mail)
	return args.Error(0)
}

func TestMail(t *testing.T) {
	t.Run("IsValid", func(t *testing.T) {
		t.Run("ValidMail", func(t *testing.T) {
			m := &mail.Mail{
				To:      []string{GOOD_EMAIL},
				Subject: "Welcome to The Tip Top",
				Text:    []byte("hello"),
			}

			assert.True(t, m.IsValid())
		})

		t.Run("InvalidMail", func(t *testing.T) {
			m := &mail.Mail{
				To:      []string{GOOD_EMAIL},
				Subject: "Welcome to The Tip Top",
			}

			assert.False(t, m.IsValid())
		})
	})

	t.Run("Prepare", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {

			mockService := new(MockService)
			mockService.On("From").Return("from@example.com")
			mockService.On("Expeditor").Return("Expeditor Name")
			mockService.On("Send", mock.Anything).Return(nil)

			mail := &mail.Mail{
				To:          []string{"to@example.com"},
				Cc:          []string{"cc@example.com"},
				Bcc:         []string{"bcc@example.com"},
				Subject:     "Test Subject",
				Text:        []byte("This is a plain text body"),
				Html:        []byte("<p>This is an HTML body</p>"),
				Attachments: map[string][]byte{"test.txt": []byte("This is a test attachment")},
			}

			// Appel de la fonction à tester
			msgBytes, recipients, err := mail.Prepare(mockService)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			// Vérification des résultats
			if len(msgBytes) == 0 {
				t.Errorf("Expected non-empty msgBytes")
			}

			expectedRecipients := []string{"to@example.com", "cc@example.com", "bcc@example.com"}
			for i, recipient := range expectedRecipients {
				if recipients[i] != recipient {
					t.Errorf("Expected recipient %s, got %s", recipient, recipients[i])
				}
			}
		})

		t.Run("Failure", func(t *testing.T) {
			m := &mail.Mail{
				To:      []string{GOOD_EMAIL},
				Subject: "Welcome to The Tip Top",
			}

			msg, to, err := m.Prepare(nil)
			assert.Nil(t, msg)
			assert.Nil(t, to)
			assert.Error(t, err)
		})
	})

}
