package main

import "testing"

func Test_parseReference(t *testing.T) {
	type args struct {
		arg string
	}
	tests := []struct {
		name     string
		args     args
		wantRepo string
		wantTag  string
		wantErr  bool
	}{
		{
			name: "ok",
			args: args{
				arg: "busybox:latest",
			},
			wantRepo: "busybox",
			wantTag:  "latest",
			wantErr:  false,
		},
		{
			name: "ok tag completion",
			args: args{
				arg: "busybox",
			},
			wantRepo: "busybox",
			wantTag:  "latest",
			wantErr:  false,
		},
		{
			name: "ng",
			args: args{
				arg: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRepo, gotTag, err := parseReference(tt.args.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseReference() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRepo != tt.wantRepo {
				t.Errorf("parseReference() gotRepo = %v, want %v", gotRepo, tt.wantRepo)
			}
			if gotTag != tt.wantTag {
				t.Errorf("parseReference() gotTag = %v, want %v", gotTag, tt.wantTag)
			}
		})
	}
}
