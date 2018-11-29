package env

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	exc "github.com/ContinuumLLC/platform-common-lib/src/exception"
)

const (
	//ErrEnvUnableToGetExeDirectory error code for unable to get executable's directory
	ErrEnvUnableToGetExeDirectory = "EnvUnableToGetExeDirectory"
)

//FactoryEnv interface to return the Env
type FactoryEnv interface {
	GetEnv() Env
}

//FactoryEnvImpl returns the concrete implementation of Factory
type FactoryEnvImpl struct {
}

//GetEnv returns Env
func (FactoryEnvImpl) GetEnv() Env {
	return envImpl{}
}

//Env interface defines the methods of Env functionality
type Env interface {
	GetExeDir() (string, error)
	GetFileReader(filePath string) (io.ReadCloser, error)
	GetCommandReader(command string, arg ...string) (io.ReadCloser, error)
	GetDirectoryFileCount(dirPathExpr string, args ...[]string) (io.ReadCloser, error)
	ExecuteBash(cmd string) (string, error)
	//getFileCount(dirPath string) int
}

//envImpl retruns the concrete
type envImpl struct{}

//GetExeDir returns the executable's absolute path.
//go run will return a different path than exepcted because it keeps the exe into user's temp folder
func (envImpl) GetExeDir() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", exc.New(ErrEnvUnableToGetExeDirectory, err)
	}
	return dir, nil
}

//GetFileReader returns a reader for the specified file
func (envImpl) GetFileReader(filePath string) (io.ReadCloser, error) {
	return os.Open(filePath)
}

//GetCommandReader returns a reader for the specified command
func (envImpl) GetCommandReader(command string, arg ...string) (io.ReadCloser, error) {
	cmd := exec.Command(command, arg...)
	outData, err := cmd.Output()
	byteReader := bytes.NewReader(outData)
	reader := ioutil.NopCloser(byteReader)
	return reader, err
}

func (envImpl) GetDirectoryFileCount(dirPathExpr string, args ...[]string) (io.ReadCloser, error) {
	var strFinalParameter []byte
	var strParameters string

	iArgumentLen := len(args)
	if iArgumentLen > 0 {
		iArrayLen := len(args[0])
		strInterfaceArr := make([]interface{}, iArgumentLen)

		for iCtr := 0; iCtr < iArrayLen; iCtr++ {
			for iSubCtr := 0; iSubCtr < iArgumentLen; iSubCtr++ {
				strInterfaceArr[iSubCtr] = args[iSubCtr][iCtr]
			}

			dPath := fmt.Sprintf(dirPathExpr, strInterfaceArr...)
			fileCount := getFileCount(dPath)

			strParameters = fmt.Sprintf("%s ", dPath)
			for iIndex := 0; iIndex < iArgumentLen; iIndex++ {
				strParameters += fmt.Sprint(strInterfaceArr[iIndex])
				strParameters += " "
			}
			strParameters += fmt.Sprintf("%d\n", fileCount)

			strFinalParameter = append(strFinalParameter, strParameters...)
		}
	} else {
		fileCount := getFileCount(dirPathExpr)
		strParameters = fmt.Sprintf("%s ", dirPathExpr)
		strParameters += fmt.Sprintf("%d\n", fileCount)

		strFinalParameter = append(strFinalParameter, strParameters...)
	}

	objReadCloser := ioutil.NopCloser(bytes.NewReader(strFinalParameter))
	return objReadCloser, nil
}

// ExecuteBash is to execute bash command and return it's output
func (envImpl) ExecuteBash(cmd string) (string, error) {
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func getFileCount(dirPath string) int {
	fdir, err := os.Open(dirPath)
	if nil != err {
		return 0
	}
	fileArr, err := fdir.Readdir(-1)
	if nil != err {
		fdir.Close()
		return 0
	}
	fdir.Close()
	return len(fileArr)
}
