# Gorm系列之1

## 特别指出

**特别指出的是，该系列基础代码来自git上的开源项目7days-golang，项目地址：https://github.com/geektutu/7days-golang。**

**原项目作者：极客兔兔，个人主页：https://geektutu.com/。**

**除基础代码外，部分解释内容也摘自作者的系列博文，地址：https://geektutu.com/post/gee.html**



大神**极客兔兔**在他的博客中对该项目有自底向上的详细讲解，并将每个项目分成7天来学习，希望深入分析源码的朋友可以移步上面的传送门。相比原作者，我将更多的**自顶向下**入手，先从整体分析项目结构，再深入其中一些关键的部分，适用于**希望快速了解项目结构的人和初学者**；同时在原项目基础上做了一定的增添。

## 项目传送门

https://github.com/CAGeng/Gorm

建议clone后再继续阅读。好！

3~

2~

1~

Gorm，启程！

## 介绍

对象关系映射（Object Relational Mapping，简称ORM）是通过使用描述对象和数据库之间映射的元数据，将面向对象语言程序中的对象自动持久化到关系数据库中。Gorm是这样用Go语言写成的ORM框架，可以支持多种数据库。

## 配置环境

作为起笔的第一个系列，这里花一些篇幅简单介绍环境配置和项目导入的问题，有问题在本贴下方留言看见会速回。

### git clone到本地

### Ide

我使用的是**Goland**

### go.mod配置

我使用了go新推出的go.mod的包管理方式来取代GoPath，可以解决项目目录带来的大部分困扰，保证能“跑起来”。

**1**

需要将file->settings中的这里勾上

![image-20210421235751087](Gorm1.assets/image-20210421235751087.png)

**2**

git中的项目本身已经配置好了go.mod，如果有导入问题，可以将其删掉，重新执行下面的步骤（不是必要的）

打开terminal，在项目根目录下输入：

```
go mod init Gorm
```

需要使用自己喜欢的项目名字，修改最后一个参数就好啦。

**3**

之后在项目中类似下面这样就可以导入自己的包

```go
import (
	"Gorm/clause"
	"Gorm/log"
)
```

**4**

对于不是本项目中包导入的问题，鼠标放上去根据提示一般都能导入。不行的话可能是源的问题，要配置一下代理，就不写了。

## 安装sqlite3数据库

Gorm的目标是支持多个数据库，在开发过程中使用的demo是sqlite3，安装过程不赘述，极客兔兔的博文中也有：https://geektutu.com/post/geeorm-day1.html

## 输出器

现在可以开始动手了，先写一个Tprinter类，支持缩进输出，方便之后调试时使用。**代码位置：/log/TPrinter.go**

为了方便输出具有层级结构（这样更容易看清楚执行逻辑），Tprinter实现了两个函数用来提高和降低输出层级（缩进），同时维护一个全局的Tprinter实例，使用起来是这样：

```go
log.Mytprinter.IndentLvUp()
defer log.Mytprinter.IndentLvDown()
log.Mytprinter.Print("testRecordInit begin")
```

效果：

![image-20210422002154893](Gorm1.assets/image-20210422002154893.png)

### log

原项目中还实现了log.go,是兔兔封装的一些输出函数，使用可以参考原博客。

## Gorm初探

ORM框架是用来做什么的前面已有简单介绍，现在具体来看看Gorm是怎么用到。**Talking is cheap.Show me the code.**废话不多说，直接看代码：**位置: /session/record_test.go**。

```go
//    ./session/record_test.go
var (
	user1 = &User{"Tom", 18}
	user2 = &User{"Sam", 25}
	user3 = &User{"Jack", 25}
)
func testRecordInit(t *testing.T) *Session {
	...
    //确定表格式
	s := NewSession().Model(&User{})
    
    //删表建表
	err1 := s.DropTable()
	err2 := s.CreateTable()
    
    //插入数据
	_, err3 := s.Insert(user1, user2)
	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatal("failed init test records")
	}
	return s
}
```

简言之，利用ORM，操作数据库不需要手动使用“SELECT”等语句，同时开发语言（go）与数据库语言（mysql、sqlite等）的对象类型的转换也是自动的。

然后来看一下这个test的执行逻辑，我插入了一些Tprinter，先输出一些相对顶层的逻辑（带颜色的log输出的先不要在意）

![image-20210422010250157](Gorm1.assets/image-20210422010250157.png)

## 认识 session

看来testRecordInit进行的创建和insert操作，是由session这个目录负责执行，事实确实如此：

* table.go负责表相关操作
* record.go负责条目相关的操作
* raw.go是上面两个的基础，封装一些基础的操作

所以，**session主要的工作是维护打开的数据库对象，并封装数据库操作**，当然这只是顶层模块，不可能负责实现所有的底层逻辑，后面还会逐步深入。看一下session类的定义：

```go
//  /session/raw.go
type Session struct {
	db      *sql.DB
	dialect dialect.Dialect
	refTable *schema.Schema
	clause clause.Clause
	sql     strings.Builder
	sqlVars []interface{}
}
```

里面用到的底层类后面也会一一展开介绍。

## 这一节先到这里