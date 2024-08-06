package utils

type Path struct {
	ImageLink string `json:"image-link"`
	Name      string `json:"name"`
	Paths     *Path  `json:"path"`
}
