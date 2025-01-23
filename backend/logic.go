package backend

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileInfo struct {
	Path    string    `json:"path"`
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	IsDir   bool      `json:"is_dir"`
	Mode    string    `json:"mode"`
	ModTime time.Time `json:"mod_time"`
	Hash    string    `json:"hash,omitempty"`
	Ext     string    `json:"ext"`
}

type DirInfo struct {
	ExtensionKeys []string       `json:"extension_keys"`
	Extensions    map[string]int `json:"extensions"`
	TotalFiles    int            `json:"total_files"`
}

func CalculateHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func matchesExtension(fileName string, extensions []string) bool {
	for _, ext := range extensions {
		if strings.HasSuffix(strings.ToLower(fileName), strings.ToLower(ext)) {
			return true
		}
	}
	return false
}

func CountFiles(root string) (int, error) {
	count := 0
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		count++
		return nil
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func ParseDirectory(root string, extensions []string) ([]FileInfo, error) {
	var result []FileInfo

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf(path, err)
		}

		ext := filepath.Ext(info.Name())

		item := FileInfo{
			Path:    path,
			Name:    info.Name(),
			Size:    info.Size(),
			IsDir:   info.IsDir(),
			Mode:    info.Mode().String(),
			ModTime: info.ModTime(),
			Ext:     ext,
		}

		if info.IsDir() {
			result = append(result, item)
			return nil
		}

		if !matchesExtension(info.Name(), extensions) {
			return nil
		}

		hash, err := CalculateHash(path)
		if err != nil {
			return fmt.Errorf(path, err)
		}
		item.Hash = hash

		result = append(result, item)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func WriteToJSON(data []FileInfo, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func SelectExt(reader *bufio.Reader) []string {
	fmt.Println("Input extensions:")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	extensions := []string{}
	for _, ext := range strings.Split(input, ",") {
		extensions = append(extensions, strings.TrimSpace(ext))
	}

	return extensions
}

func GetFileExt(filename string) string {
	ext := filepath.Ext(filename)
	return strings.ToLower(ext)
}

func GetStat(root string) (DirInfo, error) {
	fmt.Println("Get statistic for dir: ", root)
	var dirInfo DirInfo
	dirInfo.Extensions = make(map[string]int)

	err := filepath.Walk(root, func(root string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		ext := GetFileExt(info.Name())
		dirInfo.Extensions[ext]++
		dirInfo.TotalFiles++

		return nil
	})

	if err != nil {
		return DirInfo{}, err
	}

	dirInfo.ExtensionKeys = make([]string, 0, len(dirInfo.Extensions))
	for key := range dirInfo.Extensions {
		dirInfo.ExtensionKeys = append(dirInfo.ExtensionKeys, key)
	}

	return dirInfo, err
}

func SaveStatToJson(stat DirInfo, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	return encoder.Encode(stat)
}

func LoadStat(filename string) (DirInfo, error) {
	file, err := os.Open(filename)
	if err != nil {
		return DirInfo{}, fmt.Errorf("%v", err)
	}

	defer file.Close()

	var stat DirInfo
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&stat); err != nil {
		return DirInfo{}, fmt.Errorf("%v", err)
	}
	return stat, nil
}
