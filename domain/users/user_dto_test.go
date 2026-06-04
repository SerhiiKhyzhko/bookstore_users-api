package users

import (
	"strings"
	"testing"

	"github.com/SerhiiKhyzhko/bookstore_users-api/v2/user_errors"
	"github.com/stretchr/testify/assert"
)

func TestUser_Validate(t *testing.T) {
	tooLongStr := strings.Repeat("a", 46)

	testCases := []struct {
		name         string
		inputUser    User
		expectedErr  string 
		expectedUser User   
	}{
		//Positive cases
		{
			name:         "Valid user",
			inputUser:    User{FirstName: "John", LastName: "Doe", Email: "john@email.com", Password: "password123"},
			expectedErr:  "",
			expectedUser: User{FirstName: "John", LastName: "Doe", Email: "john@email.com", Password: "password123"},
		},
		{
			name:         "Valid user with spaces and uppercase email",
			inputUser:    User{FirstName: "  John  ", LastName: "\tDoe \n", Email: " JOHN@EMAIL.COM ", Password: "pass"},
			expectedErr:  "",
			expectedUser: User{FirstName: "John", LastName: "Doe", Email: "john@email.com", Password: "pass"},
		},

		//First Name
		{
			name:        "Empty first name",
			inputUser:   User{FirstName: "", LastName: "Doe", Email: "john@email.com", Password: "password"},
			expectedErr: "first name is required",
		},
		{
			name:        "First name with only spaces",
			inputUser:   User{FirstName: "   ", LastName: "Doe", Email: "john@email.com", Password: "password"},
			expectedErr: "first name is required", // Після TrimSpace довжина стане < 1
		},
		{
			name:        "First name too long",
			inputUser:   User{FirstName: tooLongStr, LastName: "Doe", Email: "john@email.com", Password: "password"},
			expectedErr: "first name is too long",
		},

		//Last Name
		{
			name:        "Empty last name",
			inputUser:   User{FirstName: "John", LastName: "", Email: "john@email.com", Password: "password"},
			expectedErr: "last name is required",
		},
		{
			name:        "Last name too long",
			inputUser:   User{FirstName: "John", LastName: tooLongStr, Email: "john@email.com", Password: "password"},
			expectedErr: "last name is too long",
		},

		//Email
		{
			name:        "Empty email",
			inputUser:   User{FirstName: "John", LastName: "Doe", Email: "", Password: "password"},
			expectedErr: "email is required",
		},
		{
			name:        "Email too long",
			inputUser:   User{FirstName: "John", LastName: "Doe", Email: tooLongStr, Password: "password"},
			expectedErr: "email is too long",
		},

		//Password
		{
			name:        "Password too short",
			inputUser:   User{FirstName: "John", LastName: "Doe", Email: "john@email.com", Password: "123"},
			expectedErr: "password must be at least 4 characters",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user := tc.inputUser

			err := user.Validate()
			
			//Positive block
			if tc.expectedErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedUser, user)
			//negative block
			} else {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectedErr)
				assert.ErrorIs(t, err, user_errors.BadRequestErr)
			}
		})
	}
}
