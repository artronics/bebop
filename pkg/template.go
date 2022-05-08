package pkg

import (
	"bufio"
	"github.com/nhsdigital/bebop-cli/internal"
	"github.com/valyala/fasttemplate"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type SourceData struct {
	Url           string
	OutputDir     string
	TemplateData  map[string]interface{}
	ExcludedDirs  []string
	ExcludedFiles []string
}

type renderer func(path string) error

func Template(sourceData SourceData) (err error) {
	tmpDir, err := gitClone(sourceData)
	if err != nil {
		return err
	}

	r := makeRenderer(sourceData.TemplateData)
	err = walkFiles(tmpDir, r, sourceData.ExcludedDirs, sourceData.ExcludedFiles)
	if err != nil {
		return err
	}

	err = os.Rename(tmpDir, sourceData.OutputDir)

	return err
}

func gitClone(data SourceData) (string, error) {
	tempDir, err := ioutil.TempDir("", "template-*")
	if err != nil {
		return "", err
	}

	args := []string{"git", "clone", data.Url, tempDir}
	err = internal.ExecBlocking("git", args)

	return tempDir, err
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

func walkFiles(tmpDir string, renderer renderer, excludeDir []string, excludeFileExt []string) error {
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
