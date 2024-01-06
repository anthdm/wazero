package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/experimental/sysfs"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	internalsys "github.com/tetratelabs/wazero/internal/sys"
)

func validateMounts(mount string) (rc int, rootPath string, config wazero.FSConfig) {
	config = wazero.NewFSConfig()
	if len(mount) == 0 {
		return 1, rootPath, config
	}

	readOnly := false
	if trimmed := strings.TrimSuffix(mount, ":ro"); trimmed != mount {
		mount = trimmed
		readOnly = true
	}

	// TODO: Support wasm paths with colon in them.
	var dir, guestPath string
	if clnIdx := strings.LastIndexByte(mount, ':'); clnIdx != -1 {
		dir, guestPath = mount[:clnIdx], mount[clnIdx+1:]
	} else {
		dir = mount
		guestPath = dir
	}

	// Eagerly validate the mounts as we know they should be on the host.
	if abs, err := filepath.Abs(dir); err != nil {
		return 1, rootPath, config
	} else {
		dir = abs
	}

	if stat, err := os.Stat(dir); err != nil {
		return 1, rootPath, config
	} else if !stat.IsDir() {
	}

	root := sysfs.DirFS(dir)
	if readOnly {
		root = &sysfs.ReadFS{FS: root}
	}

	config = config.(sysfs.FSConfig).WithSysFSMount(root, guestPath)

	if internalsys.StripPrefixesAndTrailingSlash(guestPath) == "" {
		rootPath = dir
	}
	return 0, rootPath, config
}

func main() {
	b, err := os.ReadFile("../ffaas/examples/go/app.wasm")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	runtime := wazero.NewRuntime(ctx)
	wasi_snapshot_preview1.MustInstantiate(ctx, runtime)
	mod, err := runtime.CompileModule(ctx, b)
	if err != nil {
		log.Fatal(err)
	}
	_, _, fsconfig := validateMounts("/")
	config := wazero.NewModuleConfig().
		WithFSConfig(fsconfig).
		WithStdin(os.Stdin).
		WithStdout(os.Stdout).
		WithStderr(os.Stderr)
	_, err = runtime.InstantiateModule(ctx, mod, config)
	if err != nil {
		log.Fatal(err)
	}
}
