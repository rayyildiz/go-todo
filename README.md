# An implementation  [TodoMVC](http://todomvc.com) by using [Go Lang](//golang.org/)

[![Build Status](https://travis-ci.org/rayyildiz/go-todo.svg?branch=master)](https://travis-ci.org/rayyildiz/go-todo)


Demo [http://rayyildiz-todo.appspot.com](http://rayyildiz-todo.appspot.com)


Install
====

Install [go AppEngine](https://cloud.google.com/appengine/docs/go) then set GOPATH and GOROOT

Install  [http://www.gorillatoolkit.org/pkg/mux](http://www.gorillatoolkit.org/pkg/mux)

    go get github.com/gorilla/mux


Run app engine with

    goapp serve frontend/frontend.yaml  backend/backend.yaml


Deploy
====

If you want to deploy to change `application` in `backend/backend.yaml` and `frontend/frontend.yaml`

    application: rayyildiz-todo

Then change your api url in `frontend/static/js/services/todoStorage.js`  

    var baseApiPath = "//api.rayyildiz-todo.appspot.com"

Don't forget to change frontend url in `backend/tasks.go` like

    const frontendUrl = "http://localhost:8080"

Now you can deploy your application

    goapp  deploy --oauth frontend/frontend.yaml  backend/backend.yaml
