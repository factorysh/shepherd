package du

import (
	"os"
	_path "path"
)

// Size returns size and inode
func Size(path string) (int64, int64, error) {
	return du(path)
}

func du(path string) (int64, int64, error) {
	var size, inodes int64

	f, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()
	fstat, err := f.Stat()
	if err != nil {
		return 0, 0, err
	}
	if fstat.IsDir() {
		files, err := f.Readdir(-1)
		if err != nil {
			return 0, 0, err
		}
		for _, file := range files {
			inodes++
			if file.IsDir() {
				s, i, err := du(_path.Join(path, file.Name()))
				if err != nil {
					return 0, 0, err
				}
				size += s
				inodes += i
			} else {
				size += file.Size()
			}
		}
	} else {
		size += fstat.Size()
	}
	return size, inodes, nil
}
