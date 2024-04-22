package main

/**
  Gin  session存储的实现,包括：
   cookie
   gorm
   memcached
   memstore :基于内存的实现
   mongodb
   postgres
   redis
   tester : 基于测试实现
  注意：sess_id肯定是放在cookie里面的，但是Session里面的数据，比如代码中的user_id才是store存储的
*/

/**
  一般单机单实力部署，可以考虑memstore
  多实例部署，应该选择redis实现
*/
