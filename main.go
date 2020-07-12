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
	"strings"
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

	editor := EditorOptions{
		editorCmd:  "/usr/local/bin/emacsclient",
		editorOpts: []string{"--no-wait", "--quiet"},
	}
	mapper.Register("journal", JournalHandler{
		baseDir:       filepath.Join(home, "/Dropbox/org/journal/"),
		EditorOptions: editor,
	})
	mapper.Register("book", BookHandler{
		bookFile:      filepath.Join(home, "/Dropbox/org/books.org"),
		EditorOptions: editor,
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

// EditorOptions is options for editor
type EditorOptions struct {
	editorCmd  string
	editorOpts []string
}

func (e *EditorOptions) OpenFileWithLine(filename string, line int) error {
	args := append([]string{}, e.editorOpts...)
	args = append(args, fmt.Sprintf("+%d", line), filename)

	cmd := exec.Command(e.editorCmd, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintln(os.Stderr, string(out))
		return err
	}
	return nil
}

// ServiceHandler handles requests to service.
type ServiceHandler interface {
	Handle(path string, params url.Values) error
}

var _ ServiceHandler = JournalHandler{}

// JournalHandler handles journal service.
// format: go://journal/<name>?title=<title>
type JournalHandler struct {
	EditorOptions
	baseDir string
}

func (h JournalHandler) Handle(path string, params url.Values) error {
	filename := filepath.Join(h.baseDir, fmt.Sprintf("%s.org", path))

	var line int
	t, ok := params["title"]
	if ok && len(t) > 0 {
		l, err := findHeaderLine(filename, []byte(t[0]), orgHeaderPrefix)
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("WARNING: %v", err))
		} else {
			line = l
		}
	} else {
		fmt.Fprintln(os.Stderr, "WARNING: request parameter does not contain title")
	}

	return h.OpenFileWithLine(filename, line)
}

var _ ServiceHandler = BookHandler{}

// BookHandler handles book service.
// format: go://book/<title>
type BookHandler struct {
	EditorOptions
	bookFile string
}

func (h BookHandler) Handle(path string, params url.Values) error {
	name := []byte(strings.TrimPrefix(path, "/"))
	line, err := findHeaderLine(h.bookFile, name, filePropPrefix)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("WARNING: %v", err))
	}
	return h.OpenFileWithLine(h.bookFile, line)
}

var (
	orgHeaderPrefix = []byte("*")
	filePropPrefix  = []byte(":EXPORT_FILE_NAME:")
)

func findHeaderLine(filename string, header, prefix []byte) (int, error) {
	f, err := os.Open(filename)
	if err != nil {
		return 0, err
	}

	scanner := bufio.NewScanner(f)
	for line := 1; scanner.Scan(); line++ {
		if !bytes.HasPrefix(scanner.Bytes(), prefix) {
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
