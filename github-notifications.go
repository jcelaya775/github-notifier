package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const BaseNotificationsURL = "https://api.github.com/notifications"

type GitHubAPIService struct {
	httpClient *http.Client
	token      string
}

type Notification struct {
	ID              string     `json:"id"`
	Unread          bool       `json:"unread"`
	Reason          string     `json:"reason"`
	UpdatedAt       time.Time  `json:"updated_at"`
	LastReadAt      *time.Time `json:"last_read_at,omitempty"`
	Subject         Subject    `json:"subject"`
	Repository      Repository `json:"repository"`
	URL             string     `json:"url"`
	SubscriptionURL string     `json:"subscription_url"`
}

type Subject struct {
	Title            string `json:"title"`
	URL              string `json:"url"`
	LatestCommentURL string `json:"latest_comment_url"`
	Type             string `json:"type"`
}

type Repository struct {
	ID          int64  `json:"id"`
	NodeID      string `json:"node_id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Private     bool   `json:"private"`
	Owner       Owner  `json:"owner"`
	HTMLURL     string `json:"html_url"`
	Description string `json:"description"`
}

type Owner struct {
	Login     string `json:"login"`
	ID        int64  `json:"id"`
	NodeID    string `json:"node_id"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
	Type      string `json:"type"`
	SiteAdmin bool   `json:"site_admin"`
}

func (ghClient *GitHubAPIService) GetNotifications() ([]Notification, error) {
	fmt.Println("Fetching GitHub notifications...")

	req, err := http.NewRequest(http.MethodGet, BaseNotificationsURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ghClient.token))
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotModified {
		return nil, fmt.Errorf("GitHub API returned status: %s", resp.Status)
	}
	defer resp.Body.Close()

	var notifications []Notification
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &notifications); err != nil {
		return nil, err
	}
	return notifications, nil
}
