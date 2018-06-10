package cloudinit

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type File struct {
	Name   string
	Target string
}

type ZipFile struct {
	Files []File
}

func (z ZipFile) Bytes() ([]byte, error) {
	buf := &bytes.Buffer{}
	gw := gzip.NewWriter(buf)
	tw := tar.NewWriter(gw)
	for _, file := range z.Files {
		info, err := os.Stat(file.Name)
		if err != nil {
			return nil, err
		}
		hdr := &tar.Header{
			Name: file.Target,
			Mode: int64(info.Mode()),
			Size: info.Size(),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return nil, err
		}
		lf, err := os.Open(file.Name)
		if err != nil {
			return nil, err
		}
		if _, err := io.Copy(tw, lf); err != nil {
			return nil, err
		}
	}
	if err := tw.Close(); err != nil {
		return nil, err
	}
	if err := gw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func CopyDirs(localToRemote map[string]string) ZipFile {
	var files []File
	for local, remote := range localToRemote {
		prefix := strings.TrimSuffix(local, "/")
		if err := filepath.Walk(local, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			f := File{
				Name:   p,
				Target: path.Join(remote, strings.TrimPrefix(p, prefix+"/")),
			}
			files = append(files, f)
			return nil
		}); err != nil {
			log.Print(err)
		}
	}
	return ZipFile{Files: files}
}
