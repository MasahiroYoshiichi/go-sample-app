package controller

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"

	"github.com/MasahiroYoshiichi/go-workout-app/controller/dto"
	"github.com/MasahiroYoshiichi/go-workout-app/model/entity"
	"github.com/MasahiroYoshiichi/go-workout-app/model/repository"
)

// 外部パッケージに公開するインタフェース
type TodoController interface {
	GetTodos(w http.ResponseWriter, r *http.Request)
	PostTodo(w http.ResponseWriter, r *http.Request)
	PutTodo(w http.ResponseWriter, r *http.Request)
	DeleteTodo(w http.ResponseWriter, r *http.Request)
}

// 非公開のTodoController構造体
type todoController struct {
	tr repository.TodoRepository
}

// TodoControllerのコンストラクタ。
// 引数にTodoRepositoryを受け取り、TodoController構造体のポインタを返却する。
func NewTodoController(tr repository.TodoRepository) TodoController {
	return &todoController{tr}
}

// TODOの取得
func (tc *todoController) GetTodos(w http.ResponseWriter, r *http.Request) {
	// リポジトリの取得処理呼び出し
	todos, err := tc.tr.GetTodos()
	if err != nil {
		w.WriteHeader(500)
		return
	}

	// 取得したTODOのentityをDTOに詰め替え
	var todoResponses []dto.TodoResponse
	for _, v := range todos {
		todoResponses = append(todoResponses, dto.TodoResponse{Id: v.Id, Title: v.Title, Content: v.Content})
	}

	var todosResponse dto.TodosResponse
	todosResponse.Todos = todoResponses

	// JSONに変換
	output, _ := json.MarshalIndent(todosResponse.Todos, "", "\t\t")

	// JSONを返却
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

// TODOの追加
func (tc *todoController) PostTodo(w http.ResponseWriter, r *http.Request) {
	// リクエストbodyのJSONをDTOにマッピング
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	var todoRequest dto.TodoRequest
	json.Unmarshal(body, &todoRequest)

	// DTOをTODOのEntityに変換
	todo := entity.TodoEntity{Title: todoRequest.Title, Content: todoRequest.Content}

	// リポジトリの追加処理呼び出し
	id, err := tc.tr.InsertTodo(todo)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	// LocationにリソースのPATHを設定し、ステータスコード２０１を返却
	w.Header().Set("Location", r.Host+r.URL.Path+strconv.Itoa(id))
	w.WriteHeader(201)
}

// TODOの更新
func (tc *todoController) PutTodo(w http.ResponseWriter, r *http.Request) {
	// URLのPATHに含まれるTODOのIDを取得
	todoId, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		w.WriteHeader(400)
		return
	}

	// リクエストbodyのJSONをDTOにマッピング
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	var todoRequest dto.TodoRequest
	json.Unmarshal(body, &todoRequest)

	// DTOをTODOのEntityに変換
	todo := entity.TodoEntity{Id: todoId, Title: todoRequest.Title, Content: todoRequest.Content}

	// リポジトリの更新処理呼び出し
	err = tc.tr.UpdateTodo(todo)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	// ステータスコード204を返却
	w.WriteHeader(204)
}

// TODOの削除
func (tc *todoController) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	// URLのPATHに含まれるTODOのIDを取得
	todoId, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		w.WriteHeader(400)
		return
	}

	// リポジトリの削除処理呼び出し
	err = tc.tr.DeleteTodo(todoId)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	// ステータスコード204を返却
	w.WriteHeader(204)
}
