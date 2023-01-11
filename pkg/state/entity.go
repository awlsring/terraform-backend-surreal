package state

type Entity struct {
	Id string `json:"id"`
	Locked bool `json:"locked"`
	State TfState `json:"state"`
}