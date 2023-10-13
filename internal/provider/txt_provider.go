package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/callmemhz/godash/model"
)

type TxtProviderFactory struct {
	RepoPath string
}

func (factory *TxtProviderFactory) NewProvider() (Provider, error) {
	if _, err := os.Stat(factory.RepoPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("repo path %s not exsits", factory.RepoPath)
	}
	p := &TxtProvider{repo: factory.RepoPath}
	return p, nil
}

type TxtProvider struct {
	repo string
}

type TxtEntry struct {
	title    string
	subtitle string
	content  string
	typ      model.EntryType
	viewed   uint
}

func (entry *TxtEntry) Title() string         { return entry.title }
func (entry *TxtEntry) Subtitle() string      { return entry.subtitle }
func (entry *TxtEntry) Content() string       { return entry.content }
func (entry *TxtEntry) Type() model.EntryType { return entry.typ }
func (entry *TxtEntry) Viewed() uint          { return 0 }

func (p *TxtProvider) Search(namespace, keyword string) ([]model.Entry, error) {
	var entries []model.Entry
	files, err := os.ReadDir(p.repo)
	if err != nil {
		return nil, fmt.Errorf("tranversal %s error: %v", p.repo, err)
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if !strings.HasSuffix(f.Name(), ".txt") {
			continue
		}
		if !strings.HasPrefix(f.Name(), namespace) {
			continue
		}

		filename := filepath.Join(p.repo, f.Name())
		ns := f.Name()[:len(f.Name())-len(filepath.Ext(f.Name()))] // namespace
		data, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("cannot read %s: %v\n", filename, err)
			continue
		}

		recs := strings.Split(string(data), ";;;")
		for _, rec := range recs {
			if rec = strings.TrimSpace(rec); rec == "" {
				continue
			}
			lines := strings.Split(rec, "\n")
			title := strings.TrimSpace(lines[0])
			if (namespace != "" && keyword == "-") || fuzzyMatch(title, keyword) {
				content := strings.TrimSpace(strings.Join(lines[1:], "\n"))
				desc := "@" + ns
				entry := &TxtEntry{
					title:    title,
					subtitle: desc,
					content:  content,
					typ:      model.EntryTypeMemo,
					viewed:   0,
				}
				entries = append(entries, entry)
			}
		}
	}
	return entries, nil
}
