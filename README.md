# fiber-skeleton

## 前言

我从去年开始参与 [Goravel](https://github.com/goravel/goravel) 的开发，在这一年多的时间里学到了很多新东西，我认为一个优秀的脚手架应满足以下几点：

* 目录简单（只包含 CURD 应用所需的最小目录结构）
* 代码简洁（只包含 CURD 应用所需的最小代码量）
* 功能完备（从表单验证到 ORM 等 CURD 应用所需的功能都具备）
* 易于拓展（在 CURD 的基础上可以快速拓展）
* 优雅重启（支持开发环境的 air 热重启和线上环境的优雅重启）

Goravel 是一个非常容易上手的优秀框架，但其完全复制 Laravel 的设计导致了它是一个重量级框架，不符合以上设计理念。因此便有了本脚手架。

和 [chi-skeleton](https://github.com/go-rat/chi-skeleton) 不同，此脚手架使用了速度奇快的 [Fiber](https://gofiber.io/) 框架，通常建议使用此脚手架。

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
