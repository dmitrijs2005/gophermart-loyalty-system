package config

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestParseFlags(t *testing.T) {

	// Test cases
	tests := []struct {
		name     string
		args     []string
		expected *Config
		wantErr  bool
	}{
		{"Test1 iP:port", []string{"cmd", "-a=:8080", "-d", "uri", "-r", ":9001", "-k", "secretkey", "-v", "1m"},
			&Config{":8080", "uri", ":9001", "secretkey", 1 * time.Minute}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

			os.Args = tt.args

			config := &Config{}
			parseFlags(config)

			if diff := cmp.Diff(config, tt.expected); diff != "" {
				t.Errorf("Structs mismatch (-config +expected):\n%s", diff)
			}
		})
	}
}
