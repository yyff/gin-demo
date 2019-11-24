FROM golang:stretch
ENV GO111MODULE=on
#COPY go.mod go.sum /src/gin-demo/
COPY . /src/gin-demo/
COPY db /
WORKDIR /src/gin-demo
RUN go mod download
RUN go build -v
RUN mv /src/gin-demo/gin-demo /
RUN rm -rf /src/gin-demo
EXPOSE 8000
HEALTHCHECK NONE
CMD ["/gin-demo", "-listen=:8000", "-db=mysql:apm_user:apm_passwd@tcp(localhost:3306)/apm_db?charset=utf8"]
