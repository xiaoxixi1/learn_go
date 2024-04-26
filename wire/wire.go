//go:build wireinject

package wire

import (
	"github.com/google/wire"
	"project_go/wire/repository"
	"project_go/wire/repository/dao"
)

/**
  依赖注入的写法：不关心依赖是如何构造的
   缺点：在系统初始化过程中（main函数）完成复杂冗长的构造链表
   解决方案：借助中间件wire
  依赖注入是IOC控制反转的一种实现形式。
  控制反转：UserHandler不会去初始化UserService，外面直接传入userService给UserHandler(依赖注入)
  非依赖注入：必须自己初始化依赖，比如repository需要知道如何初始化DAO和cache
  缺点：
   1 深度耦合依赖的初始化过程
   2 往往需要定义额外的config类型来传递依赖所需的配置
   3 一旦依赖增加新的配置，或者更改了初始化过程，都要跟着修改
   4 缺乏扩展性
   5 测试不友好
   6 难以服用公用组件，例如DB和redis之类的客户端
*/

/**
  wire 分成两部分：
   1 在项目中使用的依赖
   2 命令行工具
  安装命令行工具：
   go install github.com/google/wire/cmd/wire@latest
  使用步骤：
   1 往wire里面注册各种初始化方法（如何初始化）
   2 运行wire命令（开始组装）
  原理：抽象语法树
  但是wire不知单例模式，如果有多此调用initDB，就会生成多个gorm.DB
  缺点：
   1 缺乏根据环境使用不用实现的能力
   2 缺乏根据接口查找实现的能力
   3 缺乏根据类型查找所有实力的能力
  好处：
    很清晰，可控性非常强
*/

func InitUserRepository() *repository.UserRepository {
	wire.Build(repository.NewUserRepository, dao.NewUserDao, intiDb)
	return &repository.UserRepository{}
}
