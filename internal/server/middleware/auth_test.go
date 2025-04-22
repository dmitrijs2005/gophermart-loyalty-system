package middleware

import (
	"testing"
)

func TestExtractAuthToken(t *testing.T) {
	type args struct {
		header string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"ok", args{"Bearer token"}, "token", false},
		{"wrong format", args{"Some token"}, "", true},
		{"wrong format", args{"Sometoken"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractAuthToken(tt.args.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractAuthToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExtractAuthToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestAuthMiddleware(t *testing.T) {
// 	type args struct {
// 		next http.Handler
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want http.Handler
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := AuthMiddleware(tt.args.next); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("AuthMiddleware() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
