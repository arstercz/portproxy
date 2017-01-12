# portproxy
A TCP port proxy utility inspired by qtunnel(https://github.com/getqujing/qtunnel)

## How to Build

 ```
 git clone https://github.com/chenzhe07/portproxy.git
 cd portproxy
 go build -o portproxy *.go
 ```

## Usage

```
Usage of ./portproxy:
  -backend="127.0.0.1:8003": backend server ip and port
  -bind=":8002": locate ip and port
  -buffer=4096: buffer size
  -daemon=false: run as daemon process
  -logTo="stdout": stdout or syslog
```
portproxy can also log mysql queries:
```
./portproxy -backend="10.0.21.5:3301" -bind=":3316" -buffer=16384 
2017/01/12 17:27:23 portproxy started.
2017/01/12 17:27:32 client: 10.0.21.7:29110 ==> 10.0.21.5:3316
2017/01/12 17:27:32 proxy: 10.0.21.5:18386 ==> 10.0.21.5:3301
2017/01/12 17:27:32 From 10.0.21.7:29110 To 10.0.21.5:3301; Query: select @@version_comment limit 1
2017/01/12 17:27:48 From 10.0.21.7:29110 To 10.0.21.5:3301; Query: SELECT DATABASE()
2017/01/12 17:27:48 From 10.0.21.7:29110 To 10.0.21.5:3301; schema: use percona
2017/01/12 17:27:48 From 10.0.21.7:29110 To 10.0.21.5:3301; Query: show databases
2017/01/12 17:27:49 From 10.0.21.7:29110 To 10.0.21.5:3301; Query: show tables
2017/01/12 17:27:49 From 10.0.21.7:29110 To 10.0.21.5:3301; Query: table columns list: item
2017/01/12 17:27:49 From 10.0.21.7:29110 To 10.0.21.5:3301; Query: table columns list: stock
2017/01/12 17:27:56 From 10.0.21.7:29110 To 10.0.21.5:3301; Query: show tables
2017/01/12 17:28:01 From 10.0.21.7:29110 To 10.0.21.5:3301; Query: show create table item
2017/01/12 17:28:04 From 10.0.21.7:29110 To 10.0.21.5:3301; Query: kill 2
```

## changelog:
```
20170112: log mysql query
```
