package sound

type Manifest struct {
	Name       string                      `json:"name"`
	DisplayName string                     `json:"display_name,omitempty"`
	Categories map[string]ManifestCategory `json:"categories"`
}

type ManifestCategory struct {
	Sounds []ManifestSound `json:"sounds"`
}

type ManifestSound struct {
	File  string `json:"file"`
	Label string `json:"label,omitempty"`
}
