package helm

import "time"

// https://helm.sh/docs/topics/charts/
type Repo struct {
	ApiVersion string             `yaml:"apiVersion"`
	EntriesMap map[string][]Chart `yaml:"entries"`
	Generated  time.Time          `yaml:"generated"`
}

type Maintainer struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
	Url   string `yaml:"url"`
}

type Dependency struct {
	Name         string        `yaml:"name"`
	Version      string        `yaml:"version"`
	Repository   string        `yaml:"repository"`
	Condition    string        `yaml:"condition"`
	Tags         []string      `yaml:"tags"`
	ImportValues []interface{} `yaml:"import-values"`
	Alias        string        `yaml:"alias"`
}

type Chart struct {
	ApiVersion  string       `yaml:"apiVersion"`
	AppVersion  string       `yaml:"appVersion"`
	Created     time.Time    `yaml:"created"`
	Description string       `yaml:"description"`
	Digest      string       `yaml:"digest"`
	Engine      string       `yaml:"engine"`
	Home        string       `yaml:"home"`
	Icon        string       `yaml:"icon"`
	Keywords    []string     `yaml:"keywords"`
	Maintainers []Maintainer `yaml:"maintainers"`
	Name        string       `yaml:"name"`
	Sources     []string     `yaml:"sources"`
	Urls        []string     `yaml:"urls"`
	Version     string       `yaml:"version"`
}
