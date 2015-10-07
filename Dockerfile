FROM golang:1.4

ADD . /go/src/github.com/paked/steel

RUN go get github.com/codegangsta/negroni
RUN go get github.com/dgrijalva/jwt-go
RUN go get github.com/gorilla/mux
RUN go get github.com/gorilla/context
RUN go get github.com/mattn/go-sqlite3
RUN go get github.com/paked/configure
RUN go get github.com/paked/gerrycode/communicator
RUN go get github.com/paked/restrict
RUN go get -v github.com/codegangsta/gin

WORKDIR /go/src/github.com/paked/steel

RUN rm database.db
RUN go run cmd/create_db/create_db.go
