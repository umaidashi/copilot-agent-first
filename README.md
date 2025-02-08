# Todo API サーバー

このリポジトリは、Go 言語と標準ライブラリのみを使用して実装されたシンプルな Todo API サーバーです。タスクの作成、一覧取得が可能です。データは SQLite データベースに保存されます。

## エンドポイント

### Hello World エンドポイント

- **URL**: `/`
- **メソッド**: `GET`
- **レスポンス**: `{"message": "hello world"}`

### タスク作成エンドポイント

- **URL**: `/tasks`
- **メソッド**: `POST`
- **リクエストボディ**: `{"title": "タスクのタイトル", "due_date": "YYYY-MM-DDTHH:MM:SSZ"}`
- **レスポンス**: 作成されたタスクの JSON 表現

### タスク一覧取得エンドポイント

- **URL**: `/tasks/list`
- **メソッド**: `GET`
- **レスポンス**: タスクのリストの JSON 表現

## 使用方法

1. リポジトリをクローンします。
2. 必要な依存関係をインストールします。
3. サーバーを起動します。

```sh
$ git clone <リポジトリのURL>
$ cd <リポジトリのディレクトリ>
$ go run main.go
```

## やり取りの履歴

### 2023-10-11

#### ユーザーリクエスト

Go 言語で標準ライブラリのみを使用して `{"message": "hello world"}` とレスポンスを返す REST API サーバーを実装してください。

#### アシスタントの応答

`main.go` ファイルを作成し、以下の内容を追加しました。

```go
package main

import (
	"encoding/json"
	"net/http"
)

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "hello world"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/", helloWorldHandler)
	http.ListenAndServe(":8080", nil)
}
```

#### ユーザーリクエスト

Todo アプリのタスク作成エンドポイントを実装してください。タスクの情報はタスクのタイトルと完了日だけです。

#### アシスタントの応答

`main.go` ファイルを更新し、タスク作成エンドポイントを追加しました。

```go
package main

import (
	"encoding/json"
	"net/http"
	"time"
)

// Task represents a todo task
type Task struct {
	Title     string    `json:"title"`
	DueDate   time.Time `json:"due_date"`
}

var tasks []Task

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "hello world"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

tasks = append(tasks, task)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func main() {
	http.HandleFunc("/", helloWorldHandler)
	http.HandleFunc("/tasks", createTaskHandler)
	http.ListenAndServe(":8080", nil)
}
```

#### ユーザーリクエスト

新しいタスクを作成するための `curl` コマンドのサンプルを提供してください。

#### アシスタントの応答

新しいタスクを作成するための `curl` コマンドのサンプルを以下に示します。

```sh
curl -X POST http://localhost:8080/tasks -H "Content-Type: application/json" -d '{"title": "Buy groceries", "due_date": "2023-12-31T23:59:59Z"}'
```

#### ユーザーリクエスト

タスクの一覧を取得する API を実装してください。データは SQLite でローカルに保持します。標準パッケージのみを使用してください。

#### アシスタントの応答

`main.go` ファイルを更新し、SQLite データベースの初期化、タスク作成、タスク一覧取得を追加しました。

```go
package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
	_ "github.com/mattn/go-sqlite3"
)

// Task represents a todo task
type Task struct {
	Title   string    `json:"title"`
	DueDate time.Time `json:"due_date"`
}

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		panic(err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS tasks (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"title" TEXT,
		"due_date" TEXT
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		panic(err)
	}
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "hello world"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO tasks (title, due_date) VALUES (?, ?)", task.Title, task.DueDate.Format(time.RFC3339))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func listTasksHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT title, due_date FROM tasks")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var dueDate string
		if err := rows.Scan(&task.Title, &dueDate); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		task.DueDate, _ = time.Parse(time.RFC3339, dueDate)
		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func main() {
	initDB()

	http.HandleFunc("/", helloWorldHandler)
	http.HandleFunc("/tasks", createTaskHandler)
	http.HandleFunc("/tasks/list", listTasksHandler)
	http.ListenAndServe(":8080", nil)
}
```
