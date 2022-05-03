package template

import (
	"bufio"
	"errors"
	"github.com/valyala/fasttemplate"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

type SourceData struct {
	Url string
}

type renderer func(path string) error

var excludeDir = []string{".git", "node_modules", ".idea", ".vscode"}
var excludeFileExt = []string{".zip", ".exe", ".tar", ".tar.gz", ".jar"}

func Template(sourceData SourceData) (err error) {
	tmpDir, err := gitClone(sourceData)
	if err != nil {
		return err
	}

	m := map[string]interface{}{
		"SERVICE_NAME": "Jalal test",
	}

	r := makeRenderer(m)
	err = walkFiles(tmpDir, r)
	if err != nil {
		return err
	}

	err = os.Rename(tmpDir, "./build/rendered")

	return err
}

func gitClone(data SourceData) (string, error) {
	binary, err := exec.LookPath("git")
	if err != nil {
		return "", errors.New("couldn't find git executable. Make sure git is installed")
	}

	tempDir, err := ioutil.TempDir("", "template-*")
	if err != nil {
		return "", err
	}

	args := []string{binary, "clone", data.Url, tempDir}
	atr := syscall.ProcAttr{}

	pid, err := syscall.ForkExec(binary, args, &atr)
	if err != nil {
		return "", err
	}

	s, err := os.FindProcess(pid)
	if err != nil {
		return "", err
	}
	_, err = s.Wait()
	if err != nil {
		return "", err
	}

	return tempDir, nil
}

func makeRenderer(data map[string]interface{}) renderer {
	return func(path string) error {
		f, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		err = os.Remove(path)
		if err != nil {
			return err
		}

		template := fasttemplate.New(string(f), "{{ ", " }}")

		renderedFile, err := os.Create(path)
		if err != nil {
			return err
		}
		defer func() {
			err = renderedFile.Close()
		}()

		w := bufio.NewWriter(renderedFile)
		_, err = template.Execute(w, data)
		if err != nil {
			return err
		}
		defer func() {
			err = w.Flush()
		}()

		return err
	}
}

func walkFiles(tmpDir string, renderer renderer) error {
	isDirExcluded := func(s string) bool {
		for _, v := range excludeDir {
			if v == s {
				return true
			}
		}
		return false
	}
	isFileExcluded := func(s string) bool {
		for _, v := range excludeFileExt {
			if strings.HasSuffix(s, v) {
				return true
			}
		}
		return false
	}

	err := filepath.Walk(tmpDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			ps := strings.Split(path[len(tmpDir)+1:], string(os.PathSeparator))
			if isDirExcluded(ps[0]) {
				return nil
			}

			if isFileExcluded(info.Name()) {
				return nil
			}

			err = renderer(path)
			return err
		})

	return err
}
