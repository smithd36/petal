package models

type Plant struct {
	ID             int      `json:"id"`
	CommonName     string   `json:"common_name"`
	Slug           string   `json:"slug"`
	ScientificName string   `json:"scientific_name"`
	Year           int      `json:"year"`
	Bibliography   string   `json:"bibliography"`
	Author         string   `json:"author"`
	Status         string   `json:"status"`
	Rank           string   `json:"rank"`
	GenusID        int      `json:"genus_id"`
	ImageURL       string   `json:"image_url"`
	Synonyms       []string `json:"synonyms"`
	Genus          string   `json:"genus"`
	Family         string   `json:"family"`
}
