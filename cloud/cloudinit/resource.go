package cloudinit

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type Resource interface {
	Content() ([]byte, error)
	yaml.Marshaler
}

type StringResource struct {
	Resource
	Value string
}

func (r StringResource) Content() ([]byte, error) {
	return []byte(r.Value), nil
}

func (r StringResource) MarshalYAML() (interface{}, error) {
	return r.Value, nil
}

type FileResource struct {
	Resource
	Filename string
}

func (r FileResource) Content() ([]byte, error) {
	bs, err := ioutil.ReadFile(r.Filename)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func (r FileResource) MarshalYAML() (interface{}, error) {
	bs, err := r.Content()
	if err != nil {
		return nil, err
	}
	return string(bs), nil
}

type URLResource struct {
	Resource
	URL string
}

func (r URLResource) Content() ([]byte, error) {
	log.Printf("resolving URLResource from %s", r.URL)
	bs, err := get(r.URL)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func (r URLResource) MarshalYAML() (interface{}, error) {
	bs, err := r.Content()
	if err != nil {
		return nil, err
	}
	return string(bs), nil
}

func NewStringResource(str string) Resource {
	return StringResource{Value: str}
}

func NewFileResource(str string) Resource {
	return FileResource{Filename: str}
}

func NewURLResource(url string) Resource {
	return URLResource{URL: url}
}

func get(url string) ([]byte, error) {
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return bs, nil
}
