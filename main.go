package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/squ1d123/go-embed/cmd/assert"
)

//go:embed ansible
var f embed.FS

func main() {
	_, err := exec.LookPath("ansible-playbook")
	assert.AssertNoErr("ansible-playbook not found in path!! You must install this first", err)

	tempDir, err := os.MkdirTemp("", "")
	assert.AssertNoErr("Could not create temp dir", err)
	defer removeDir(tempDir)

	err = syncFsToTempDir(tempDir)
	assert.AssertNoErr("Cloud not sync embeded files to tmp dir:", err)

	fmt.Println("Temp dir: ", tempDir)

	cmd := exec.Command("tree", tempDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	ansiblePlay := filepath.Join(tempDir, "ansible", "site.yaml")
	cmd = exec.Command("ansible-playbook", ansiblePlay)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func syncFsToTempDir(tempDir string) error {
	return fs.WalkDir(f, "ansible", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		destPath := filepath.Join(tempDir, path)
		if d.IsDir() {
			if err := os.MkdirAll(destPath, 0766); err != nil {
				return err
			}
		} else {
			fileContents, err := f.ReadFile(path)
			if err != nil {
				return err
			}
			if err = os.WriteFile(destPath, fileContents, 0666); err != nil {
				return err
			}
		}
		return nil
	})
}

func removeDir(dir string) {
	fmt.Println("Removing dir: ", dir)
	err := os.RemoveAll(dir)
	assert.AssertNoErr("Could not remove dir : ", err)
}
