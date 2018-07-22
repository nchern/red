// Code generated by go-bindata.
// sources:
// assets/template.txt
// DO NOT EDIT!

package app

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

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _assetsTemplateTxt = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x84\x8f\xcd\x4a\xf4\x30\x18\x85\xf7\xb9\x8a\x43\xba\x98\xcd\x07\xd3\xcf\x9d\xa5\x94\x59\x18\x1d\x41\x19\x91\x71\x2d\xef\xd4\xb7\xd3\x60\xda\x74\x92\x74\xf0\x87\xb9\x77\x49\xa2\x82\x1b\xdd\x04\xce\x73\x7e\x92\xac\x8c\x6d\xc9\xf4\xd6\x87\xea\xfc\xac\x2c\x85\x28\xb0\x56\x37\x77\x95\x28\xb0\x58\xd5\x91\x8f\x34\x70\x93\xcc\x05\x28\x20\xf4\x8c\x1d\xef\xf5\xa8\xc7\x3d\x6c\x07\x82\xd1\x23\xc3\x73\xf0\xf0\xec\x8e\xec\x90\xd6\x26\xeb\x02\x82\xc5\x40\xcf\x0c\xc7\x87\x99\x7d\xd4\xa2\x40\x7d\xab\xb6\xeb\xcd\x45\x83\x65\x67\xed\x72\x47\xee\xaf\xd9\xaf\x76\xee\x2d\x1f\xee\xaf\x41\xe8\xac\x4b\xa5\xc3\xcc\xee\x55\x14\x68\x9a\x46\x6d\x2e\xeb\xba\xfe\x6d\x2d\xd8\xc9\x63\x22\xe7\x23\x8f\x21\x6f\x67\xd7\xb2\x88\xdf\x56\x2f\x34\x4c\x86\x63\x5c\x19\xf2\x41\xb7\x9e\xc9\xb5\xfd\xe7\x0d\xe2\x4a\x6d\xf3\x8b\x1f\x33\x17\xef\x02\x00\xa4\xd7\x6f\x2c\x2b\xfc\x2f\xff\x65\x9d\xe2\xb2\x42\xb6\x13\x0a\xec\x86\x1f\x24\xd1\x4e\xb3\x79\x92\x15\xe4\x91\xcc\xcc\xf2\xdb\x3c\x89\x7c\x9e\xc4\x47\x00\x00\x00\xff\xff\x44\x4e\x77\xe8\x9d\x01\x00\x00")

func assetsTemplateTxtBytes() ([]byte, error) {
	return bindataRead(
		_assetsTemplateTxt,
		"assets/template.txt",
	)
}

func assetsTemplateTxt() (*asset, error) {
	bytes, err := assetsTemplateTxtBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/template.txt", size: 413, mode: os.FileMode(420), modTime: time.Unix(1492038386, 0)}
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
	"assets/template.txt": assetsTemplateTxt,
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
	"assets": &bintree{nil, map[string]*bintree{
		"template.txt": &bintree{assetsTemplateTxt, map[string]*bintree{}},
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
