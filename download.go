package insco

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func DownloadFile(url, dir string) (filePath string, err error) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	downloadedFilePath := filepath.Join(dir, fileName)
	fmt.Println("Downloading", url, "to", downloadedFilePath)

	file, err := os.Create(downloadedFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, err)
		return "", err
	}
	defer response.Body.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, err)
		return "", err
	}

	fmt.Println("Downloaded:", fileName)
	return downloadedFilePath, nil
}
