package changelog

import (
	"os"

	"gopkg.in/yaml.v3"
)

type projectConfig struct {
	Template string `yaml:"template"`
	Info     struct {
		Title         string `yaml:"title"`
		RepositoryURL string `yaml:"repository_url"`
	} `yaml:"info"`
	Options struct {
		TagFilterPattern string `yaml:"tag_filter_pattern"`
		Sort             string `yaml:"sort"`
		Commits          struct {
			Filters map[string][]string `yaml:"filters"`
			SortBy  string              `yaml:"sort_by"`
		} `yaml:"commits"`
		CommitGroups struct {
			GroupBy    string            `yaml:"group_by"`
			SortBy     string            `yaml:"sort_by"`
			TitleOrder []string          `yaml:"title_order"`
			TitleMaps  map[string]string `yaml:"title_maps"`
		} `yaml:"commit_groups"`
		Header struct {
			Pattern     string   `yaml:"pattern"`
			PatternMaps []string `yaml:"pattern_maps"`
		} `yaml:"header"`
		Notes struct {
			Keywords []string `yaml:"keywords"`
		} `yaml:"notes"`
	} `yaml:"options"`
}

func loadProjectConfig(path string) (projectConfig, error) {
	if path == "" {
		return projectConfig{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return projectConfig{}, nil
		}
		return projectConfig{}, err
	}

	var cfg projectConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return projectConfig{}, err
	}

	return cfg, nil
}
