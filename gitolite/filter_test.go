package gitolite

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFilter(t *testing.T) {
	var tests = []struct {
		names   []string
		pattern string
		want    []string
	}{
		{
			names:   []string{"foo", "bar", "baz"},
			pattern: "",
			want:    []string{"foo", "bar", "baz"},
		},
		{
			names:   []string{"user/foo", "bar"},
			pattern: "user",
			want:    []string{"user/foo"},
		},
		{
			names:   []string{"user/foo", "bar"},
			pattern: "user/",
			want:    []string{"user/foo"},
		},
		{
			names:   []string{"user/foo", "bar"},
			pattern: "Foo",
			want:    []string{"user/foo"},
		},
		{
			names:   []string{"user/foo", "bar"},
			pattern: "user/*",
			want:    []string{"user/foo"},
		},
		{
			names:   []string{"user/foo/repo", "bar"},
			pattern: "user/*",
			want:    []string{"user/foo/repo"},
		},
		{
			names:   []string{"user/foo/repo", "bar"},
			pattern: "user/f*",
			want:    []string{"user/foo/repo"},
		},
		{
			names:   []string{"user/foo/repo", "bar"},
			pattern: "user/[a-f]*",
			want:    []string{"user/foo/repo"},
		},
		{
			names:   []string{"user/foo/repo", "bar"},
			pattern: "user/[x-z]*",
		},
		{
			names:   []string{"user/foo/repo", "bar"},
			pattern: "user/[xy]??",
		},
		{
			names:   []string{"user/foo/repo", "bar"},
			pattern: "user/[a-f]??",
			want:    []string{"user/foo/repo"},
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			res := filter(test.names, test.pattern)
			sort.Strings(res)
			sort.Strings(test.want)

			if !cmp.Equal(test.want, res) {
				t.Error(cmp.Diff(test.want, res))
			}
		})
	}
}
