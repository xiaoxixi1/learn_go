//go:build !k8s

package config

/*
*

	因为此处还没有引入配置模块，所以使用一个叫k8s的编译标签来控制不同参数
	然后build的时候使用：
	go build -tags=k8s
*/
var Config = config{
	DB: DbConfig{
		DSN: "root:Jike1504240602*@tcp(localhost:13306)/webook",
	},
	Redis: RedisConfig{
		Addr: "localhost:6379",
	},
}
