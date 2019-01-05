package models

type Runner struct {
	Name     string `json:"name"`
	Position string `json:"position"`
	Category string `json:"category"`
	Club     string `json:"club"`
	Time     string `json:"time"`
}

// Result data structure of a result
type Result struct {
	ID              int      `json:"id"`
	Race            string   `json:"race"`
	Date            string   `json:"date"`
	NumberOfRunners int      `json:"numberOfRunners"`
	Runners         []Runner `json:"runners"`
}
