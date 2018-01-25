#### beetest

#### 背景

1. 由于beego的运行机制，造成了beego无法进行测试覆盖率的分析，现在这里模拟了beego的调度方法，使用静态分发来支持代码覆盖率的测试。

#### 使用方法

```sh
cd beetest/testpkg

go test -cover -v

=== RUN   TestSimpleController

  SimpleGet ✔✔✔✔


4 total assertions


  SimplePost ✔✔✔


7 total assertions


  SimpleUpload ✔


8 total assertions

--- PASS: TestSimpleController (0.00s)
PASS
coverage: 90.9% of statements
ok  	github.com/imiskolee/beetest/testpkg	0.020s
```

#### Roadmap

1. 增加路由层判断
2. 增加filter的支持

