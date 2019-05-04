package projects

import (
	"time"
)

type Project struct {
	Path       string
	Name       string
	LastEdited time.Time
	Created    time.Time
	SizeBytes  int64
	Active     bool
}

type Workspace struct {
	Path   string
	Name   string
	Active bool
}

type Configuration struct {
	Workspaces []Workspace
}

// ByWorkSpaceName is a type that provides a sort interface implementation for sorting by Workspace name
type ByWorkSpaceName []Workspace

// Len is part of sort.Interface.
func (w ByWorkSpaceName) Len() int {
	return len(w)
}

// Swap is part of sort.Interface.
func (w ByWorkSpaceName) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (w ByWorkSpaceName) Less(i, j int) bool {
	return w[i].Name < w[j].Name
}

// ByProjectLastEdited is a type that provides a sort interface implementation for sorting by Workspace name
type ByProjectLastEdited []Project

// Len is part of sort.Interface.
func (p ByProjectLastEdited) Len() int {
	return len(p)
}

// Swap is part of sort.Interface.
func (p ByProjectLastEdited) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (p ByProjectLastEdited) Less(i, j int) bool {
	if p[i].Active && p[j].Active || !p[i].Active && !p[j].Active {
		return p[i].LastEdited.Before(p[j].LastEdited)
	}
	if p[i].Active && !p[j].Active {
		return false
	}
	return true
}
