package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	// Copy the main.go and types.go files to a test directory.
	dataloaderPath, _ := filepath.Abs("./dataloader.go")
	tmpDir, err := ioutil.TempDir("", "generator")
	require.NoError(t, err)
	err = os.Mkdir(filepath.Join(tmpDir, "pkg"), 0755)
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Run dataloader to generate the dataloader files.
	copy(t, filepath.Join(".", "testdata/main.go"), filepath.Join(tmpDir, "main.go"))
	copy(t, filepath.Join(".", "testdata/pkg/types.go"), filepath.Join(tmpDir, "pkg/main.go"))
	copy(t, filepath.Join(".", "testdata/go.mod"), filepath.Join(tmpDir, "go.mod"))
	run(t, tmpDir+"/pkg", fmt.Sprintf("go run %s -type=Pet -table=pets -okey=id -provider=mysql", dataloaderPath), "GO111MODULE=auto")
	run(t, tmpDir, "go run "+filepath.Join(tmpDir, "main.go"), "GO111MODULE=auto DB_HOST=localhost DB_USER=root DB_PASSWORD=qwerty DB_NAME=pets DB_PORT=3306")
}

func copy(t *testing.T, src, dest string) {
	t.Helper()

	srcContent, err := ioutil.ReadFile(src)
	require.NoError(t, err)
	err = ioutil.WriteFile(dest, srcContent, 0644)
	require.NoError(t, err)
}

func run(t *testing.T, dir, command, env string) {
	t.Helper()

	fmt.Println(command)

	args := strings.Split(command, " ")
	envVars := strings.Split(env, " ")

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Env = append(os.Environ(), envVars...)
	require.NoError(t, cmd.Run())
}
