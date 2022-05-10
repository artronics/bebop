package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func ExecBlocking(bin string, args []string) error {
	// FIXME: this can be simplified using exec.Command see pkg/aws.go for reference
	binary, err := exec.LookPath(bin)
	if err != nil {
		return fmt.Errorf("couldn't find %s executable. Make sure %s is installed", bin, bin)
	}
	atr := syscall.ProcAttr{}

	pid, err := syscall.ForkExec(binary, args, &atr)
	if err != nil {
		return err
	}

	s, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	_, err = s.Wait()
	if err != nil {
		return err
	}

	return nil
}

var NoZeroExit = errors.New("non zero exit code")

func WaitExec(cmd *exec.Cmd) error {
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if _, ok := exiterr.Sys().(syscall.WaitStatus); ok { // TODO: this should work for both all plat. Test it for win
				return NoZeroExit
			}
		} else {
			return err
		}
	}
	return nil
}

func JustRun(bin string, args []string) (string, error) {
	binary, err := exec.LookPath(bin)
	if err != nil {
		return "", fmt.Errorf("couldn't find %s executable. Make sure %s is installed", bin, bin)
	}

	cmd := exec.Command(binary, args...)
	errOut := new(strings.Builder)
	stdOut := new(strings.Builder)
	cmd.Stderr = errOut
	cmd.Stdout = stdOut

	if err := cmd.Start(); err != nil {
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		if ext, ok := err.(*exec.ExitError); ok {
			// TODO: this should work for both all plat. Test it for win
			if _, ok := ext.Sys().(syscall.WaitStatus); ok {
				return "", fmt.Errorf(errOut.String())
			}
		} else {
			return "", err
		}
	}

	return stdOut.String(), nil
}

func ParseDataFile(path string) (map[string]interface{}, error) {
	var data map[string]interface{}

	file, err := os.Open(path)
	if err != nil {
		return data, err
	}
	defer func() {
		err = file.Close()
	}()
	bb, err := ioutil.ReadAll(file)
	if err != nil {
		return data, err
	}

	ext := filepath.Ext(file.Name())
	if ext == ".json" {
		err = json.Unmarshal(bb, &data)
		if err != nil {
			return data, err
		}

	} else if ext == ".yml" || ext == ".yaml" {
		err = yaml.Unmarshal(bb, &data)
		if err != nil {
			return data, err
		}

	} else {
		return data, errors.New("this file format is not supported")
	}

	return data, err
}
