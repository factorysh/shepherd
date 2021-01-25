package du

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDu(t *testing.T) {
	tempDir, err := ioutil.TempDir(os.TempDir(), "du-")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)
	os.Mkdir(path.Join(tempDir, "folder"), 0700)
	for i := 0; i < 10; i++ {
		err = ioutil.WriteFile(path.Join(tempDir, "folder", fmt.Sprintf("file_%d", i)), []byte("popo"), 0600)
		assert.NoError(t, err)
	}
	s, err := Size(tempDir)
	assert.NoError(t, err)
	assert.Equal(t, int64(40), s)
}
