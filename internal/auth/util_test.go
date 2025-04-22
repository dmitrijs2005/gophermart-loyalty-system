package auth

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	type args struct {
		password string
		salt     string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Ok1", args{password: "MySuperSecretPassword", salt: "HHxqo4TGJzrcIBcv4Z1fAw=="}, "zHFUFI4QtHWwLIJ1jpYEzQ65fhSk0/Yp+yyAY/DkE2I="},
		{"Ok2", args{password: "MyOtherSecretPassword", salt: "DLRmyu5WrvxTvOpvRG89CQ=="}, "ybgSrD3nnHaYnthQBosGa7aL5F9Uy14kqlwC/E4MY5E="},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			saltBytes, err := base64.StdEncoding.DecodeString(tt.args.salt)
			assert.Nil(t, err)
			if got := HashPassword(tt.args.password, saltBytes); got != tt.want {
				t.Errorf("HashPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
