package git

import (
	"path/filepath"
	"testing"

	"github.com/ahodieb/git-helpers/git/scripts"
	"github.com/stretchr/testify/assert"
)

// (*Git).Exec must not treat stderr output as an error when the command
// itself succeeded. Git uses stderr as a status channel ("Switched to
// branch...", progress) on exit 0, and earlier code mistakenly surfaced
// that as an error.
func TestExecDoesNotErrorOnStderrWhenCommandSucceeds(t *testing.T) {
	r := scripts.Init(t)
	g := Git{WorkingDir: r.Dir}

	_, err := g.Exec("checkout", "-b", "feature")
	assert.NoError(t, err)
}

// (*Git).Exec must apply WorkingDir via -C. The earlier implementation
// built the exec.Command before prepending -C, so the flag was silently
// dropped and the command ran against the test process's cwd instead.
func TestExecRespectsWorkingDir(t *testing.T) {
	r := scripts.Init(t)
	g := Git{WorkingDir: r.Dir}

	out, err := g.Exec("rev-parse", "--show-toplevel")
	assert.NoError(t, err)

	expected, err := filepath.EvalSymlinks(r.Dir)
	assert.NoError(t, err)
	got, err := filepath.EvalSymlinks(out.TrimmedStdout())
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestListBranches(t *testing.T) {
	r := scripts.Init(t)
	r.CreateBranch(t, "b1", "b2", "b3")
	g := Git{WorkingDir: r.Dir}

	expected := []string{"main", "b1", "b2", "b3"}
	got, err := g.ListBranches()
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected, got)
}

func TestCurrentBranch(t *testing.T) {
	r := scripts.Init(t)
	r.CreateBranch(t, "b1")
	g := Git{WorkingDir: r.Dir}

	got, err := g.CurrentBranch()
	assert.NoError(t, err)
	assert.Equal(t, "main", got)

	r.Checkout(t, "b1")
	got, err = g.CurrentBranch()
	assert.NoError(t, err)
	assert.Equal(t, "b1", got)
}

func TestFindMainBranch(t *testing.T) {
	cases := []struct {
		name      string
		branches  []string
		expected  string
		shouldErr bool
	}{
		{
			name:     "repo with main branch",
			branches: []string{"main", "b1", "b2"},
			expected: "main",
		},
		{
			name:     "repo with master branch",
			branches: []string{"master", "b1", "b2"},
			expected: "master",
		},
		{
			name:     "repo with main and master branch",
			branches: []string{"main", "master", "b1", "b2"},
			expected: "main",
		},
		{
			name:      "repo with neither main and master",
			branches:  []string{"b1", "b2"},
			expected:  "",
			shouldErr: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := scripts.Init(t, scripts.WithBranch(tt.branches[0]))
			r.CreateBranch(t, tt.branches[1:]...)
			g := Git{WorkingDir: r.Dir}
			got, err := g.FindMainBranch()
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}
