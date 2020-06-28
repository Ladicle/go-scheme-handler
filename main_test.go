package main

import (
	"errors"
	"testing"
)

func TestValidation(t *testing.T) {
	testCases := []struct {
		desc    string
		args    []string
		wantErr error
	}{
		{
			desc: "valid argument",
			args: []string{"go://journal/20200728"},
		},
		{
			desc:    "no arguments",
			wantErr: errors.New("URL is required arguments"),
		},
		{
			desc:    "unknown URL scheme",
			args:    []string{"foo://journal/20200728"},
			wantErr: errors.New("\"foo\" is unexpected URL scheme"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			h := Handler{}
			err := h.validation(tc.args)
			if err != nil {
				if tc.wantErr == nil || err.Error() != tc.wantErr.Error() {
					t.Fatalf("unexpected %v error has occurred", err)
				}
				return
			}
			if tc.wantErr != nil {
				t.Fatalf("expect to occur %v error but there is no error", tc.wantErr)
			}
		})
	}
}
