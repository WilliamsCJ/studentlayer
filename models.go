package main

// University object, mirrors JSON schema in data/universities_clean.json.
type University struct {
	Name         string   `json:"name"`
	Domains      []string `json:"domain"`
	AlphaTwoCode []string `json:"alpha_two_code"`
}

// UniversityCollection represents an array of University objects. Mirrors the array of JSON objects seen in
// data/universities_clean.json.
type UniversityCollection []University

// Email represents an email address and its associated information, such as whether it is verified and the University
// it is associated to
type Email struct {
	Email        string               `json:"email"`
	Verified     bool                 `json:"verified"`
	Universities UniversityCollection `json:"universities,omitempty"`
}
