handlers [![Build Status](https://travis-ci.org/issue9/handlers.svg?branch=master)](https://travis-ci.org/issue9/handlers)
======

handlers 是对 http.Handler 接口的一系列实现，提供了大部分实用的功能。

- version 匹配从 Accept 报头中的版本号信息；
- comporess 对内容进行压缩；
- host 匹配指定的域名；
- recovery 对 Panic 的处理；


### 安装

```shell
go get github.com/issue9/handlers
```


### 文档

[![Go Walker](https://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/issue9/handlers)
[![GoDoc](https://godoc.org/github.com/issue9/handlers?status.svg)](https://godoc.org/github.com/issue9/handlers)


### 版权

本项目采用 [MIT](https://opensource.org/licenses/MIT) 开源授权许可证，完整的授权说明可在 [LICENSE](LICENSE) 文件中找到。
