package boondoggle

import (
	"strings"
	"github.com/spf13/viper"
	"os"
	"fmt"
)

// Requirements represents the yaml file requirements.yaml as used in a Helm chart.
type Requirements struct {
	Dependencies []Dependency `yaml:"dependencies"`
}

// Dependency is part of the requirements file
type Dependency struct {
	Name         string        `yaml:"name,omitempty"`
	Version      string        `yaml:"version,omitempty"`
	Repository   string        `yaml:"repository,omitempty"`
	Condition    string        `yaml:"condition,omitempty"`
	Tags         []string      `yaml:"tags,omitempty"`
	Enabled      bool          `yaml:"enabled,omitempty"`
	Importvalues []interface{} `yaml:"importvalues,omitempty"`
	Alias        string        `yaml:"alias,omitempty"`
}

//BuildRequirements converts a Boondoggle into a Helm Requirements struct.
func BuildRequirements(b Boondoggle) Requirements {
	var r Requirements
	var repoLocation string
	for _, service := range b.Services {

		if service.Repository == "localdev" {
			repoLocation = fmt.Sprintf("file://%s/%s/%s", os.Getenv("PWD"), service.Path ,service.Chart)
		} else {
			repoLocation = service.Repository
		}

		version := getVersionFlag(service)

		var dependency = Dependency{
			Name:         service.Chart,
			Version:      version,
			Repository:   repoLocation,
			Condition:    service.Condition,
			Tags:         service.Tags,
			Enabled:      service.Enabled,
			Importvalues: service.Importvalues,
			Alias:        service.Alias,
		}

		r.Dependencies = append(r.Dependencies, dependency)
	}
	return r
}

func getVersionFlag(service Service) string {
	// for each of the --state-v-override flags...
	for _, override := range viper.GetStringSlice("state-v-override") {
		//Split into slice on the "="
		splitvalue := strings.Split(override, "=")
		// if the left side of the value equals the service name, return the version.
		if splitvalue[0] == service.Name {
			return splitvalue[1]
		}
	}
	// if no match found, return the original version.
	return service.Version
}

func getLocalRepoLocation(s Service, l string) string {
	var repo string
	repo = fmt.Sprintf("file://%s/%s/%s", l, s.Name, s.Chart)
	return repo
}
