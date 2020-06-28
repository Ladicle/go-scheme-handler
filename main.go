package main

import (
	"bufio"
	"bytes"
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

	var line int
	t, ok := params["title"]
	if ok && len(t) > 0 {
		l, err := findHeaderLine(filename, []byte(t[0]))
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("WARNING: %v", err))
		} else {
			line = l
		}
	} else {
		fmt.Fprintln(os.Stderr, "WARNING: request parameter does not contain title")
	}

	args := append([]string{}, h.editorOpts...)
	args = append(args, fmt.Sprintf("+%d", line), filename)

	cmd := exec.Command(h.editorCmd, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintln(os.Stderr, string(out))
		return err
	}
	return nil
}

var orgHeaderPrefix = []byte("*")

func findHeaderLine(filename string, header []byte) (int, error) {
	f, err := os.Open(filename)
	if err != nil {
		return 0, err
	}

	scanner := bufio.NewScanner(f)
	for line := 1; scanner.Scan(); line++ {
		if !bytes.HasPrefix(scanner.Bytes(), orgHeaderPrefix) {
			continue
		}
		if bytes.Contains(scanner.Bytes(), header) {
			return line, nil
		}
	}
	return 0, fmt.Errorf("failed to find %q header", header)
}

func isOrgHeader() bool {
	return false
}
