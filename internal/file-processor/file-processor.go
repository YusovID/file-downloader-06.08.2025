package fileprocessor

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
)

const TempDirPath = "./temp/"

var (
	ErrTaskLimitExceeded  = errors.New("task limit exceeded")
	ErrFilesLimitExceeded = errors.New("files limit exceeded")
)

type FileInfo struct {
	File     []byte
	FilePath string
	FileExt  string
}

// DownloadFile 
func DownloadFile(taskID string, URL string) error {
	fileInfo := &FileInfo{}

	err := fileInfo.getFile(URL)
	if err != nil {
		return fmt.Errorf("failed to get file: %v", err)
	}

	filePath, err := fileInfo.saveFile(taskID)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	fmt.Println(filePath)
	return nil
}

func (fi *FileInfo) getFile(URL string) error {
	// TODO вынести в другую функцию
	// validate URL
	_, err := url.ParseRequestURI(URL)
	if err != nil {
		return fmt.Errorf("invalid URL: %s", URL)
	}

	// doing request
	resp, err := http.Get(URL)
	if err != nil {
		return fmt.Errorf("can't send GET request: %v", err)
	}
	defer resp.Body.Close()

	// check request status
	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// get file extension
	contentType := resp.Header.Get("Content-Type")
	switch contentType {
	case "image/jpeg":
		fi.FileExt = ".jpeg"
	case "application/pdf":
		fi.FileExt = ".pdf"
	default:
		return fmt.Errorf("expected '.jpeg' or '.pdf', got: '%s'", contentType)
	}

	// read file body
	file, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("can't read body")
	}

	fi.File = file
	return nil
}

func (fi *FileInfo) saveFile(taskID string) (string, error) {
	// create task dir path
	taskDirPath := path.Join(TempDirPath, taskID)

	taskDir, err := os.Open(taskDirPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			mkdirErr := os.MkdirAll(taskDirPath, 0700)
			if mkdirErr != nil {
				return "", fmt.Errorf("can't make task dir: %v", mkdirErr)
			}
		} else {
			return "", fmt.Errorf("can't open task dir: %v", err)
		}
	}
	defer taskDir.Close()

	// get files count in task dir
	files, err := taskDir.ReadDir(0)
	if err != nil {
		if !errors.Is(err, os.ErrInvalid) {
			return "", fmt.Errorf("can't read task dir '%s': %v", taskDirPath, err)
		}
	}

	filesCount := len(files)
	if filesCount == 3 {
		return "", ErrFilesLimitExceeded
	}

	filePath := path.Join(taskDirPath, strconv.Itoa(filesCount)+fi.FileExt)

	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("can't create file: %v", err)
	}
	defer out.Close()

	_, err = out.Write(fi.File)
	if err != nil {
		return "", fmt.Errorf("can't write file: %v", err)
	}

	return filePath, nil
}
