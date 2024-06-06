//go:build !windows
// +build !windows

package agent

import "testing"

func TestTargetPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		dlPath, destPath, want string
	}{
		// If the destination ends in /,
		// then no part of destination is stripped
		// before joining with path
		{dlPath: "a", destPath: "a/", want: "a/a"},
		{dlPath: "a/b", destPath: "c/a/", want: "c/a/a/b"},

		// if the last part of destination
		// does not match the first part of path,
		// then just join path to destination
		{dlPath: "a", destPath: ".", want: "a"},
		{dlPath: "a/b/c", destPath: ".", want: "a/b/c"},
		{dlPath: "a/b", destPath: "a/b", want: "a/b/a/b"},

		// If the last part of destination
		// matches the first part of path,
		// then remove the last part of destination
		// before joining with path
		{dlPath: "a", destPath: "a", want: "a"},
		{dlPath: "a/b", destPath: "c/a", want: "c/a/b"},
		{dlPath: "lambda.zip", destPath: "lambda.zip", want: "lambda.zip"},

		// Gotcha: this is not what you want.
		{dlPath: "a/lambda.zip", destPath: "a/lambda.zip", want: "a/lambda.zip/a/lambda.zip"},

		// Test absolute paths
		// no match, no trailing
		{dlPath: "app/a.log", destPath: "/var/logs", want: "/var/logs/app/a.log"},
		// match, no trailing
		{dlPath: "app/a.log", destPath: "/var/logs/app", want: "/var/logs/app/a.log"},
		// match, trailing
		{dlPath: "app/a.log", destPath: "/var/logs/app/", want: "/var/logs/app/app/a.log"},

		// artifact_download documentation examples
		{dlPath: "app/logs/a.log", destPath: "foo/app/", want: "foo/app/app/logs/a.log"},
		{dlPath: "app/logs/a.log", destPath: "foo/app", want: "foo/app/logs/a.log"},
		{dlPath: "app/logs/a.log", destPath: ".", want: "app/logs/a.log"},
	}

	for _, test := range tests {
		got := targetPath(test.dlPath, test.destPath)
		if got != test.want {
			t.Errorf("targetPath(%q, %q) = %q, want %q", test.dlPath, test.destPath, got, test.want)
		}
	}
}
