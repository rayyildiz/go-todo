package main

import (
    "encoding/json"
    "net/http"

    "github.com/gorilla/mux"
    "appengine"
    "fmt"
    "appengine/memcache"
    "time"
)

const frontendUrl = "http://localhost:8080"
//const frontendUrl = "http://rayyildiz-todo.appspot.com"

type Todo struct {
    Completed   bool    `json:"completed"`
    Id          int     `json:"id"`
    Title       string  `json:"title"`
}

type Seq struct {
    LatestTodoId int    `json:"seq_no"`
}

func init() {
    r := mux.NewRouter()

    r.HandleFunc("/api",status)
    r.HandleFunc("/api/todos", options).Methods("OPTIONS")
    r.HandleFunc("/api/todos/{id:[0-9]+}", options).Methods("OPTIONS")
    r.HandleFunc("/api/todos", getTodos).Methods("GET")
    r.HandleFunc("/api/todos", insertTodo).Methods("POST")
    r.HandleFunc("/api/todos/{id:[0-9]+}", updateTodo).Methods("PUT")
    r.HandleFunc("/api/todos/{id:[0-9]+}", deleteTodo).Methods("DELETE")

    http.Handle("/",r)
}

func defaultHeader(rw http.ResponseWriter){
    rw.Header().Set("Access-Control-Allow-Origin", frontendUrl)
    rw.Header().Set("Access-Control-Allow-Methods","POST, GET, OPTIONS, PUT, DELETE")
    rw.Header().Set("Access-Control-Allow-Headers","Origin, X-Requested-With, Content-Type, Accept")
}

func options(rw http.ResponseWriter,req* http.Request) {
    defaultHeader(rw)


}

func insertTodo(rw http.ResponseWriter,req* http.Request){
    defaultHeader(rw)

    c := appengine.NewContext(req)

    nextId,err1 := next_seq(c)
    if ( err1!=nil){
        c.Errorf("error at %v",err1)
    }

    var todo Todo
    decoder := json.NewDecoder(req.Body)
    err := decoder.Decode(&todo)
    if (err!=nil){
        c.Errorf("error at %v",err)
        return
    }

    todo.Id = nextId
    c.Infof("Id %v",todo)
}


func updateTodo(rw http.ResponseWriter,req* http.Request){
    defaultHeader(rw)
    c := appengine.NewContext(req)

    //params := mux.Vars(req)
    //id := params["id"]

    var todo Todo
    decoder := json.NewDecoder(req.Body)
    err := decoder.Decode(&todo)
    if (err!=nil){
        c.Errorf("error at %v",err)
        return
    }

    c.Infof("Todo for updating %v",todo)
}

func deleteTodo(rw http.ResponseWriter,req* http.Request){
    defaultHeader(rw)
    c := appengine.NewContext(req)

    params := mux.Vars(req)
    id := params["id"]

    c.Infof("Todo for deleting %v",id)
}


func getTodos(rw http.ResponseWriter, req* http.Request) {
    defaultHeader(rw)

    c := appengine.NewContext(req)

    var res struct {
        Todos   []  *Todo    `json:"todos"`
        Errors  []  string   `json:"err"`
    }


    for i:=1 ; i < 2 ; i++ {
        todo, err1 := fetch(c,i)
        if ( err1!= nil){
            c.Errorf("encode response: %v", err1)
            res.Errors = append(res.Errors,err1.Error())
        } else {
            res.Todos = append(res.Todos, todo)
        }
    }

    enc := json.NewEncoder(rw)
    err := enc.Encode(res)

    // And if encoding fails we log the error
    if err != nil {
        c.Errorf("encode response: %v", err)
    }
}

func status(rw http.ResponseWriter, req *http.Request){
    defaultHeader(rw)
}

const SequenceName = "TodoSeq"

func next_seq(c appengine.Context)(int,error){
    seq := &Seq{}
    _, err := memcache.JSON.Get(c, SequenceName, seq)
    if ( err == nil) {
        return increment(c,*seq), nil
    }

    if err != memcache.ErrCacheMiss {
        c.Errorf("memcache get %q: %v", "latestId", err)
        return 0,err
    }

    todos  := defaultTodo(c)

    // insert todos
    c.Infof(fmt.Sprintf("%v",todos))


    return increment(c,Seq{LatestTodoId:cap(todos)}),nil
}

func increment(c appengine.Context,seq Seq) int {
    next := seq.LatestTodoId + 1
    item := &memcache.Item{
        Key:        SequenceName,
        Object:     Seq{LatestTodoId:next} ,
        Expiration: time.Hour,
    }

    err := memcache.JSON.Set(c, item)
    if ( err != nil){
        c.Errorf("memcache set %q: %v", SequenceName, err)
    }

    return next
}


func fetch(c appengine.Context,id int) (*Todo, error) {
    return &Todo{
        Completed:      false,
        Id:             id,
        Title:          "This is a test todo",
    }, nil
}


func defaultTodo(c appengine.Context)([] Todo) {
    todos := []Todo{
        {false,1,"My first todo"} ,
        {false,2,"My second todo"},
        {true,3,"My comleted todo"},
    }

    return todos
}