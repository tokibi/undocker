package undocker

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

func Test_isDockerHub(t *testing.T) {
	mux, serverURL, teardown := setupMockAPI()
	defer teardown()
	// Reproduce authentication failed response on DockerHub
	mux.HandleFunc("/hub/v2/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Www-Authenticate", `Bearer realm="https://auth.docker.io/token",service="registry.docker.io"`)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, `{"errors":[{"code":"UNAUTHORIZED","message":"authentication required","detail":null}]}`)
	})
	mux.HandleFunc("/other/v2/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{}`)
	})

	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "true docker hub",
			url:  serverURL + "/hub",
			want: true,
		},
		{
			name: "false other something",
			url:  serverURL + "/other",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.url)
			if err != nil {
				t.Error(err)
			}
			if got := isDockerHub(u); got != tt.want {
				t.Errorf("isDockerHub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_complementOfficialRepoName(t *testing.T) {
	type args struct {
		repository string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok complemented",
			args: args{
				repository: "busybox",
			},
			want: "library/busybox",
		},
		{
			name: "ok did not complement",
			args: args{
				repository: "library/busybox",
			},
			want: "library/busybox",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := complementOfficialRepoName(tt.args.repository); got != tt.want {
				t.Errorf("complementOfficialRepoName() = %v, want %v", got, tt.want)
			}
		})
	}
}
