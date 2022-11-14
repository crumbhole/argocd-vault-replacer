package fsimplvaluesource

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	fsimpl "github.com/hairyhenderson/go-fsimpl"
	"github.com/hairyhenderson/go-fsimpl/awssmfs"
	"github.com/hairyhenderson/go-fsimpl/awssmpfs"
	"github.com/hairyhenderson/go-fsimpl/blobfs"
	"github.com/hairyhenderson/go-fsimpl/consulfs"
	"github.com/hairyhenderson/go-fsimpl/filefs"
	"github.com/hairyhenderson/go-fsimpl/gitfs"
	"github.com/hairyhenderson/go-fsimpl/httpfs"
	"github.com/hairyhenderson/go-fsimpl/vaultfs"
)

const (
	envCheck   = "SECRET_URL_PREFIX"
	argoPrefix = "ARGOCD_ENV_"
)

// FsimplOK returns true of SECRET_URL_PREFIX or ARGOCD_ENV_SECRET_URL_PREFIX are set
// If ARGOCD_ENV_SECRET_URL_PREFIX is set the value is copied to SECRET_URL_PREFIX
func FsimplURL() *string {
	val, got := os.LookupEnv(argoPrefix + envCheck)
	if !got {
		val, got = os.LookupEnv(envCheck)
		if !got {
			return nil
		}
	}
	return &val
}

func (m FsimplValueSource) init() {
	if m.mux == nil {
		mux := fsimpl.NewMux()
		m.mux = &mux
		m.mux.Add(awssmfs.FS)
		m.mux.Add(awssmpfs.FS)
		m.mux.Add(blobfs.FS)
		m.mux.Add(consulfs.FS)
		m.mux.Add(filefs.FS)
		m.mux.Add(gitfs.FS)
		m.mux.Add(httpfs.FS)
		m.mux.Add(vaultfs.FS)
	}
}

// FsimplValueSource is a value source getting values from bitwarden
type FsimplValueSource struct {
	mux *fsimpl.FSMux
}

// func (FsimplValueSource) getItemSplitPath(path string) (*bwwrap.BwItem, error) {
// 	pathParts := strings.Split(string(path), `/`)
// 	keyUsed := pathParts[len(pathParts)-1]
// 	pathUsed := strings.Join(pathParts[:len(pathParts)-1], `/`)
// 	return bwwrap.GetItemFromFolder(keyUsed, pathUsed)
// }

// GetValue returns a value from a path+key in bitwarden or null if it doesn't exist
func (m FsimplValueSource) GetValue(path []byte, key []byte) (*[]byte, error) {
	if FsimplURL() == nil {
		return nil, errors.New("SECRET_URL_PREFIX not set")
	}
	m.init()
	url := fmt.Sprintf("%s%s", *FsimplURL(), path)
	fsys, err := m.mux.Lookup(url)
	if err != nil {
		return nil, err
	}
	value, err := fs.ReadFile(fsys, string(key))
	return &value, err
}
