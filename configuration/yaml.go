package configuration

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/iwittkau/projects"
	yaml "gopkg.in/yaml.v2"
)

func ReadConfigFromHomeDir() (conf projects.Configuration, err error) {

	path := ".projects"

	usr, err := user.Current()
	if err != nil {
		return conf, err
	}

	path = usr.HomeDir + string(os.PathSeparator) + path

	return ReadConfigFromFile(path)

}

func ReadConfigFromFile(path string) (conf projects.Configuration, err error) {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return conf, err
	}
	err = yaml.Unmarshal(data, &conf)
	return conf, err
}

func WriteToFile(conf projects.Configuration) error {

	data, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}

	return ioutil.WriteFile("projects.yaml", data, os.ModePerm)

}

func WriteToHomeDir(conf projects.Configuration) error {
	data, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}

	path := ".projects"

	usr, err := user.Current()
	if err != nil {
		return err
	}

	path = usr.HomeDir + string(os.PathSeparator) + path

	return ioutil.WriteFile(path, data, os.ModePerm)
}

func AddProject(conf *projects.Configuration, ws projects.Workspace) error {
	for i := range conf.Workspaces {
		if conf.Workspaces[i].Path == ws.Path {
			return fmt.Errorf("path '%s' already exists in configuration as '%s'", ws.Path, conf.Workspaces[i].Name)
		}
	}

	conf.Workspaces = append(conf.Workspaces, ws)

	return nil
}
