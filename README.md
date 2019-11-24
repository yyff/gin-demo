# gin-demo
go gin demo using Elastic apm


## build
```
go build -v
```

## run

### run dependency with docker
```
docker-compose up -d
```


### run gin-demo

```
./gin-demo -listen=:8000 -db="mysql:apm_user:apm_passwd@tcp(localhost:3306)/apm_db?charset=utf8"
```

