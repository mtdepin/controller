package api

import (
	"controller/pkg/logger"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

// HasSuffix - Suffix matcher string matches suffix in a platform specific way.
// For example on windows since its case insensitive we are supposed
// to do case insensitive checks.
func HasSuffix(s string, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

// MustGetUUID - get a random UUID.
func MustGetUUID() string {
	u, err := uuid.NewRandom()
	if err != nil {
		//logger.CriticalIf(GlobalContext, err)
	}

	return u.String()
}

func GetRequestId() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%d", r.Intn(100000000000))
}

// GetRandString generates random string with length
func GetRandString(len int) string {
	return randomString(len)
}

// Utility to create random string of strlen length
func randomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	seed := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(seed)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rnd.Intn(len(chars))]
	}
	return string(result)
}

// FileExists checks to see if a file exists
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func writeJson(w http.ResponseWriter, httpStatus int, obj interface{}) (err error) {
	var bytes []byte
	bytes, err = json.Marshal(obj)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	_, err = w.Write(bytes)
	return
}

func WriteJsonQuiet(w http.ResponseWriter, httpStatus int, obj interface{}) {
	if err := writeJson(w, httpStatus, obj); err != nil {
		logger.Errorf("ctlerror writing JSON %s: %v", obj, err)
	}
}

// GetFileModTime 获取文件的最后修改时间
func GetFileModTime(path string) time.Time {
	f, err := os.Open(path)
	if err != nil {
		return time.Now()
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return time.Now()
	}
	return fi.ModTime()
}

func WriteAllGoroutineStacks(w io.Writer) error {
	// this is based on pprof.writeGoroutineStacks, and removes the 64 MB limit
	buf := make([]byte, 1<<20)
	for i := 0; ; i++ {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			buf = buf[:n]
			break
		}
		// if len(buf) >= 64<<20 {
		// 	// Filled 64 MB - stop there.
		// 	break
		// }
		buf = make([]byte, 2*len(buf))
	}
	_, err := w.Write(buf)
	return err
}
