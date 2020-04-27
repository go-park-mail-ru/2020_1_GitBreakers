package git

type Branch struct {
	Name   string `json:"name" valid:"-"`
	Commit Commit `json:"commit" valid:"-"`
}
//easyjson:json
type BranchSet []Branch