package janitor

import (
	"errors"
	"fmt"
	_path "path"
	"sort"
	"time"
)

type Janitor struct {
	paths []string
	cfg   map[string]time.Duration
}

func New(cfg map[string]time.Duration) *Janitor {
	paths := make([]string, len(cfg))
	i := 0
	for k := range cfg {
		paths[i] = k
		i++
	}
	sort.Sort(sort.Reverse(sort.StringSlice(paths)))
	return &Janitor{
		paths: paths,
		cfg:   cfg,
	}
}

func (j *Janitor) Get(path string) (time.Duration, error) {
	for _, p := range j.paths {
		fmt.Println(p, path)
		m, err := _path.Match(p, path)
		if err != nil {
			return 0, err
		}
		if m {
			return j.cfg[p], nil
		}
	}
	return 0, errors.New("Not found")
}
