package validators

import (
	"fmt"
	"testing"
)

func TestValidateNonEmptyStringField(t *testing.T) {
	tests := loadTestCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNonEmptyStringField(tt.fieldName, tt.fieldValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expected error: %v, got: %v", tt.wantErr, err != nil)
			}
			if tt.wantErr && err.Error() != tt.wantMsg {
				t.Errorf("Expected message: '%s', got: '%s'", tt.wantMsg, err.Error())
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
			wantMsg:    fmt.Sprintf(MESSAGE_ERROR, "nombre"),
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
			wantMsg:    fmt.Sprintf(MESSAGE_ERROR, "nombre"),
		},
		{
			name:       "Nil *string",
			fieldName:  "nombre",
			fieldValue: (*string)(nil),
			wantErr:    true,
			wantMsg:    fmt.Sprintf(MESSAGE_ERROR, "nombre"),
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
