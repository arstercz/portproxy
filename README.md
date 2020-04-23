# portproxy

A TCP port proxy utility inspired by qtunnel(https://github.com/getqujing/qtunnel).

**note:** `portproxy` does not suport ssl mode(mysql 5.7/8.0 client), only used in test environments.

## How to Install

 ```
 go get github.com/arstercz/portproxy
 ```

## Usage

```
Usage of ./portproxy:
  -backend string
        backend server ip and port (default "127.0.0.1:8003")
  -bind string
        locate ip and port (default ":8002")
  -buffer uint
        buffer size (default 4096)
  -conf string
        config file to verify database and record sql query
  -daemon
        run as daemon process
  -logTo string
        stdout or syslog (default "stdout")
  -verbose
        print verbose sql query
```

portproxy only print mysql queries when `conf` is not set:
```
./portproxy -backend="10.0.21.5:3301" -bind=":3316" -buffer=16384  --verbose
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
20200423: skip error when does not set conf option
20170112: log mysql query
```
