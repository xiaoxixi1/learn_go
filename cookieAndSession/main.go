package main

import (
	"fmt"
	"net/http"
)

/*
*

				   cookie的关键配置:
						    Path       string    表示cookie可以用在什么路径下，按照最小化原则设定
							Domain     string    表示cookie可以用在什么域名下，按照最小化原则设定
							Expires    time.Time // optional
							// MaxAge=0 means no 'Max-Age' attribute specified.
							// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
							// MaxAge>0 means Max-Age attribute present and given in seconds
							MaxAge   int
					        // 过期时间，只保留必要时间
							Secure   bool // 只能用于https协议，生产环境永远设置成true
							HttpOnly bool // 如果设置为True，浏览器上的jS代表将无法使用这个cookie，永远设为true
							SameSite SameSite //是否允许跨站发送cookie，尽量避免
				   cookie是存储在浏览器本地的，所以很不安全，所以一般只存放一些不太关键的数据
			       关键的东西放后端：session,可以使用session来记录登录状态
	               sessionId一般可以放header ,或者查询参数，或者cookie里面
		           通过将sessionId存放在cookie的方式来记录当前用户的登录状态，同时sessionID设置有效期
*/
func main() {
	ck := http.Cookie{}
	fmt.Printf("%v \n", ck)
}
