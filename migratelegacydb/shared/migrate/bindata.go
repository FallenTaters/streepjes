// Code generated for package main by go-bindata DO NOT EDIT. (@generated)
// sources:
// files/0001.sql
package migrate

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _files0001Sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x54\x4f\x8f\x9b\x3e\x10\xbd\xf3\x29\x46\x7b\x49\x90\xb2\x12\xfb\xbb\xee\x89\x6c\x9c\x15\xfa\xb1\xa6\xa5\x8e\xb4\x7b\x42\x06\xbb\xa9\x55\x82\x23\x63\xd4\xf6\xdb\x57\x60\x43\x60\xc1\x28\xea\x75\xfe\x3c\xbf\xf7\x66\x3c\x2f\x29\x0a\x09\x02\x12\xee\x63\x04\x4d\xcd\xd5\xd6\x03\x00\x10\x0c\x22\x4c\xd0\x2b\x4a\xe1\x4b\x1a\xbd\x85\xe9\x07\xfc\x8f\x3e\x76\x5d\xae\xa2\x17\x0e\x04\xbd\x13\xc0\x09\x01\x7c\x8a\x63\x13\x6f\xbb\x5d\xb9\x2b\xad\xeb\x5f\x52\x31\xd8\xc7\xc9\xfe\x53\x4e\xc9\x92\x0f\xaf\xd9\x54\x69\x52\xb4\xd1\x3f\x32\x2d\x7f\xf2\x6a\x0a\x0a\x07\x74\x0c\x4f\x31\x81\xcd\x66\x54\xc8\xa8\xe6\x5a\x5c\x38\x1c\x42\x82\x48\xf4\x86\x16\xea\xff\x0b\x82\xe0\x31\x78\x7a\x0c\x9e\x36\x9e\xff\xec\x59\xf9\x27\x1c\x7d\x3d\x21\x88\xf0\x01\xbd\x83\x60\xbf\xb3\x56\x4b\x36\x08\x4a\xb0\xb1\xa6\x0f\xf8\xcf\x9e\xe7\x4d\x9c\xbb\xf0\x4b\x7e\x8f\x77\x45\xd9\xe4\x9f\xb5\xc6\xeb\xb6\xe6\xb4\xa4\x55\x31\x73\xe8\xa6\x29\x58\x13\x62\x88\x65\x2d\x78\xd6\x3d\x9e\xe0\x9e\x6c\x1b\xdb\x75\x8c\xe6\x7a\x0a\xaa\xf9\x59\xaa\x3f\xf0\x2f\xeb\xb0\xc6\xa7\x07\xce\x7a\x63\xfb\xc0\x76\xd9\xd8\xab\x92\xac\x29\xf4\x1d\x3c\x06\xe4\x51\xd1\x7d\x0e\x5f\x95\x28\x78\x76\xa5\x8a\xe6\x52\x96\x8e\x66\x53\x74\x2e\x29\x13\x54\x4b\x55\x2f\x94\x75\x75\xc7\x24\x45\xd1\x2b\x6e\x79\x6d\x47\x94\x7c\x48\xd1\x11\xa5\x08\xbf\xa0\x6f\x37\xcd\x82\xf9\x6b\x66\x59\xf1\x73\xd3\x6c\x62\xfc\xc0\x0e\xac\x81\x16\xcc\x8d\xe2\x00\x68\xbd\x9f\x58\xff\x20\x15\xe3\xea\xe1\x0e\xeb\x73\xaa\x34\xaf\x18\x57\x6e\xef\xed\x1e\x3a\xf3\x85\xac\x34\xaf\x74\xbd\x74\x20\x3a\xeb\x1d\x7d\x1d\xc7\x95\x9f\x6f\xaa\x6a\x4d\x75\xb3\x34\x32\x73\x9b\x04\x5b\x42\x58\x9e\xe9\x58\xeb\x64\xa8\xdd\x85\x10\xcc\xdf\xcd\x5a\x06\xed\x93\x7a\xfb\x0b\xa7\x2b\x70\x9b\x9a\x11\x66\x89\x27\xb8\x9f\xc6\xd6\x44\xdc\x1d\x06\x76\xdc\x71\x7b\xde\xd9\x34\x88\x1a\xf7\x4d\x94\xce\x2f\x9e\x38\x2b\xaa\x85\xac\xec\x7e\x7c\x17\x25\x5f\x3c\x05\x7f\x03\x00\x00\xff\xff\xf7\xd0\xbf\x7c\x62\x06\x00\x00")

func files0001SqlBytes() ([]byte, error) {
	return bindataRead(
		_files0001Sql,
		"files/0001.sql",
	)
}

func files0001Sql() (*asset, error) {
	bytes, err := files0001SqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "files/0001.sql", size: 1634, mode: os.FileMode(420), modTime: time.Unix(1616337764, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"files/0001.sql": files0001Sql,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"files": {nil, map[string]*bintree{
		"0001.sql": {files0001Sql, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
