package file

import (
	"blockchain-transactions/internal/logger"
	"bufio"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func CreateFile(path string) error {
	var _, err = os.Stat(path)

	if os.IsExist(err) {
		return nil
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	file.Close()

	return nil
}

func WriteFile(path, data string) error {
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	err = file.Sync()
	if err != nil {
		return err
	}

	return nil
}

func CreateFileOfBase64(fileB64, path string) error {
	dec, err := base64.StdEncoding.DecodeString(fileB64)
	if err != nil {
		return err
	}

	_, err = os.Stat(path)
	if os.IsExist(err) {
		return nil
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	if _, err := file.Write(dec); err != nil {
		return err
	}
	if err := file.Sync(); err != nil {
		return err
	}

	file.Close()
	return nil
}

func FileToB64(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	reader := bufio.NewReader(f)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	defer f.Close()
	encoded := base64.StdEncoding.EncodeToString(content)

	return encoded, nil
}

func GenerateHash(value string) string {
	hasher := md5.New()
	hasher.Write([]byte(value))
	return hex.EncodeToString(hasher.Sum(nil))
}

func MoveFile(fileName, destinyPath, originPath string) error {
	_, err := os.Stat(destinyPath)

	if os.IsNotExist(err) {
		err = os.MkdirAll(destinyPath, 0777)
		if err != nil {
			logger.Error.Println("no se pudo crear la carpeta")
			fmt.Println("No se pudo crear la carpera '" + destinyPath + " en el equipo local, contactese con el administrador.")
			return err
		}
	}

	filePathTemp := destinyPath + "/" + fileName
	err = os.Rename(originPath+fileName, filePathTemp)

	if err != nil {
		logger.Error.Println("El archivo no se ha podido mover a la carpeta de destino: %V", err)
		fmt.Println("No se ha podido mover el archivo a la carpeta de destino, contactese con el administrador.")
		return err
	}
	return nil
}

func RemoveFile(path string) error {
	e := os.Remove(path)
	return e
}

func GetMineType(data []byte) string {
	return http.DetectContentType(data)
}

func FileBytesToB64(filePDF io.Reader) string {
	reader := bufio.NewReader(filePDF)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return ""
	}

	encoded := base64.StdEncoding.EncodeToString(content)

	return encoded
}
