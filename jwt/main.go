package main

/**
 jwt：Json web Token，主要用于身份认证，也就是登录
基本原理就是通过加密生成一个token,而后客户端每次访问的时候都带上这个token
由三部分组成：
   Header:头部，JWT的元数据，也就是描述这个token本身的数据，一个JSON对象
   Payload:负载，数据内容，一个JSON对象
   Signature:签名，根据Header和Payload生成
登录过程中使用JWT也是两个步骤：
   1 加密和解密JWT数据
   2 登录校验
*/
