package template

import (
	"bufio"
	"errors"
	"github.com/valyala/fasttemplate"
	"io/ioutil"
	"log"
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
	path, err := gitClone(sourceData)
	if err != nil {
		return err
	}

	m := map[string]interface{}{
		"SERVICE_NAME": "Jalal test",
	}

	r := makeRenderer(m)
	err = walker(path, r)
	//err = walkFiles(path, m)
	if err != nil {
		return err
	}

	_ = os.Rename(path, "./build/rendered")

	return nil
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

func renderTemplate(path string, data map[string]interface{}) (err error) {
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

func walker(path string, renderer renderer) error {
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			ps := strings.Split(path, string(os.PathSeparator))
			if isExcluded(ps[0], true) {
				return nil
			}

			if isExcluded(info.Name(), false) {
				return nil
			}

			log.Println("printing", path)
			err = renderer(path)
			return err
		})

	return err
}
func walkFiles(path string, data map[string]interface{}) error {
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			ps := strings.Split(path, string(os.PathSeparator))
			if isExcluded(ps[0], true) {
				return nil
			}

			if isExcluded(info.Name(), false) {
				return nil
			}

			log.Println("printing", path)
			err = renderTemplate(path, data)
			if err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		return err
	}
	return nil

}

func isExcluded(s string, isDir bool) bool {
	if isDir {
		for _, v := range excludeDir {
			if v == s {
				return true
			}
		}
	} else {
		for _, v := range excludeFileExt {
			if strings.HasSuffix(s, v) {
				return true
			}
		}
	}
	return false
}
