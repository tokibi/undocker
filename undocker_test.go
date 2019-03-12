package undocker

import (
	"testing"
)

func Test_parseReference(t *testing.T) {
	type args struct {
		arg string
	}
	tests := []struct {
		name           string
		args           args
		wantRepository string
		wantTag        string
		wantErr        bool
	}{
		{
			name: "ok",
			args: args{
				arg: "busybox:latest",
			},
			wantRepository: "busybox",
			wantTag:        "latest",
			wantErr:        false,
		},
		{
			name: "ok tag completion",
			args: args{
				arg: "busybox",
			},
			wantRepository: "busybox",
			wantTag:        "latest",
			wantErr:        false,
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
			gotRepository, gotTag, err := parseReference(tt.args.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseReference() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRepository != tt.wantRepository {
				t.Errorf("parseReference() gotRepository = %v, want %v", gotRepository, tt.wantRepository)
			}
			if gotTag != tt.wantTag {
				t.Errorf("parseReference() gotTag = %v, want %v", gotTag, tt.wantTag)
			}
		})
	}
}
