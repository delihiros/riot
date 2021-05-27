package v1

type PlatformDataDto struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Locales      []string     `json:"locales"`
	Maintenances []*StatusDto `json:"maintenances"`
	Incidents    []*StatusDto `json:"incidents"`
}

type StatusDto struct {
	ID                int           `json:"id"`
	MaintenanceStatus string        `json:"maintenance_status"`
	IncidentSeverity  string        `json:"incident_severity"`
	Titles            []*ContentDto `json:"titles"`
	Updates           []*UpdateDto  `json:"updates"`
	CreatedAt         string        `json:"created_at"`
	ArchiveAt         string        `json:"archive_at"`
	UpdatedAt         string        `json:"updated_at"`
	Platforms         []string      `json:"platforms"`
}

type ContentDto struct {
	Locale  string `json:"locale"`
	Content string `json:"content"`
}

type UpdateDto struct {
	ID               int           `json:"id"`
	Author           string        `json:"author"`
	Publish          bool          `json:"publish"`
	PublishLocations []string      `json:"publish_locations"`
	Translations     []*ContentDto `json:"translations"`
	CreatedAt        string        `json:"created_at"`
	UpdatedAt        string        `json:"updated_at"`
}
