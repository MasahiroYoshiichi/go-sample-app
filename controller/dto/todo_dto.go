package dto

type TodoResponse struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type TodoRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type TodosResponce struct {
	Todos []TodoResponse `json:"todos"`
}
