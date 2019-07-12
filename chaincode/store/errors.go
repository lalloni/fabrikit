package store

type MemberError struct {
	Kind  string `json:"kind,omitempty"`
	Tag   string `json:"tag,omitempty"`
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}
