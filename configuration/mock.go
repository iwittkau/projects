package configuration

import (
	"github.com/iwittkau/projects"
)

func Mock() projects.Configuration {
	return projects.Configuration{
		Workspaces: []projects.Workspace{
			{
				Path:   "/Users/name/projects/project",
				Name:   "Name of the Project",
				Active: true,
			},
			{
				Path:   "/Users/name/projects/project-inactive",
				Name:   "Name of the Project",
				Active: false,
			},
		},
	}
}
