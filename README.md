# fiber-skeleton

## Preface

I started to participate in the development of [Goravel](https://github.com/goravel/goravel) last year. I have learned a lot of new things in this year. I think an excellent skeleton should meet the following points:

* Simple directory (only contains the minimum directory structure required for CURD application)
* Concise code (only contains the minimum amount of code required for CURD application)
* Complete functions (from form validation to ORM and other functions required by CURD application)
* Easy to expand (can be expanded quickly based on CURD)
* Graceful restart (supports air hot restart of development environment and graceful restart of online environment)

Goravel is an excellent framework that is very easy to use, but its design that completely copies Laravel makes it a heavyweight framework, which does not meet the above design concepts. Therefore, this skeleton was created.

Unlike [chi-skeleton](https://github.com/go-rat/chi-skeleton), this skeleton uses the incredibly fast [Fiber](https://gofiber.io/) framework, which is generally recommended.

## Design

According to the above design concept, the final directory structure of this skeleton is as follows:

* The cmd directory stores the entry file of the application, one file for each application
* The config directory stores configuration files, which can have multiple configuration files
* The internal directory stores various codes of the application
* The mocks directory stores the generated mock code for testing
* The pkg directory stores some packages that can be reused by the application
* The storage directory stores files generated when the application is running
* The web directory stores the front-end code of the application
* go.mod and go.sum are used to manage dependencies

The internal directory refers to the design of [Kratos](https://go-kratos.dev/), dividing the application into three layers: biz, data, and service, which are responsible for business logic, data access, and service layers respectively.

## TODO

* [x] support protobuf
* [x] code generation tool

## Thanks

The development of this project refers to the following projects, I would like to express my gratitude:

* [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
* [Kratos](https://go-kratos.dev/)
* [Goravel](https://github.com/goravel/goravel)
* [Fiber backend template](https://github.com/create-go-app/fiber-go-template)
* [GinSkeleton](https://github.com/qifengzhang007/GinSkeleton)
* [gin-layout](https://github.com/wannanbigpig/gin-layout)
