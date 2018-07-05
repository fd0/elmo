package gitolite

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseInfo(t *testing.T) {
	var testdata = `
{
   "repos" : {
      "user/foo/bar" : {
         "perms" : {
            "R" : 1
         }
      },
      "rwrepo" : {
         "perms" : {
            "W" : 1,
            "R" : 1
         }
      },
      "user/CREATOR/[a-zA-Z0-9].*" : {
         "perms" : {
            "C" : 1
         }
      }
   },
   "gitolite_version" : "3.6.7-5.el7",
   "GL_USER" : "foo",
   "git_version" : "1.8.3.1",
   "USER" : "gitolite@server"
}
	`

	var want = Info{
		GitoliteVersion: "3.6.7-5.el7",
		GitoliteUser:    "foo",
		GitVersion:      "1.8.3.1",
		User:            "gitolite@server",
		Repos: map[string]Repo{
			"user/foo/bar": Repo{
				Perms: Perms{
					Read: true,
				},
			},
			"rwrepo": Repo{
				Perms: Perms{
					Read:  true,
					Write: true,
				},
			},
			"user/CREATOR/[a-zA-Z0-9].*": Repo{
				Perms: Perms{
					Create: true,
				},
			},
		},
	}

	var info Info
	err := json.Unmarshal([]byte(testdata), &info)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(info, want) {
		t.Error(cmp.Diff(info, want))
	}
}

func TestServerInfo(t *testing.T) {
	target, ok := os.LookupEnv("GITOLITE_SERVER")
	if !ok {
		t.Skip("target server ($GITOLITE_SERVER) not available")
	}

	s := Server{Hostname: target}

	info, err := s.Info()
	if err != nil {
		t.Fatal(err)
	}

	// make sure versions are not empty
	if info.GitoliteVersion == "" {
		t.Errorf("gitolite version field is empty")
	}

	if info.GitVersion == "" {
		t.Errorf("git version field is empty")
	}
}
