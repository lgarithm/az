package cloudinit

import (
	"encoding/base64"
	"fmt"
	"path"

	yaml "gopkg.in/yaml.v2"
)

// ScriptSection is a section of http://cloudinit.readthedocs.io/en/latest/topics/modules.html#runcmd
type ScriptSection interface {
	yaml.Marshaler
}

type OneLineScript struct {
	ScriptSection

	cmd string
}

func (s OneLineScript) MarshalYAML() (interface{}, error) {
	return s.cmd, nil
}

func NewOneLineScript(cmd string) ScriptSection {
	return &OneLineScript{cmd: cmd}
}

type FileScript struct {
	ScriptSection

	file FileResource
}

func (s FileScript) MarshalYAML() (interface{}, error) {
	return s.file.MarshalYAML()
}

func NewFileScript(name string) ScriptSection {
	return &FileScript{file: FileResource{Filename: name}}
}

type UnzipFileScript struct {
	ScriptSection

	ZipFile ZipFile
	Owner   string
}

func (s UnzipFileScript) MarshalYAML() (interface{}, error) {
	src, err := s.ZipFile.Bytes()
	if err != nil {
		return nil, err
	}
	const name = "data"
	script := fmt.Sprintf(`
%s='
%s
'
cd %s && echo $%s | tr -d " " | base64 --decode | sudo -u %s tar -xz

`, name, base64endode(src), homeOf(s.Owner), name, s.Owner)
	return script, nil
}

func NewUnzipFileScript(zipFile ZipFile, owner string) UnzipFileScript {
	return UnzipFileScript{
		ZipFile: zipFile,
		Owner:   owner,
	}
}

type WriteFileScript struct {
	Content  Resource
	FilePath string
}

func (s WriteFileScript) MarshalYAML() (interface{}, error) {
	bs, err := s.Content.Content()
	if err != nil {
		return nil, err
	}
	script := fmt.Sprintf(`
data='
%s
'
echo $data | tr -d " " | base64 --decode > %s
`, base64endode(bs), s.FilePath)
	return script, nil
}

func NewWriteFileScript(r Resource, filapath string) WriteFileScript {
	return WriteFileScript{
		Content:  r,
		FilePath: filapath,
	}
}

func homeOf(user string) string {
	if user == "root" {
		return "/"
	}
	return path.Join("/home", user)
}

func base64endode(src []byte) string {
	return lineWrap(76, base64.StdEncoding.EncodeToString(src))
}

func lineWrap(width int, text string) string {
	var output string
	n := len(text)
	for i := 0; i < n; i += width {
		j := i + width
		if j > n {
			j = n
		}
		output += text[i:j] + "\n"
	}
	return output
}
