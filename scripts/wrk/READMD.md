# Basic

教程： https://zhuanlan.zhihu.com/p/552756287

```shell
wrk -t1 -d1s -c2 -s ./scripts/wrk/profile.lua  http://localhost:8080/users/profile
```

- -t 线程数量
- -d 持续时间 1s/1m
- -c 并发数
- -s 测试脚本

## 安装

```shell
sudo apt-get install build-essential libssl-dev git -y
git clone https://github.com/wg/wrk.git wrk
cd wrk
make
# 将可执行文件移动到 /usr/local/bin 位置
sudo cp wrk /usr/local/bin
```

# profile

## MySQL

```shell
❯ wrk -t2 -d10s -c200 -s ./scripts/wrk/profile.lua  http://localhost:8080/users/profile                              ─╯
Running 10s test @ http://localhost:8080/users/profile
  2 threads and 200 connections (2个线程，200个并发连接)
             （平均值）（标准差）（最大值）（正负一个标准差所占比例）
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    81.05ms  188.89ms   1.23s    92.39%
    Req/Sec     3.37k   540.16     4.65k    80.81%
  66557 requests in 10.03s, 19.92MB read
  Socket errors: connect 0, read 0, write 0, timeout 4
Requests/sec:   6638.57  (QPS ,即平均每秒处理请求数为 6638.57)
Transfer/sec:      1.99MB (平均每秒流量)

# 出现了较多慢日志，> 200ms
```

## Redis

改用缓存

```shell
❯ wrk -t2 -d10s -c200 -s ./scripts/wrk/profile.lua  http://localhost:8080/users/profile                                 ─╯
Running 10s test @ http://localhost:8080/users/profile
  2 threads and 200 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    14.79ms   12.38ms  97.99ms   72.58%
    Req/Sec     7.61k     0.97k    9.25k    76.26%
  150147 requests in 10.03s, 45.53MB read
Requests/sec:  14975.83
Transfer/sec:      4.54MB
```
