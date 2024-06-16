package upload

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func UploadImage(file *multipart.FileHeader, folder string) (string, error) {
	file.Filename = uuid.New().String() + filepath.Ext(file.Filename)
	folderName := fmt.Sprintf("uploads/%s", folder)

	src, err := file.Open()

	if err != nil {
		return "", err
	}

	defer src.Close()

	if file.Size > 5*1024*1024 {
		return "", fmt.Errorf("ukuran file melebihi 5 MB")
	}

	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		err := os.MkdirAll(folderName, 0755)
		if err != nil {
			fmt.Println("Error creating directory: ", err)
			return "", err
		}
		fmt.Println("Directory created successfully.")
	}

	validExtensions := []string{".jpg", ".jpeg", ".png", ".webp"}
	fileExtension := strings.ToLower(filepath.Ext(file.Filename))
	isValidExtension := false
	for _, ext := range validExtensions {
		if fileExtension == ext {
			isValidExtension = true
			break
		}
	}

	if !isValidExtension {
		return "", fmt.Errorf("ekstensi file tidak valid. Hanya .jpg, .jpeg, .png, .webp yang diperbolehkan")
	}
	dst, err := os.Create(filepath.Join(folderName, file.Filename))

	if err != nil {
		return "", err
	}

	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", folderName, file.Filename), nil
}
func UploadFile(file *multipart.FileHeader, folder string) (string, error) {
	file.Filename = uuid.New().String() + filepath.Ext(file.Filename)
	folderName := fmt.Sprintf("uploads/%s", folder)

	src, err := file.Open()

	if err != nil {
		return "", err
	}

	if file.Size > 5*1024*1024 {
		return "", fmt.Errorf("ukuran file melebihi 5 MB")
	}

	defer src.Close()

	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		err := os.MkdirAll(folderName, 0755)
		if err != nil {
			fmt.Println("Error creating directory: ", err)
			return "", err
		}
		fmt.Println("Directory created successfully.")
	}

	fileExtension := strings.ToLower(filepath.Ext(file.Filename))

	if fileExtension != ".pdf" {
		return "", fmt.Errorf("ekstensi file tidak valid. Hanya .pdf yang diperbolehkan")
	}

	fileName := file.Filename
	dst, err := os.Create(filepath.Join(folderName, fileName))

	if err != nil {
		return "", err
	}

	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", folderName, fileName), nil
}

func DeleteFile(filePath string) error {

	absPath, err := filepath.Abs(filePath)
	fmt.Println(absPath)
	if err != nil {
		return fmt.Errorf("gagal mendapatkan path absolut: %v", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("file tidak ditemukan: %s", absPath)
	}

	err = os.Remove(absPath)
	if err != nil {
		return fmt.Errorf("gagal menghapus file: %v", err)
	}
	return nil
}
