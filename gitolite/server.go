package gitolite

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
)

// Server is a gitolite server, contacted by calling `ssh`.
type Server struct {
	Hostname string

	// sshCommand can be set to something different to mock gitolite server in tests
	sshCommand []string
}

// Info contains information about the server and repositories.
type Info struct {
	GitoliteVersion string `json:"gitolite_version"`
	GitVersion      string `json:"git_version"`
	GitoliteUser    string `json:"GL_USER"`
	User            string `json:"USER"`
	Repos           Repos  `json:"repos"`
}

// Repos is a list of repositories.
type Repos map[string]Repo

// Names returns a sorted list of repository names.
func (repos Repos) Names() []string {
	res := make([]string, 0, len(repos))
	for name := range repos {
		res = append(res, name)
	}
	sort.Strings(res)
	return res
}

// Repo contains information about a repository.
type Repo struct {
	Perms `json:"perms"`
}

// Perms encodes permissions a user has on a repo.
type Perms struct {
	Read, Write, Create bool
}

// UnmarshalJSON decodes permissions from a JSON document.
func (p *Perms) UnmarshalJSON(data []byte) error {
	// decode as map[string]int
	var temp struct {
		Perms map[string]int `json:"perms"`
	}
	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}

	p.Read = temp.Perms["R"] > 0
	p.Write = temp.Perms["W"] > 0
	p.Create = temp.Perms["C"] > 0

	return nil
}

func (p Perms) String() string {
	res := ""
	if p.Read {
		res += "R "
	}
	if p.Write {
		res += "W "
	}
	if p.Create {
		res += "C "
	}
	return res
}

// run executes a gitolite command via `ssh` and returns the output.
func (s *Server) run(gitoliteCommand ...string) ([]byte, error) {
	var args = []string{"ssh", s.Hostname}

	// allow overriding the ssh command for tests
	if len(s.sshCommand) > 0 {
		args = s.sshCommand
	}

	args = append(args, gitoliteCommand...)

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	return cmd.Output()
}

// Info gets the basic information from the server.
func (s *Server) Info() (Info, error) {
	buf, err := s.run("info", "-json")
	if err != nil {
		return Info{}, err
	}

	var info Info
	err = json.Unmarshal(buf, &info)
	if err != nil {
		return Info{}, err
	}

	return info, nil
}

// Clone clones a repository from the server.
func (s *Server) Clone(name, targetdir string) error {
	repo := fmt.Sprintf("%v:%v", s.Hostname, name)
	cmd := exec.Command("git", "clone", "--quiet", repo, targetdir)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
