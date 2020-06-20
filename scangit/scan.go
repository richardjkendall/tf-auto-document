package scangit

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// ScanGit scans git repos looking for information to include in tf module documentation
type ScanGit struct {
	repo *git.Repository
	tags map[string]string
}

// GitCommit stores commits
type GitCommit struct {
	Hash    string
	Tag     string
	Message string
}

// New creates a new instance of GitScanner
func New() *ScanGit {
	return &ScanGit{tags: make(map[string]string)}
}

// Open opens a local git repo
func (scanner *ScanGit) Open(path string) error {
	r, err := git.PlainOpen(path)
	if err != nil {
		return err
	}
	scanner.repo = r
	return nil
}

// GetCommits gets list of commits
func (scanner *ScanGit) GetCommits(subpath string) ([]GitCommit, error) {
	var r []GitCommit
	lopts := git.LogOptions{
		PathFilter: func(s string) bool {
			if strings.HasPrefix(s, subpath) && (strings.HasSuffix(s, "tf") || strings.HasSuffix(s, "TF")) {
				return true
			}
			return false
		},
	}
	commits, err := scanner.repo.Log(&lopts)
	if err != nil {
		return r, err
	}
	commits.ForEach(func(c *object.Commit) error {
		r = append(r, GitCommit{
			Hash:    hex.EncodeToString(c.Hash[:]),
			Message: c.Message,
			Tag:     scanner.tags[hex.EncodeToString(c.Hash[:])],
		})
		return nil
	})
	return r, nil
}

// LoadTags populates an in-memory list of tags for later use
func (scanner *ScanGit) LoadTags() error {
	tags, err := scanner.repo.Tags()
	if err != nil {
		return err
	}
	tags.ForEach(func(ref *plumbing.Reference) error {
		hash := ref.Hash()
		tag, err := scanner.repo.TagObject(hash)
		if err != nil {
			// nothing to do
			fmt.Printf("error on this tag %s\n", hex.EncodeToString(hash[:]))
		} else {
			target := tag.Target
			scanner.tags[hex.EncodeToString(target[:])] = tag.Name
		}
		return nil
	})
	//fmt.Printf("tags %+v\n", scanner.tags)
	return nil
}

func (scanner *ScanGit) GetTags() map[string]string {
	return scanner.tags
}
