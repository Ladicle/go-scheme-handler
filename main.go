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

const (
	schemeName  = "go"
	journalPath = "/Dropbox/org/journal/"
)

func main() {
	h := Handler{}
	if err := h.validation(os.Args[1:]); err != nil {
		log.Fatalln(err)
	}
	if err := h.handle(); err != nil {
		log.Fatalln(err)
	}
}

type Handler struct {
	url *url.URL
}

func (h *Handler) validation(args []string) error {
	if len(args) == 0 {
		return errors.New("URI is required arguments")
	}
	u, err := url.Parse(args[0])
	if err != nil {
		return err
	}
	if u.Scheme != schemeName {
		return fmt.Errorf("%q is unexpected URI scheme: %+v", u.Scheme, *u)
	}
	h.url = u
	return nil
}

func (h *Handler) handle() error {
	switch h.url.Host {
	// "journal/<filename>"
	case "journal":
		return openJournalEditor(h.url.Path)
	default:
		return fmt.Errorf("%q is unknown hostname", h.url.Host)
	}
}

func openJournalEditor(path string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	filename := filepath.Join(home, journalPath, fmt.Sprintf("%s.org", path))
	if _, err := os.Stat(filename); err != nil {
		return err
	}
	cmd := exec.Command("emacsclient", "-qn", filename)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintln(os.Stderr, string(out))
		return err
	}
	return nil
}
