[![Build Status](https://api.travis-ci.org/gowebtw/goslim.svg?branch=master)](https://api.travis-ci.org/gowebtw/goslim.svg)

# Goslim
Goslim is a Web framework written in Go. And this is not a total complete project.

I just did a minimum workable architechture.

The workable feature is:

- Router (based on [httprouter] package)
- Context (It can easy to render view and print json string)
- Middleware
- Logger middleware
- Log package
- Read config from a yaml file

And Model is not complete yet. Just finished MySQL SELECT statement Builder

##### all feature you can see in Example Directory.

# Install
```
$ go get github.com/gowebtw/goslim
```

# Run demo
```
$ cd $GOPATH/src/github.com/gowebtw/goslim/Example
$ go run main.go
```

# Your project file structure
    |-- $GOPATH
    |   |-- src
    |       |--Your_Project_Name
    |          |-- config
    |              |-- default.yaml
    |          |-- controllers
    |              |-- default.go
    |          |-- log
    |          |-- models
    |          |-- routers
    |              |-- routers.go
    |          |-- static
    |          |-- views
    |          |-- main.go

# Quick start

### Import
```go
import (
    "Your_Project_Name/routers"
    "github.com/gowebtw/goslim"
    "github.com/gowebtw/goslim/middleware"
)
```

### New
```go
g := goslim.New()
```

### Register Routes
```go
routers.RegistRout(g.Router)
```
Then in your routers.go

```go
package routers

import (
    "Your_Project_Name/controllers"
    "github.com/gowebtw/goslim"
)

func RegistRout(r *goslim.Router)  {
    
    r.Get("/", controllers.IndexPage)
    r.Post("/post/:param", controllers.PostTest)
    
    rc := &controllers.RestController{}
    r.REST("/User", rc)
    
}
```

### Register middleware
```go
g.Router.Use(middleware.LogMiddleware)
```

And you can write your own middleware function

```go
func LogMiddleware(next goslim.CHandler) goslim.CHandler {
    return func (c *goslim.Context) error  {
       
       // do something before next handler
       
       err := next(c)
       
       // do something after next handler
       
       return err
    }
}
```

### And done

```go
g.Run()
```

### MySQL DB BuilderInterface works progress
- [x] SetDB
- [x] GetDB
- [x] Select
 - [x] Distinct
 - [x] From
 - [x] Where
 - [x] AndWhere
 - [x] OrWhere
 - [x] GroupBy
 - [x] Having
 - [x] OrderBy
 - [x] Desc
 - [x] Asc
 - [x] Limit
 - [x] Get
 - [x] Count
 - [x] Insert
- [ ] Join
 - [ ] Inner join
 - [ ] Outer join
 - [ ] Left join
 - [ ] Right join
- [ ] Update
- [ ] Delete
- [x] Transaction
 - [x] begin
 - [x] commit
 - [x] rollback

The BuilderInterface at model/interface_builder.go

### Roadmap
- [ ] Models
 - [ ] Model fields mapping
 - [ ] ORM
 - [ ] Relation mapping
 - [x] Transaction
 - [ ] QueryBuilder
- [ ] Session
 - [ ] Filesystem
 - [ ] Database
 - [ ] Redis
 - [ ] Memcache
- [ ] Cache
 - [ ] Memory
 - [ ] File
 - [ ] Redis
 - [ ] Memcache
- [ ] i18n
- [ ] Command line tools
- [ ] Form handler (maybe next version)
- [ ] Security check features(csrf, xss filter...etc)


[httprouter]: <https://github.com/julienschmidt/httprouter>
