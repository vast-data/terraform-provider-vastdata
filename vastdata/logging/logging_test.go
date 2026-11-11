// Copyright (c) HashiCorp, Inc.

package logging

import (
	"context"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type entry struct {
	level  string
	msg    string
	fields []map[string]any
}

type recorder struct {
	mu      sync.Mutex
	entries []entry
}

func (r *recorder) add(level, msg string, f []map[string]any) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries = append(r.entries, entry{level, msg, f})
}
func (r *recorder) Trace(_ context.Context, m string, f ...map[string]any) { r.add("trace", m, f) }
func (r *recorder) Debug(_ context.Context, m string, f ...map[string]any) { r.add("debug", m, f) }
func (r *recorder) Info(_ context.Context, m string, f ...map[string]any)  { r.add("info", m, f) }
func (r *recorder) Warn(_ context.Context, m string, f ...map[string]any)  { r.add("warn", m, f) }
func (r *recorder) Error(_ context.Context, m string, f ...map[string]any) { r.add("error", m, f) }

func TestDefaultIsTFLogAndIsSafeWithoutTerraformLogger(t *testing.T) {
	SetLogger(nil)
	_, isTF := FromContext(context.Background()).(TFLogger)
	assert.True(t, isTF)
	// tflog drops messages when the context carries no Terraform logger; must not panic.
	assert.NotPanics(t, func() {
		Trace(context.Background(), "t")
		Debug(context.Background(), "d")
		Info(context.Background(), "i")
		Warn(context.Background(), "w")
		Error(context.Background(), "e", map[string]any{"k": "v"})
	})
}

func TestContextLoggerWinsOverGlobal(t *testing.T) {
	global, scoped := &recorder{}, &recorder{}
	SetLogger(global)
	t.Cleanup(func() { SetLogger(nil) })

	Info(context.Background(), "to-global")
	Info(WithLogger(context.Background(), scoped), "to-scoped")

	require.Len(t, global.entries, 1)
	assert.Equal(t, "to-global", global.entries[0].msg)
	require.Len(t, scoped.entries, 1)
	assert.Equal(t, "to-scoped", scoped.entries[0].msg)
}

func TestAllLevelsAndFieldsAreForwarded(t *testing.T) {
	r := &recorder{}
	ctx := WithLogger(context.Background(), r)

	Trace(ctx, "a")
	Debug(ctx, "b", map[string]any{"x": 1})
	Info(ctx, "c")
	Warn(ctx, "d")
	Error(ctx, "e", map[string]any{"err": "boom"}, map[string]any{"y": 2})

	var levels []string
	for _, e := range r.entries {
		levels = append(levels, e.level)
	}
	assert.Equal(t, []string{"trace", "debug", "info", "warn", "error"}, levels)
	assert.Equal(t, 1, r.entries[1].fields[0]["x"])
	assert.Len(t, r.entries[4].fields, 2)
}

func TestSetLoggerNilRestoresDefault(t *testing.T) {
	SetLogger(&recorder{})
	SetLogger(nil)
	_, isTF := FromContext(context.Background()).(TFLogger)
	assert.True(t, isTF)
}

func TestNilContextFallsBackToGlobal(t *testing.T) {
	r := &recorder{}
	SetLogger(r)
	t.Cleanup(func() { SetLogger(nil) })
	//nolint:staticcheck // a nil context must not panic.
	assert.NotPanics(t, func() { FromContext(nil).Info(nil, "x") })
	assert.Len(t, r.entries, 1)
}

// Every resource log must go through this package, otherwise a host that installs its own
// Logger (Crossplane) silently loses it: tflog drops everything outside a Terraform plugin.
func TestNoDirectTflogOutsideThisPackage(t *testing.T) {
	root := ".."
	fset := token.NewFileSet()
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == "logging" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		f, perr := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		require.NoError(t, perr, path)
		for _, imp := range f.Imports {
			p, _ := strconv.Unquote(imp.Path.Value)
			assert.NotEqual(t, "github.com/hashicorp/terraform-plugin-log/tflog", p,
				"%s imports tflog directly; use vastdata/logging instead", path)
		}
		return nil
	})
	require.NoError(t, err)
}
