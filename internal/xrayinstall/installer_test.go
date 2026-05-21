// SPDX-License-Identifier: AGPL-3.0-or-later
package xrayinstall

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/amirrezakm/zeroone/internal/stack"
)

func TestValidateVersionToken(t *testing.T) {
	good := []string{"v25.1.30", "v25.2.0-rc1", "v1.0.0+build.2", "1.0.0", "25.2"}
	for _, v := range good {
		if err := ValidateVersionToken(v); err != nil {
			t.Errorf("expected %q to be valid, got %v", v, err)
		}
	}
	bad := []string{
		"",
		".",
		"..",
		"../foo",
		"v1.0.0/../etc",
		`v1.0.0\..\etc`,
		"v1.0.0 garbage",
		strings.Repeat("v", 100),
	}
	for _, v := range bad {
		if err := ValidateVersionToken(v); err == nil {
			t.Errorf("expected %q to be rejected", v)
		}
	}
}

func TestCompareVersions(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"v25.1.30", "v25.1.30", 0},
		{"v25.1.30", "v25.2.0", -1},
		{"v25.2.0", "v25.1.30", 1},
		{"v1.0.0", "v2.0.0", -1},
		{"v1.0", "v1.0.0", 0},
		{"v1.0.0-rc1", "v1.0.0", 0},
		{"", "v1.0.0", -1},
		{"v1.0.0", "", 1},
		{"garbage", "garbage", 0},
	}
	for _, tc := range cases {
		if got := CompareVersions(tc.a, tc.b); got != tc.want {
			t.Errorf("CompareVersions(%q,%q)=%d want %d", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestParseDigest(t *testing.T) {
	body := []byte("SHA1=abc\nSHA2-256= 0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef\nSHA2-512= def\n")
	got, err := ParseDigest(body)
	if err != nil {
		t.Fatal(err)
	}
	want := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}

	// Empty / malformed
	if _, err := ParseDigest([]byte("not a digest")); err == nil {
		t.Fatal("expected error on no digest line")
	}
}

func TestParseVersionToken(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Xray 25.1.30 (Xray, ...)", "v25.1.30"},
		{"Xray v25.1.30 build", "v25.1.30"},
		{"", ""},
		{"no version here", ""},
	}
	for _, tc := range cases {
		if got := parseVersionToken(tc.in); got != tc.want {
			t.Errorf("parseVersionToken(%q)=%q want %q", tc.in, got, tc.want)
		}
	}
}

func TestResolvePrefersOverride(t *testing.T) {
	dir := t.TempDir()
	imageBin := filepath.Join(dir, "image-xray")
	imageAssets := filepath.Join(dir, "image-assets")
	if err := os.WriteFile(imageBin, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(imageAssets, 0o755); err != nil {
		t.Fatal(err)
	}
	root := filepath.Join(dir, "override")
	i := New(root, imageBin, imageAssets, "v1.0.0", nil)
	if err := i.EnsureDirs(); err != nil {
		t.Fatal(err)
	}
	// Empty override → image wins.
	got := i.Resolve()
	if got.Source != "image" {
		t.Fatalf("empty override should pick image, got %+v", got)
	}
	// Populate override binary + symlink.
	verDir := filepath.Join(root, "versions", "v1.1.0")
	if err := os.MkdirAll(verDir, 0o755); err != nil {
		t.Fatal(err)
	}
	overrideBin := filepath.Join(verDir, "xray")
	if err := os.WriteFile(overrideBin, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join("..", "versions", "v1.1.0", "xray"), filepath.Join(root, "bin", "xray")); err != nil {
		t.Fatal(err)
	}
	got = i.Resolve()
	if got.Source != "override" {
		t.Fatalf("populated override should win, got %+v", got)
	}
	// AssetsDir falls back to image because override has no geoip.dat.
	if got.AssetsDir != imageAssets {
		t.Fatalf("expected assets fallback to image, got %s", got.AssetsDir)
	}
	// Add geoip.dat under override → assets switches to override tree.
	if err := os.WriteFile(filepath.Join(root, "assets", "geoip.dat"), []byte{0}, 0o644); err != nil {
		t.Fatal(err)
	}
	got = i.Resolve()
	if got.AssetsDir != filepath.Join(root, "assets") {
		t.Fatalf("expected assets to switch to override, got %s", got.AssetsDir)
	}
}

func TestEffectiveSourcesPanelOverride(t *testing.T) {
	i := New(t.TempDir(), "/dev/null", "/dev/null", "", nil)
	i.EnvSources = Sources{ReleaseMirror: "https://env.example/xray"}
	i.LoadConfig = func() stack.Config {
		return stack.Config{XrayUpdate: stack.XrayUpdateConfig{ReleaseMirror: "https://panel.example/xray"}}
	}
	src := i.EffectiveSources()
	if src.ReleaseMirror != "https://panel.example/xray" {
		t.Fatalf("panel value should win, got %q", src.ReleaseMirror)
	}
	// Clear panel → env wins.
	i.LoadConfig = func() stack.Config { return stack.Config{} }
	src = i.EffectiveSources()
	if src.ReleaseMirror != "https://env.example/xray" {
		t.Fatalf("env value should win after panel cleared, got %q", src.ReleaseMirror)
	}
}

func TestStateRoundTrip(t *testing.T) {
	i := New(t.TempDir(), "/dev/null", "/dev/null", "", nil)
	if err := i.EnsureDirs(); err != nil {
		t.Fatal(err)
	}
	st, err := i.LoadState()
	if err != nil {
		t.Fatal(err)
	}
	if st.InstalledVersion != "" {
		t.Fatalf("expected empty state, got %+v", st)
	}
	st.InstalledVersion = "v25.2.0"
	st.Source = "online"
	st.InstalledAt = 1700000000
	if err := i.SaveState(st); err != nil {
		t.Fatal(err)
	}
	got, err := i.LoadState()
	if err != nil {
		t.Fatal(err)
	}
	if got != st {
		t.Fatalf("round-trip mismatch: got %+v want %+v", got, st)
	}
}

// TestCommitInstallHonorsIncludeGeo locks in the panel toggle wiring:
// when stack.xray_update.include_geo is false, commitInstall must leave
// the assets dir alone even though the release zip carried geo files.
func TestCommitInstallHonorsIncludeGeo(t *testing.T) {
	dir := t.TempDir()
	i := New(filepath.Join(dir, "ov"), filepath.Join(dir, "img"), filepath.Join(dir, "imga"), "", nil)
	if err := i.EnsureDirs(); err != nil {
		t.Fatal(err)
	}
	// Pre-existing assets that the operator curated — must survive.
	assetsDir := filepath.Join(i.Root, "assets")
	if err := os.MkdirAll(assetsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	preserved := []byte("operator-curated-geoip")
	if err := os.WriteFile(filepath.Join(assetsDir, "geoip.dat"), preserved, 0o644); err != nil {
		t.Fatal(err)
	}

	// Stage the inputs commitInstall expects: a binary plus geo files.
	stage := filepath.Join(i.Root, "tmp", "stage")
	if err := os.MkdirAll(stage, 0o755); err != nil {
		t.Fatal(err)
	}
	bin := filepath.Join(stage, "xray")
	if err := os.WriteFile(bin, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	geoIP := filepath.Join(stage, "geoip.dat")
	if err := os.WriteFile(geoIP, []byte("from-release-zip"), 0o644); err != nil {
		t.Fatal(err)
	}
	extracted := extractedAssets{BinaryPath: bin, GeoIPPath: geoIP}

	off := false
	i.LoadConfig = func() stack.Config {
		return stack.Config{XrayUpdate: stack.XrayUpdateConfig{IncludeGeo: &off}}
	}
	if err := i.commitInstall(extracted, "v1.2.3", "online", ""); err != nil {
		t.Fatalf("commitInstall: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(assetsDir, "geoip.dat"))
	if err != nil {
		t.Fatalf("read geoip.dat after toggle-off install: %v", err)
	}
	if string(got) != string(preserved) {
		t.Fatalf("geoip.dat overwritten despite include_geo=false: got %q", got)
	}
	if _, err := os.Stat(filepath.Join(i.Root, "bin", "xray")); err != nil {
		t.Fatalf("binary swap should still complete: %v", err)
	}

	// Flip include_geo on (and stage a fresh geo file again — the previous
	// run did the os.Rename so the source file is gone) → assets refresh.
	if err := os.WriteFile(geoIP, []byte("from-release-zip"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Re-stage the binary too: prior commitInstall moved it out.
	if err := os.WriteFile(bin, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	on := true
	i.LoadConfig = func() stack.Config {
		return stack.Config{XrayUpdate: stack.XrayUpdateConfig{IncludeGeo: &on}}
	}
	if err := i.commitInstall(extracted, "v1.2.4", "online", ""); err != nil {
		t.Fatalf("commitInstall (geo on): %v", err)
	}
	got, err = os.ReadFile(filepath.Join(assetsDir, "geoip.dat"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "from-release-zip" {
		t.Fatalf("geoip.dat should be refreshed when include_geo=true; got %q", got)
	}
}

func TestResetToImageWipesOverride(t *testing.T) {
	dir := t.TempDir()
	i := New(filepath.Join(dir, "ov"), filepath.Join(dir, "img"), filepath.Join(dir, "imga"), "", nil)
	if err := i.EnsureDirs(); err != nil {
		t.Fatal(err)
	}
	verDir := filepath.Join(i.Root, "versions", "v1.0.0")
	if err := os.MkdirAll(verDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(verDir, "xray"), []byte("x"), 0o755); err != nil {
		t.Fatal(err)
	}
	_ = os.Symlink(filepath.Join("..", "versions", "v1.0.0", "xray"), filepath.Join(i.Root, "bin", "xray"))
	if err := i.SaveState(State{InstalledVersion: "v1.0.0"}); err != nil {
		t.Fatal(err)
	}
	if err := i.ResetToImage(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(i.Root, "bin", "xray")); err == nil {
		t.Fatal("expected bin/xray to be gone after reset")
	}
	st, err := i.LoadState()
	if err != nil {
		t.Fatal(err)
	}
	if st.InstalledVersion != "" {
		t.Fatalf("expected state cleared after reset, got %+v", st)
	}
}
