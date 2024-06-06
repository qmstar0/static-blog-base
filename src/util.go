package main

import (
	"github.com/charmbracelet/log"
	"os"
	"path/filepath"
	"sort"
	"time"
)

func TimeStamp(i int64) string {
	return time.Unix(i, 0).Format("2006-01-02")
}

func Sort(posts Posts) {
	sort.Sort(posts)
}

func CreateDir(path string) {
	dirPath := filepath.Dir(path)

	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		log.Fatalf("创建目录失败：%v", err)
	}
}

func ClearDir(dir string) {
	matches, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		log.Fatal(err)
	}
	for _, match := range matches {
		err := os.RemoveAll(match)
		if err != nil {
			log.Fatal(err)
		}
	}
}
func SetupPaginationData(paginationMap map[string]any, totalPages, currentPage int) {
	paginationMap["page"] = currentPage

	if currentPage > 1 {
		paginationMap["prev"] = currentPage - 1
	}

	if currentPage < totalPages {
		paginationMap["next"] = currentPage + 1
	}
}
