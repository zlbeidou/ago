# 概述
collector 作用是数据批量采集, 使用场景是零散数据的合并, 例如将频繁而零散的请求合并之后批量发送到后端, 避免频繁发送请求.
支持设置采集的大小和超时时间, 两种条件可以获取到合并数据, 一种是达到合并条数, 还有一种是第一条数据已经超时.

# 图解说明

```
参数: 合并大小(5), 超时(5)

time:   0--------10--------20--------30--------40--------50
input:  6          18        2 1 3 2 1      2 1 1       2 3
output: 5    1     5*3  3        5    4          4        5
```
