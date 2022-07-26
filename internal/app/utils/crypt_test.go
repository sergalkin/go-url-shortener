package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	type args struct {
		userID string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Encode can generate sha string using provided userID",
			args: args{"1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encode(tt.args.userID)
			assert.NotEmpty(t, got)
			assert.Empty(t, err)
		})
	}
}

func TestDecode(t *testing.T) {
	type args struct {
		userID string
	}
	tests := []struct {
		args    args
		wantErr assert.ErrorAssertionFunc
		name    string
	}{
		{
			name: "Decode can retrieve userID from sha string",
			args: args{"1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var uuid string
			sha, err := Encode(tt.args.userID)

			decodeErr := Decode(sha, &uuid)

			assert.NotEmpty(t, uuid)
			assert.Empty(t, err)
			assert.Empty(t, decodeErr)
		})
	}
}
