package users

import (
	"testing"

	"github.com/SerhiiKhyzhko/bookstore_users-api/v2/user_errors"
	"github.com/stretchr/testify/assert"
)

func ptr[T any](v T) *T {
	return &v
}

func TestQueryBuilder(t *testing.T) {
	rawQuery := "UPDATE users SET"
	var userID int64 = 123

	testCases := []struct {
		name          string
		userData      PartialUser
		expectedQuery string
		expectedData  []any
		expectError   bool
	}{
		{
			name: "Update single field (FirstName)",
			userData: PartialUser{
				Id:        userID,
				FirstName: ptr("John"),
			},
			expectedQuery: "UPDATE users SET first_name = ? WHERE id=?;",
			expectedData:  []any{"John", userID},
			expectError:   false,
		},
		{
			name: "Update multiple fields",
			userData: PartialUser{
				Id:       userID,
				LastName: ptr("Doe"),
				Status:   ptr("active"),
			},
			expectedQuery: "UPDATE users SET last_name = ?, status = ? WHERE id=?;",
			expectedData:  []any{"Doe", "active", userID},
			expectError:   false,
		},
		{
			name: "Update all fields",
			userData: PartialUser{
				Id:        userID,
				FirstName: ptr("John"),
				LastName:  ptr("Doe"),
				Email:     ptr("john@test.com"),
				Status:    ptr("active"),
			},
			expectedQuery: "UPDATE users SET first_name = ?, last_name = ?, email = ?, status = ? WHERE id=?;",
			expectedData:  []any{"John", "Doe", "john@test.com", "active", userID},
			expectError:   false,
		},
		{
			name: "Empty fields (should return error)",
			userData: PartialUser{
				Id: userID,
			},
			expectedQuery: "",
			expectedData:  []any{},
			expectError:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualQuery, actualData, err := queryBuilder(rawQuery, tc.userData)

			//Negative block
			if tc.expectError {
				assert.Error(t, err)
				assert.Empty(t, actualQuery)
				assert.ErrorIs(t, err, user_errors.BadRequestErr)
				//Positive block
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedQuery, actualQuery)
				assert.Equal(t, tc.expectedData, actualData)
			}
		})
	}
}
