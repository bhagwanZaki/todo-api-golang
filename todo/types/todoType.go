package types

type TodoSchema struct {
	Name string `json:"name"`
	Completed bool `json:"completed"`
}
type Todo struct{
	Id int `json:"id"`
	Name string `json:"name"`
	Completed bool `json:"completed"`
}