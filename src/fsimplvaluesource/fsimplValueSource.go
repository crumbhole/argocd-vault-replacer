package fsimplvaluesource

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

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

// FsimplURL returns the value if SECRET_URL_PREFIX or ARGOCD_ENV_SECRET_URL_PREFIX are set
// or nil if they are not
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

// New initialises a FsimplValueSource and returns it to the user
// This is the recommended way of using this module
func New() FsimplValueSource {
	mux := fsimpl.NewMux()
	mux.Add(awssmfs.FS)
	mux.Add(awssmpfs.FS)
	mux.Add(blobfs.FS)
	mux.Add(consulfs.FS)
	mux.Add(filefs.FS)
	mux.Add(gitfs.FS)
	mux.Add(httpfs.FS)
	mux.Add(vaultfs.FS)
	return FsimplValueSource{mux: mux}
}

// FsimplValueSource is a value source getting values from bitwarden
type FsimplValueSource struct {
	mux fsimpl.FSMux
}

func ensureTrailingSlash(thing string) string {
	if strings.HasSuffix(thing, `/`) {
		return thing
	}
	return fmt.Sprintf("%s/", thing)
}

// MangleURL splits a url into a prefix and suffix and ensures
// it meets rules fsimpl requires. Easier than documenting the rules.
func MangleURL(url string) (string, string) {
	if strings.HasPrefix(url, `git`) {
		split := strings.SplitN(url, "#", 2)
		prefix := ensureTrailingSlash(split[0])
		suffix := ``
		if len(split) > 1 {
			suffix = fmt.Sprintf("#%s", split[1])
		}
		return prefix, suffix
	}
	return url, ``
}

// ManglePathForURL modifies a path for a specific protocol type
// so they can all conform to a single set of rules
func ManglePathForURL(url string, path string) string {
	if strings.HasPrefix(url, `http`) {
		return ensureTrailingSlash(path)
	}
	return path
}

// GetValue returns a value from a path+key in bitwarden or null if it doesn't exist
func (m FsimplValueSource) GetValue(path []byte, key []byte) (*[]byte, error) {
	if FsimplURL() == nil {
		return nil, errors.New("SECRET_URL_PREFIX not set")
	}
	prefix, suffix := MangleURL(*FsimplURL())
	url := fmt.Sprintf("%s%s%s", prefix, ManglePathForURL(prefix, string(path)), suffix)
	fmt.Printf("%s\n", url)
	fsys, err := m.mux.Lookup(url)
	if err != nil {
		return nil, err
	}
	value, err := fs.ReadFile(fsys, string(key))
	if err != nil {
		return nil, err
	}
	fmt.Printf("URL value: %s %s\n", url, value)
	return &value, err
}
