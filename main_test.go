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
			args: []string{
				"go://journal/20200627?title=macOS%e3%81%a7%e7%8b%ac%e8%87%aaURLScheme%e3%81%a8%e3%83%8f%e3%83%b3%e3%83%89%e3%83%a9%e3%82%92%e5%ae%9f%e8%a3%85%e3%81%99%e3%82%8b",
			},
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
			_, err := validation(tc.args)
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
