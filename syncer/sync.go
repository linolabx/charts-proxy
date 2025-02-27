package syncer

import (
	"io"
	"log"
	"os"
	"path"

	"github.com/goccy/go-yaml"
	"github.com/linolabx/charts-proxy/config"
	"github.com/linolabx/charts-proxy/helm"
	"github.com/linolabx/charts-proxy/utils"
	resty "resty.dev/v3"
)

func Sync(repo_config *config.Repo, config *config.Config) error {
	log.Printf("Syncing repo %s (%s)", repo_config.Name, repo_config.Url)

	base_dir := path.Join(config.TargetDir, repo_config.Name)
	if err := os.MkdirAll(base_dir, 0755); err != nil {
		return err
	}

	var local_repo *helm.Repo

	content, err := os.ReadFile(path.Join(base_dir, "index.yaml"))
	if err != nil {
		if os.IsNotExist(err) {
			local_repo = nil
		} else {
			return err
		}
	} else {
		err = yaml.Unmarshal(content, &local_repo)
		if err != nil {
			return err
		}
	}

	client := resty.New().SetBaseURL(repo_config.Url)
	if config.DefaultProxy != "" {
		client.SetProxy(config.DefaultProxy)
	}

	repo_resp, err := client.R().Get("/index.yaml")
	if err != nil {
		return err
	}

	repo_raw, err := io.ReadAll(repo_resp.Body)
	if err != nil {
		return err
	}

	var remote_repo helm.Repo
	err = yaml.Unmarshal(repo_raw, &remote_repo)
	if err != nil {
		return err
	}

	if local_repo != nil && !remote_repo.Generated.After(local_repo.Generated) {
		return nil
	}

	if len(repo_config.ChartFilters) > 0 {
		filtered_charts_map := make(map[string][]helm.Chart)
		for _, filter := range repo_config.ChartFilters {
			charts, ok := remote_repo.EntriesMap[filter.Name]
			if !ok {
				continue
			}

			filtered_charts_map[filter.Name] = charts

			filtered_charts := make([]helm.Chart, 0)

			for _, chart := range charts {
				if filter.VersionsFrom != "" {
					if utils.SemverCompare(chart.Version, filter.VersionsFrom) < 0 {
						continue
					}
				}

				if filter.StableOnly && !utils.SemverIsStable(chart.Version) {
					continue
				}

				filtered_charts = append(filtered_charts, chart)
			}

			filtered_charts_map[filter.Name] = filtered_charts
		}

		remote_repo.EntriesMap = filtered_charts_map
	}

	for chart_name, charts := range remote_repo.EntriesMap {
		for _, chart := range charts {
			for _, file_url := range chart.Urls {
				file_name := path.Base(file_url)
				file_path := path.Join(base_dir, file_name)
				if _, err := os.Stat(file_path); err == nil {
					continue
				}

				log.Printf("Downloading %s@%s", file_name, chart_name)

				file_resp, err := client.R().Get(file_url)
				if err != nil {
					return err
				}

				file_obj, err := os.Create(file_path)
				if err != nil {
					return err
				}
				defer file_obj.Close()

				if _, err = io.Copy(file_obj, file_resp.Body); err != nil {
					return err
				}
			}
		}
	}

	repo_file_content, err := yaml.Marshal(remote_repo)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(base_dir, "index.yaml"), repo_file_content, 0644)
	if err != nil {
		return err
	}

	return nil
}
