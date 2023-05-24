package main

import (
	"testing"
)

func TestValidateEmpty(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Empty filename",
			input:   "",
			wantErr: true,
		},
		{
			name:    "Non-empty filename",
			input:   "testfile",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateEmpty(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("validateEmpty() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateLength(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Filename length exceeds 30",
			input:   "thisisaveryverylongfilenameexceedingthirtycharacters",
			wantErr: true,
		},
		{
			name:    "Filename length less than 30",
			input:   "testfile",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateLength(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("validateLength() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateInvalidChars(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Filename contains invalid characters",
			input:   "invalid/filename",
			wantErr: true,
		},
		{
			name:    "Filename does not contain invalid characters",
			input:   "valid_filename",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateInvalidChars(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("validateInvalidChars() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
