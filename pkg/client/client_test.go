package client

import (
	"testing"
)

func TestClient_MakeQueryString(t *testing.T) {
	type args struct {
		params map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "bar=xyz",
			args: args{
				params: map[string]string{
					"bar": "xyz",
				},
			},
			want: "?bar=xyz",
		},
		{
			name: "bar=xyz&foo=abc",
			args: args{
				params: map[string]string{
					"foo": "abc",
					"bar": "xyz",
				},
			},
			want: "?bar=xyz&foo=abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New("")
			if got := c.MakeQueryString(tt.args.params); got != tt.want {
				t.Errorf("MakeQueryString() = %v, want %v", got, tt.want)
			}
		})
	}
}
