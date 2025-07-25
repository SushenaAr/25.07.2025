package model

type PackageSpec struct {
	Name string `json:"name"`
	Ver  string `json:"ver,omitempty"`
}

type PackagesList struct {
	Packages []PackageSpec `json:"packages"`
}
