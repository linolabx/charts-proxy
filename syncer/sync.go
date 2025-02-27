package syncer

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/linolabx/charts-proxy/config"
	"github.com/linolabx/charts-proxy/helm"
	"github.com/linolabx/charts-proxy/utils"
	"github.com/samber/lo"
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
		chart_dir := path.Join(base_dir, chart_name)
		if err := os.MkdirAll(chart_dir, 0755); err != nil {
			return err
		}

		for chart_ver_index, chart_ver := range charts {
			new_urls := make([]string, 0)

			for _, file_url := range chart_ver.Urls {
				file_ext := filepath.Ext(file_url)
				file_path := path.Join(chart_dir, fmt.Sprintf("%s-%s-%s%s", chart_name, utils.SemverNormalize(chart_ver.Version), utils.Md5sum(file_url), file_ext))
				releative_path := lo.Must(filepath.Rel(base_dir, file_path))
				new_urls = append(new_urls, releative_path)

				if _, err := os.Stat(file_path); err == nil {
					continue
				}

				log.Printf("Downloading %s", releative_path)

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

			remote_repo.EntriesMap[chart_name][chart_ver_index].Urls = new_urls
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
