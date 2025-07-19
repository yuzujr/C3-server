package eventbus

type Command struct {
	Type string         `json:"type"`
	Data map[string]any `json:"data"`
}

type ShellOutputMsg struct {
	Type     string `json:"type"`
	ClientID string `json:"client_id"`
	Output   string `json:"output"`
}

type StatusChangeMsg struct {
	Type     string `json:"type"`
	ClientID string `json:"client_id"`
	Online   bool   `json:"online"`
}

type NewScreenshotMsg struct {
	Type     string `json:"type"`
	ClientID string `json:"client_id"`
	URL      string `json:"url"`
}

type UpdateAliasMsg struct {
	Type     string `json:"type"`
	ClientID string `json:"client_id"`
	OldAlias string `json:"old_alias"`
	NewAlias string `json:"new_alias"`
}
