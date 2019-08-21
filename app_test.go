package main

import (
	"reflect"
	"testing"

	"github.com/giorgioprevitera/redditproto"
)

func Test_getThread(t *testing.T) {
	type args struct {
		url    string
		client getter
	}

	m := &mockGetter{}
	tests := []struct {
		name    string
		args    args
		want    *redditproto.Link
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "my test",
			args: args{
				url:    "my url",
				client: m,
			},
			want:    *redditproto.Link{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getThread(tt.args.url, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("getThread() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getThread() = %v, want %v", got, tt.want)
			}
		})
	}
}
