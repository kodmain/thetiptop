package user_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	"github.com/stretchr/testify/assert"
)

func TestEmployee(t *testing.T) {
	encodingTypes := []EncodingType{FormURLEncoded, JSONEncoded}
	assert.Nil(t, start(8888, 8444))
	for _, encoding := range encodingTypes {

		var encodingName string = "FormURLEncoded"
		if encoding == JSONEncoded {
			encodingName = "JSONEncoded"
		}

		request("DELETE", "http://0.0.0.0:1080/email/all", "", encoding)
		time.Sleep(1 * time.Second)

		users := []struct {
			email     string
			password  string
			statusSU  int
			statusSI  int
			statusDel int
			statusUP  int
		}{
			// mail, pass, status-signup, status-signin
			{fmt.Sprintf("employee%v", encoding) + GOOD_EMAIL, GOOD_PASS, http.StatusCreated, http.StatusOK, http.StatusNoContent, http.StatusOK},
			{fmt.Sprintf("employee%v", encoding) + GOOD_EMAIL, GOOD_PASS + "hello", http.StatusConflict, http.StatusNotFound, http.StatusMethodNotAllowed, http.StatusBadRequest},
			{fmt.Sprintf("employee%v", encoding) + WRONG_EMAIL, WRONG_PASS, http.StatusBadRequest, http.StatusBadRequest, http.StatusMethodNotAllowed, http.StatusBadRequest},
		}

		t.Run("SignUp/"+encodingName, func(t *testing.T) {
			for _, user := range users {

				values := map[string][]any{
					"email":    {user.email},
					"password": {user.password},
				}

				RegisteredEmployee, status, err := request("POST", EMPLOYEE_REGISTER, "", encoding, values)

				employee := entities.Employee{}
				json.Unmarshal(RegisteredEmployee, &employee)
				urlwithcid := fmt.Sprintf(EMPLOYEE_WITH_ID, employee.ID)

				assert.NotNil(t, employee)
				assert.Nil(t, err)
				assert.Equal(t, user.statusSU, status)

				if status == http.StatusCreated {
					t.Run("Validation/"+encodingName, func(t *testing.T) {
						email, err := getMailFor(user.email, 100)
						assert.Nil(t, err)
						assert.Equal(t, user.email, email.To[0].Address)
					})

					t.Run("Validation/recover/"+encodingName, func(t *testing.T) {
						_, status, err := request("POST", USER_VALIDATION_RENEW, "", encoding, map[string][]any{
							"email": {user.email},
							"type":  {entities.MailValidation.String()},
						})
						assert.Nil(t, err)
						assert.Equal(t, http.StatusNoContent, status)
						email, err := getMailFor(user.email, 100)
						assert.Nil(t, err)
						assert.Equal(t, user.email, email.To[0].Address)
						token := extractToken(email.HTML)
						_, status, err = request("PUT", USER_REGISTER_VALIDATION, "", encoding, map[string][]any{
							"token": {token},
							"email": {user.email},
						})
						assert.Nil(t, err)
						assert.Equal(t, http.StatusOK, status)
					})
				}

				JWT, status, err := request("POST", USER_AUTH, "", encoding, values)
				assert.Nil(t, err)
				assert.Equal(t, user.statusSI, status)

				if status == http.StatusOK {
					var tokenData fiber.Map
					err = json.Unmarshal(JWT, &tokenData)
					assert.Nil(t, err)

					t.Run("Renew/"+encodingName, func(t *testing.T) {
						refresh_token_sting := tokenData["refresh_token"].(string)
						users := []struct {
							token  string
							status int
						}{
							{"Bearer " + refresh_token_sting, http.StatusOK},
							{"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMxMDkxMzEsImlkIjoiN2M3OTQwMGYtMDA2YS00NzVlLTk3YjYtNWRiZGUzNzA3NjAxIiwib2ZmIjo3MjAwLCJ0eXBlIjoxLCJ0eiI6IkxvY2FsIn0.5Lae56HNcQ1OHcP_FhTfcOOtHpaZVgRFy6vzzBugN7Y", http.StatusUnauthorized}, // Replace with actual expired JWT token
							{"Bearer malformed.jwt.token.here", http.StatusUnauthorized},
							{"", http.StatusUnauthorized},
						}

						for _, user := range users {
							_, status, err := request("GET", USER_AUTH_RENEW, user.token, encoding)
							assert.Nil(t, err)
							assert.Equal(t, user.status, status)
						}
					})

					authorization := "Bearer " + tokenData["access_token"].(string)

					t.Run("GetByID/"+encodingName, func(t *testing.T) {
						_, status, err := request("GET", urlwithcid, authorization, encoding, nil)
						assert.Nil(t, err)
						assert.Equal(t, http.StatusOK, status)
					})

					t.Run("Password/"+encodingName, func(t *testing.T) {
						_, status, err := request("POST", USER_VALIDATION_RENEW, authorization, encoding, map[string][]any{
							"email": {user.email + "wrong"},
							"type":  {entities.PasswordRecover.String()},
						})

						assert.NoError(t, err)
						assert.Equal(t, http.StatusNotFound, status)

						_, status, err = request("POST", USER_VALIDATION_RENEW, authorization, encoding, map[string][]any{
							"email": {user.email},
							"type":  {entities.PasswordRecover.String()},
						})

						assert.Nil(t, err)
						assert.Equal(t, http.StatusNoContent, status)
						email, err := getMailFor(user.email, 100)
						assert.Nil(t, err)
						assert.Equal(t, user.email, email.To[0].Address)

						token := extractToken(email.HTML)
						assert.NotEmpty(t, token)

						_, status, err = request("PUT", USER_PASSWORD, authorization, encoding, map[string][]any{
							"token":    {token},
							"email":    {user.email},
							"password": {GOOD_PASS_UPDATED},
						})

						assert.Nil(t, err)
						assert.Equal(t, http.StatusOK, status)

						_, status, err = request("POST", USER_AUTH, authorization, encoding, values)
						assert.NoError(t, err)
						assert.Equal(t, http.StatusNotFound, status)

						values["password"] = []any{GOOD_PASS_UPDATED}
						_, status, err = request("POST", USER_AUTH, authorization, encoding, values)
						assert.NoError(t, err)
						assert.Equal(t, user.statusSI, status)
					})

					_, status, err := request("PUT", EMPLOYEE, authorization, encoding, map[string][]any{
						"id": {employee.ID},
					})

					assert.Nil(t, err)
					assert.Equal(t, user.statusUP, status)

					t.Run("Delete/"+encodingName, func(t *testing.T) {
						_, status, err := request("DELETE", urlwithcid, authorization, encoding, nil)
						assert.Nil(t, err)
						assert.Equal(t, user.statusDel, status)
					})
				}
			}
		})
	}
	assert.Nil(t, stop())
}
