package authentication

import (
	"net/http"
	"reflect"
	"testing"
)

func TestIsAuthenticated(t *testing.T) {
	type args struct {
		keys     []string
		endpoint func(http.ResponseWriter, *http.Request)
	}
	tests := []struct {
		name string
		args args
		want http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEntitled(tt.args.keys, tt.args.endpoint); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IsAuthenticated() = %v, want %v", got, tt.want)
			}
		})
	}
}
