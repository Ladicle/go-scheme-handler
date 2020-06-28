package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
)

const schemeName = "go"

func main() {
	u, err := validation(os.Args[1:])
	if err != nil {
		log.Fatalln(err)
	}

	mapper := NewServiceMapper()

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}
	mapper.Register("journal", JournalHandler{
		baseDir:    filepath.Join(home, "/Dropbox/org/journal/"),
		editorCmd:  "/usr/local/bin/emacsclient",
		editorOpts: []string{"--no-wait", "--quiet"},
	})

	if err := mapper.Dispatch(u.Host, u.Path, u.Query()); err != nil {
		log.Fatalln(err)
	}
}

func validation(args []string) (*url.URL, error) {
	if len(args) == 0 {
		return nil, errors.New("URL is required arguments")
	}
	u, err := url.Parse(args[0])
	if err != nil {
		return nil, err
	}
	if u.Scheme != schemeName {
		return nil, fmt.Errorf("%q is unexpected URL scheme", u.Scheme)
	}
	return u, nil
}

func NewServiceMapper() *ServiceMapper {
	return &ServiceMapper{
		maps: make(map[string]ServiceHandler),
	}
}

// ServiceMapper maps service name and handler.
type ServiceMapper struct {
	maps map[string]ServiceHandler
}

func (m *ServiceMapper) Register(name string, handler ServiceHandler) {
	m.maps[name] = handler
}

func (m *ServiceMapper) Dispatch(service, path string, params url.Values) error {
	h, ok := m.maps[service]
	if !ok {
		return fmt.Errorf("%q is unknown service", service)
	}
	return h.Handle(path, params)
}

// ServiceHandler handles requests to service.
type ServiceHandler interface {
	Handle(path string, params url.Values) error
}

var _ ServiceHandler = JournalHandler{}

// JournalHandler handles journal service.
// format: go://journal/<name>?title=<title>
type JournalHandler struct {
	baseDir    string
	editorCmd  string
	editorOpts []string
}

func (h JournalHandler) Handle(path string, params url.Values) error {
	filename := filepath.Join(h.baseDir, fmt.Sprintf("%s.org", path))
	if _, err := os.Stat(filename); err != nil {
		return err
	}

	// TODO: search header line

	args := append([]string{}, h.editorOpts...)
	args = append(args, filename)

	cmd := exec.Command(h.editorCmd, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintln(os.Stderr, string(out))
		return err
	}
	return nil
}
