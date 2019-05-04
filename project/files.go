// Package project handles project folder analysis
//
// Development note:
// List dirs with `ls -d -1 $PWD/**`
package project

import (
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/iwittkau/projects"
	times "gopkg.in/djherbis/times.v1"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func FromWorkspace(ws projects.Workspace) (pro projects.Project, err error) {
	fi, err := os.Stat(ws.Path)
	if err != nil {
		return pro, err
	}
	t, err := times.Stat(ws.Path)
	if err != nil {
		return pro, err
	}

	pro.LastEdited = fi.ModTime()
	pro.SizeBytes = fi.Size()
	pro.Active = ws.Active
	pro.Path = ws.Path
	pro.Created = t.BirthTime()

	pro.Name = ws.Name
	return pro, err
}

func ListDates(path string) (ts []time.Time, err error) {

	var fi []os.FileInfo
	if fi, err = ioutil.ReadDir(path); err != nil {
		return ts, err
	}

	for i := range fi {
		if !strings.HasPrefix(fi[i].Name(), ".") {
			ts = append(ts, fi[i].ModTime())
		}
	}

	return ts, err

}

func ListCommits(path string) (ts []time.Time, err error) {
	g, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	i, err := g.CommitObjects()
	if err != nil {
		return nil, err
	}

	defer i.Close()

	aYearAgo := time.Now().AddDate(-1, -6, 0)
	for {
		c, err := i.Next()
		if err != nil {
			break
		}
		if c.Author.When.Before(aYearAgo) {
			continue
		}
		ts = append(ts, c.Author.When)
	}

	err = i.ForEach(func(c *object.Commit) error {
		ts = append(ts, c.Author.When)
		return nil
	})

	return ts, err
}
