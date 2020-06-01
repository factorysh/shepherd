package janitor

import (
	"errors"
	"fmt"
	_path "path"
	"sort"
	"time"
)

type Later struct {
	paths []string
	cfg   map[string]time.Duration
}

func NewLater(cfg map[string]time.Duration) *Later {
	paths := make([]string, len(cfg))
	i := 0
	for k := range cfg {
		paths[i] = k
		i++
	}
	sort.Sort(sort.Reverse(sort.StringSlice(paths)))
	return &Later{
		paths: paths,
		cfg:   cfg,
	}
}

func (l *Later) Get(path string) (time.Duration, error) {
	for _, p := range l.paths {
		fmt.Println(p, path)
		m, err := _path.Match(p, path)
		if err != nil {
			return 0, err
		}
		if m {
			return l.cfg[p], nil
		}
	}
	return 0, errors.New("Not found")
}
