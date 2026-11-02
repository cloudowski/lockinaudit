package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cloudowski/lockinaudit/internal/result"
	"github.com/cloudowski/lockinaudit/internal/testutil"
)

// inspectFixture runs inspect in-process and returns the summary and result JSON.
func inspectFixture(t *testing.T, name string) (summary, js []byte) {
	t.Helper()
	root := testutil.RepoRoot(t)
	out := filepath.Join(t.TempDir(), "result.json")
	var so, se bytes.Buffer
	code := Run([]string{"inspect", filepath.Join(root, "testdata", "fixtures", name),
		"--include-drafts", "--as-of", "2026-10-06", "--repo-name", name,
		"--emoji", "never", "--no-color", "--no-store", "--output", out}, &so, &se)
	if code != 0 {
		t.Fatalf("inspect %s exited %d: %s", name, code, se.String())
	}
	b, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	return so.Bytes(), b
}

func TestInspectGolden(t *testing.T) {
	for _, fx := range []string{"shop-aws", "multi-cloud", "app-only", "cycles"} {
		t.Run(fx, func(t *testing.T) {
			var sum []byte
			js := testutil.RunTwice(t, func() []byte {
				var j []byte
				sum, j = inspectFixture(t, fx)
				return j
			})
			if err := result.Validate(js); err != nil {
				t.Fatalf("result does not match schema: %v", err)
			}
			if strings.Contains(string(js), "hunter2") || strings.Contains(string(sum), "hunter2") {
				t.Fatal("seeded secret leaked into output")
			}
			if n := strings.Count(string(sum), "\n"); n > 24 {
				t.Fatalf("summary has %d lines, max 24", n)
			}
			testutil.Golden(t, "inspect/"+fx+".txt", sum)
			testutil.Golden(t, "inspect/"+fx+".json", js)
		})
	}
}
