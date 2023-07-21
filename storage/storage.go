package storage

import (
	"path/filepath"
)

var Storage storage

type storage struct {
	// public 目录的文件夹路径
	AbsPath  string
	DiskName string
	Uri      string
}

func Init(storagePath string) {
	Storage = storage{
		AbsPath: storagePath,
	}
}

// FullPath 获取绝对路径
func (s *storage) FullPath(path string) string {
	return filepath.Join(s.AbsPath, path)
}
