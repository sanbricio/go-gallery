package validators

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateNonEmptyStringField(t *testing.T) {
	tests := loadTestCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNonEmptyStringField(tt.fieldName, tt.fieldValue)

			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func loadTestCases() []struct {
	name       string
	fieldName  string
	fieldValue any
	wantErr    bool
	wantMsg    string
} {
	return []struct {
		name       string
		fieldName  string
		fieldValue any
		wantErr    bool
		wantMsg    string
	}{
		{
			name:       "Valid string",
			fieldName:  "nombre",
			fieldValue: "valor",
			wantErr:    false,
		},
		{
			name:       "Empty string",
			fieldName:  "nombre",
			fieldValue: "",
			wantErr:    true,
			wantMsg:    fmt.Sprintf(ERROR_MESSAGE, "nombre"),
		},
		{
			name:       "Valid *string",
			fieldName:  "nombre",
			fieldValue: func() *string { s := "valor"; return &s }(),
			wantErr:    false,
		},
		{
			name:       "Empty *string",
			fieldName:  "nombre",
			fieldValue: func() *string { s := ""; return &s }(),
			wantErr:    true,
			wantMsg:    fmt.Sprintf(ERROR_MESSAGE, "nombre"),
		},
		{
			name:       "Nil *string",
			fieldName:  "nombre",
			fieldValue: (*string)(nil),
			wantErr:    true,
			wantMsg:    fmt.Sprintf(ERROR_MESSAGE, "nombre"),
		},
		{
			name:       "Unsupported type",
			fieldName:  "edad",
			fieldValue: 123,
			wantErr:    true,
			wantMsg:    fmt.Sprintf(UNKNOWN_ERROR, "edad"),
		},
	}
}
