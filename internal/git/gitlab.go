package git

import "time"

type Project struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	PathWithNamespace string `json:"path_with_namespace"`
}

type Pipeline struct {
	ID        int       `json:"id"`
	Status    string    `json:"status"`
	Source    string    `json:"source"`
	Ref       string    `json:"ref"`
	WebURL    string    `json:"web_url"`
	Duration  int       `json:"duration"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
