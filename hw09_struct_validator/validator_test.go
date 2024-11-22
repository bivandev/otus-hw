package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "Valid User",
			in: User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Name:   "John",
				Age:    25,
				Email:  "user@example.com",
				Role:   "admin",
				Phones: []string{"12345678901", "10987654321"},
			},
			expectedErr: nil,
		},
		{
			name: "Valid Token",
			in: Token{
				Header:    []byte("header"),
				Payload:   []byte("payload"),
				Signature: []byte("signature"),
			},
			expectedErr: nil,
		},
		{
			name: "Invalid User: Multiple Errors",
			in: User{
				ID:     "short-id",
				Name:   "",
				Age:    15,
				Email:  "invalid_email",
				Role:   "guest",
				Phones: []string{"123", "456"},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: fmt.Errorf("must be 36 characters long")},
				{Field: "Age", Err: fmt.Errorf("must be >= 18")},
				{Field: "Email", Err: fmt.Errorf("does not match pattern ^\\w+@\\w+\\.\\w+$")},
				{Field: "Role", Err: fmt.Errorf("must be one of admin, stuff")},
				{Field: "Phones", Err: fmt.Errorf("element 0: must be 11 characters long")},
			},
		},
		{
			name: "Valid App",
			in: App{
				Version: "1.0.0",
			},
			expectedErr: nil,
		},
		{
			name: "Invalid App: Version too short",
			in: App{
				Version: "1234",
			},
			expectedErr: ValidationErrors{
				{Field: "Version", Err: fmt.Errorf("must be 5 characters long")},
			},
		},
		{
			name: "Valid Response",
			in: Response{
				Code: 200,
				Body: "OK",
			},
			expectedErr: nil,
		},
		{
			name: "Invalid Response: Code not in allowed values",
			in: Response{
				Code: 403,
				Body: "Forbidden",
			},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: fmt.Errorf("must be one of 200, 404, 500")},
			},
		},
		{
			name: "Invalid Token: Empty Payload",
			in: Token{
				Header:    []byte("header"),
				Payload:   nil,
				Signature: []byte("signature"),
			},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d: %s", i, tt.name), func(t *testing.T) {
			t.Parallel()

			err := Validate(tt.in)

			if tt.expectedErr == nil {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				return
			}

			var validationErr ValidationErrors
			if !errors.As(err, &validationErr) {
				t.Fatalf("expected validation errors, got: %v", err)
			}

			var expectedValidationErr ValidationErrors
			if !errors.As(tt.expectedErr, &expectedValidationErr) {
				t.Fatalf("expected error is not of type ValidationErrors")
			}

			if len(validationErr) != len(expectedValidationErr) {
				t.Errorf("expected %d errors, got %d", len(expectedValidationErr), len(validationErr))
			}

			for i, ve := range validationErr {
				expectedVe := expectedValidationErr[i]
				if ve.Field != expectedVe.Field || ve.Err.Error() != expectedVe.Err.Error() {
					t.Errorf("expected error for field %s: %v, got: %v", ve.Field, expectedVe, ve)
				}
			}
		})
	}
}
