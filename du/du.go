package du

import (
	"os"
	_path "path"
)

func Size(path string) (int64, error) {
	return du(path)
}

func du(path string) (int64, error) {
	var size int64
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	fstat, err := f.Stat()
	if err != nil {
		return 0, err
	}
	if fstat.IsDir() {
		files, err := f.Readdir(-1)
		if err != nil {
			return 0, err
		}
		for _, file := range files {
			if file.IsDir() {
				s, err := du(_path.Join(path, file.Name()))
				if err != nil {
					return 0, err
				}
				size += s
			} else {
				size += file.Size()
			}
		}
	} else {
		size += fstat.Size()
	}
	return size, nil
}
