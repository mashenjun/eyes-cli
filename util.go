package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync/atomic"
	"syscall"
	"time"
	"os/exec"
	"strings"
	"strconv"
)

const minGuidLen = 12

// ---------------------------------------------------------------------------

var g_incVal uint32

var (
	TrackingLog *log.Logger
	ErrorLog    *log.Logger
)

func NewTrackingLog(logpath string)(file *os.File) {
	file, err := os.OpenFile(logpath, os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	TrackingLog = log.New(file, "",log.Flags() &^ (log.Ldate | log.Ltime))
	return file
}

func NewErrorLog(logpath string)(file *os.File) {
	file, err := os.OpenFile(logpath, os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	ErrorLog = log.New(file, "",log.Flags() &^ (log.Ldate | log.Ltime))
	return file
}

func init() {
	g_incVal = uint32(time.Now().UnixNano() / 1e6)
}

func Make(n int) (v []byte, err error) {
	if n < minGuidLen {
		return nil, syscall.EINVAL
	}

	incVal := atomic.AddUint32(&g_incVal, 1)

	v = make([]byte, n)
	v[0] = byte(incVal)
	v[1] = byte(incVal >> 8)
	_, err = io.ReadFull(rand.Reader, v[2:])
	return
}

func Gen(n int) (s string, err error) {
	v, err := Make(n)
	if err != nil {
		return
	}
	s = base64.URLEncoding.EncodeToString(v)
	return
}

func NewUuid() (s string, err error) {
	s, err = Gen(12)
	return
}

func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func TouchFile(name string) (f *os.File, err error) {
	f, err = os.Create(name)
	return
}

func RemoveFile(name string) (err error) {
	return os.Remove(name)
}
func MoveFile(oldName, newName string)(err error){
	err = os.Rename(oldName, newName)
	return
}

func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

func TotalLine(fileName string)(total int, err error){
	dir, err := os.Getwd()
	if err != nil {
		return 0, err
	}
	wc := exec.Command("wc", "-l",fileName)
	wc.Dir = dir
	out, err := wc.Output()
	if err != nil {
		return 0, err
	}
	log.Println(string(out))
	total, err =strconv.Atoi(strings.Fields(string(out))[0])
	return total, err
}
