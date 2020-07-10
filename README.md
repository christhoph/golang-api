# GOLANG-API

- GraphiQL URI: http://localhost:3000/graphql

# Setup

#### Docker

Development

```
docker-compose up --build
```

#### Non-Docker

- Start by setting env variables for Go (Note: `$` implies using the terminal)

```
$ export GOPATH=$HOME/path/to/go/projects
$ export PATH=$GOPATH/bin:$PATH
```

- Install dependencies via Glide

```
glide install
```

- Build the server executable

```
$ go build
```

- Run the server

```
$ ./kyrene
```
