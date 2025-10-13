# fiber-skeleton

和 [chi-skeleton](https://github.com/libtnb/chi-skeleton) 不同，此脚手架使用了速度奇快的 [Fiber](https://gofiber.io/) 框架，通常建议使用此脚手架。

## 设计

按以上设计理念，最终本脚手架的目录结构如下：

* cmd 目录存放应用的入口文件，每个应用一个文件
* config 目录存放配置文件，可以有多种配置文件
* internal 目录存放应用的各种代码
* mocks 目录存放生成的 mock 代码，用于测试
* pkg 目录存放可以被应用重复使用的一些包
* storage 目录存放应用运行时产生的文件
* web 目录存放应用的前端代码
* go.mod 和 go.sum 用于管理依赖

其中 internal 目录参考了 [Kratos](https://go-kratos.dev/) 的设计，将应用分为 biz、data 和 service 三层，分别负责业务逻辑、数据访问和服务层。

## TODO

* [x] 支持 protobuf
* [x] 代码生成工具

## 致谢

本项目的开发中参考了以下项目，特此感谢：

* [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
* [Kratos](https://go-kratos.dev/)
* [Goravel](https://github.com/goravel/goravel)
* [Fiber backend template](https://github.com/create-go-app/fiber-go-template)
* [GinSkeleton](https://github.com/qifengzhang007/GinSkeleton)
* [gin-layout](https://github.com/wannanbigpig/gin-layout)
