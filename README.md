# Breeze-go
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/weibreeze/breeze-go/blob/master/LICENSE)
[![Build Status](https://img.shields.io/travis/weibreeze/breeze-go/master.svg?label=Build)](https://travis-ci.org/weibreeze/breeze-go)
[![codecov](https://codecov.io/gh/weibreeze/breeze-go/branch/master/graph/badge.svg)](https://codecov.io/gh/weibreeze/breeze-go)
[![GoDoc](https://godoc.org/github.com/weibreeze/breeze-go?status.svg&style=flat)](https://godoc.org/github.com/weibreeze/breeze-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/weibreeze/breeze-go)](https://goreportcard.com/report/github.com/weibreeze/breeze-go)


# 概述
[Breeze](https://github.com/weibreeze/breeze)是一个跨语言序列化协议与服务描述的schema，与protobuf类似，但更加易用并且提供对旧对象的兼容能力。
Breeze-go是Breeze的go语言版本。

# 快速入门
1. 添加依赖
```shell
    go get github.com/weibreeze/breeze-go
```

2. 基础类型编解码

```go
    // 编码
    s := "just test"
    buf := breeze.NewBuffer(256)
    breeze.WriteString(buf, s, true)
    // 解码
    var ns string
    err := breeze.ReadString(breeze.CreateBuffer(buf.Bytes()), &ns)
    fmt.Printf("result:%s, err:%v\n", ns, err)
```
在明确知道类型的场景下，基础类型使用对应的编解码方法效率最高。

3. 集合类型编解码

```go
    // 编码
    m := make(map[int][]string, 16)
    m[11] = []string{"a1", "a2"}
    m[789] = []string{"a3", "a4"}
    buf := breeze.NewBuffer(256)
    breeze.WriteValue(buf, m)
    // 解码方法一
    nm := make(map[int][]string, 16)
    _, err := breeze.ReadValue(breeze.CreateBuffer(buf.Bytes()), &nm)
    fmt.Printf("result:%v, err:%v\n", nm, err)
    // 解码方法二
    i, err := breeze.ReadValue(breeze.CreateBuffer(buf.Bytes()), nil)
    fmt.Printf("result:%v, err:%v\n", i, err)
```
`WriteValue`和`ReadValue`可以实现任意类型对象的编解码，包括基础类型。

解码方法一使用变量地址作为入参，是推荐的解码方式；解码方法二适合不知道具体解码类型的场景，此时通过方法返回值获取解码结果。

4. Breeze Message编解码
```go
    // 编码
    msg := breeze.GetBenchData(1)
    buf := breeze.NewBuffer(256)
    breeze.WriteValue(buf, msg)
    // 解码
    var result breeze.TestMsg
    _, err := breeze.ReadValue(breeze.CreateBuffer(buf.Bytes()), &result)
    fmt.Printf("result:%v, err:%v\n", result, err)
```

# 使用Breeze Schema生成Message类

参见[breeze-generator](https://github.com/weibreeze/breeze-generator)

## Breeze协议说明

参考[Breeze协议说明](https://github.com/weibreeze/breeze/wiki/zh_protocol)