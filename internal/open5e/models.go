package open5e

type Monster struct {
	Name        string `json:"name"`
	Description string `json:"desc"`

	HP          int    `json:"hit_points"`
	AC          int    `json:"armor_class"`
	Experience  int    `json:"experience_points"`

}

type MonstersResponse struct {
	Count    int       `json:"count"`
	Next     *string   `json:"next"`
	Previous *string   `json:"previous"`
	Results  []Monster `json:"results"`
}
