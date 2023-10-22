# 3DQuestBackend
Backend for a 3D Printing Service written in [Go](https://go.dev/) using [echo](https://github.com/labstack/echo) and [CouchDB](https://couchdb.apache.org/). 


## Development Setup

In order to set up this service, you need to download [CouchDB](https://couchdb.apache.org/), as well as the [Go](https://go.dev/) programming language.

### CouchDB Installation

Install CouchDB following the instructions on the site. We recommend the following login combination for the testing environment:
- user: admin
- password: admin


Once done, check that the installation is correct with the following command:

`λ curl -X GET http://admin:admin@127.0.0.1:5984/_all_dbs`

The response should be `[_replicators, _users]`. If this is not correct, you must create the databases yourself:

`λ curl -X PUT http://admin:admin@127.0.0.1:5984/_users`
`λ curl -X PUT http://admin:admin@127.0.0.1:5984/_replicator`

Additionally, you may create the global changes database, but this is unnecessary:

`λ curl -X PUT http://admin:admin@127.0.0.1:5984/_global_changes`

If it has been created and you want to erase it, you may use the following:

`λ curl -X DELETE http://admin:admin@127.0.0.1:5984/_global_changes`

#### CouchDB's Fauxton Database Managing Interface

Instead of using curl or equivalent software, you can use the builtin Fauxton. You may access it from:

[http://127.0.0.1:5984/_utils/](http://127.0.0.1:5984/_utils/)

Login using the administrator login information you introduced on install.

### Golang Install

You may install golang following the instructions [here](https://go.dev/doc/install)

### Golang Dependencies

Some go dependencies are required for 3DQuest's Backend:
- [godotenv](github.com/joho/godotenv) for the easy import of environment variables.
- OPTIONAL: [Go Compile Daemon](github.com/githubnemo/CompileDaemon) allows for constant compilation of 3DQuest's Backend.
- [echo](https://github.com/labstack/echo) as the main web framework
- [Default](https://github.com/creasty/defaults) as an addon for struct tags to automatically initialize to specified values
- [Go Querystring](https://github.com/google/go-querystring) as an addon for struct tags to automatically create query strings from struct values
- [go-jwd](https://github.com/golang-jwt/jwt) for Json Web Token Authentication

I recommend changing the Environment Variable `GO111MODULE` to `on`, according to your system's specs, and install in the order proposed.


## Running

This Go Project is organized as a module called `3DQuest`. To run it, you can do `λ go run 3DQuest` from within the directory. Ideally, you may also use the **Go Compile Daemon** throughout development to ensure that it is consistently up to date and to avoid waiting for compile times. To do so, after installing the Compile Daemon, you may do the following:

`λ CompileDaemon -command="./3DQuest.exe""`

### .env File

We are using a `.env` file for the configuration of the server and using the [godotenv](github.com/joho/godotenv) module to import it. It is necessary for the execution of the program, and looks like a typical `.ini` file following `KEY=VALUE` format:

```
# Server information
PORT=8082

# DB Basic Information
COUCHDB_USR="admin"
COUCHDB_PWD="admin"
COUCHDB_SCH="http"
COUCHDB_URL="localhost:5984"

# DB General Information
ALL_DBS_URL="_all_dbs"

# DB Design Information
USER_DESIGN_DOC=_design/test
ALL_USERS_VIEW=_view/all_user_view
BASIC_USERS_VIEW=_view/user_view
ADMIN_USERS_VIEW=_view/admin_user_view
```

## Documentation


