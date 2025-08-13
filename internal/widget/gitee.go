package widget

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// GiteeReposWidget Gitee仓库组件
type GiteeReposWidget struct {
	ChineseWidget
	Repositories []string `yaml:"repositories"`
	Token        string   `yaml:"token"`
	ShowIssues   bool     `yaml:"show-issues"`
	ShowPRs      bool     `yaml:"show-prs"`
	Limit        int      `yaml:"limit"`
}

type GiteeRepoData struct {
	Name            string    `json:"name"`
	FullName        string    `json:"full_name"`
	Description     string    `json:"description"`
	URL             string    `json:"url"`
	Language        string    `json:"language"`
	Stars           int       `json:"stars"`
	Forks           int       `json:"forks"`
	Issues          int       `json:"issues"`
	PullRequests    int       `json:"pull_requests"`
	LastCommit      time.Time `json:"last_commit"`
	LatestRelease   string    `json:"latest_release,omitempty"`
	ReleaseURL      string    `json:"release_url,omitempty"`
}

type GiteeAPIRepo struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	HTMLURL     string `json:"html_url"`
	Language    string `json:"language"`
	StargazersCount int `json:"stargazers_count"`
	ForksCount      int `json:"forks_count"`
	OpenIssuesCount int `json:"open_issues_count"`
	UpdatedAt       string `json:"updated_at"`
}

func NewGiteeReposWidget() *GiteeReposWidget {
	return &GiteeReposWidget{
		ChineseWidget: ChineseWidget{
			BaseWidget: BaseWidget{
				Type: "gitee-repos",
			},
			Region:    "cn",
			APISource: "gitee",
		},
		ShowIssues: true,
		ShowPRs:    true,
		Limit:      10,
	}
}

func (g *GiteeReposWidget) GetData(ctx context.Context, config Config) (interface{}, error) {
	var allRepos []GiteeRepoData
	
	for _, repo := range g.Repositories {
		repoData, err := g.fetchRepository(ctx, repo)
		if err != nil {
			// 记录错误但继续处理其他仓库
			continue
		}
		allRepos = append(allRepos, *repoData)
	}
	
	// 限制数量
	if len(allRepos) > g.Limit {
		allRepos = allRepos[:g.Limit]
	}
	
	return map[string]interface{}{
		"repositories": allRepos,
		"show_issues":  g.ShowIssues,
		"show_prs":     g.ShowPRs,
		"title":        g.getTitle(),
	}, nil
}

func (g *GiteeReposWidget) fetchRepository(ctx context.Context, repoPath string) (*GiteeRepoData, error) {
	// repoPath 格式: owner/repo
	parts := strings.Split(repoPath, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid repository format: %s", repoPath)
	}
	
	url := fmt.Sprintf("https://gitee.com/api/v5/repos/%s", repoPath)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	// 添加认证token（如果提供）
	if g.Token != "" {
		req.Header.Set("Authorization", "token "+g.Token)
	}
	
	req.Header.Set("User-Agent", "Glance-China/1.0")
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var apiRepo GiteeAPIRepo
	if err := json.NewDecoder(resp.Body).Decode(&apiRepo); err != nil {
		return nil, err
	}
	
	// 解析更新时间
	updatedAt, _ := time.Parse(time.RFC3339, apiRepo.UpdatedAt)
	
	repoData := &GiteeRepoData{
		Name:         apiRepo.Name,
		FullName:     apiRepo.FullName,
		Description:  apiRepo.Description,
		URL:          apiRepo.HTMLURL,
		Language:     apiRepo.Language,
		Stars:        apiRepo.StargazersCount,
		Forks:        apiRepo.ForksCount,
		Issues:       apiRepo.OpenIssuesCount,
		LastCommit:   updatedAt,
	}
	
	// 获取最新发布版本（如果需要）
	if release, err := g.fetchLatestRelease(ctx, repoPath); err == nil {
		repoData.LatestRelease = release.TagName
		repoData.ReleaseURL = release.HTMLURL
	}
	
	return repoData, nil
}

type GiteeRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
}

func (g *GiteeReposWidget) fetchLatestRelease(ctx context.Context, repoPath string) (*GiteeRelease, error) {
	url := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/releases/latest", repoPath)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	if g.Token != "" {
		req.Header.Set("Authorization", "token "+g.Token)
	}
	
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var release GiteeRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}
	
	return &release, nil
}

func (g *GiteeReposWidget) GetCacheKey(config Config) string {
	return fmt.Sprintf("gitee-repos:%d", len(g.Repositories))
}

func (g *GiteeReposWidget) getTitle() string {
	if g.Title != "" {
		return g.Title
	}
	return "Gitee 仓库"
}

func (g *GiteeReposWidget) Validate(config Config) error {
	if len(g.Repositories) == 0 {
		return fmt.Errorf("至少需要配置一个仓库")
	}
	
	for _, repo := range g.Repositories {
		if !strings.Contains(repo, "/") {
			return fmt.Errorf("仓库格式错误，应为 owner/repo: %s", repo)
		}
	}
	
	return nil
}
