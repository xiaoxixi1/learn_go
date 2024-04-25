-- 存储在redis中的key ：phone_code:biz:phone
local key = KEYS[1]
-- 要存储的验证码
local val = ARGV[1]
-- 验证次数
local cntKey = key .. ":cnt"
-- 有效期时间
local ttl = tonumber(redis.call("ttl", key))

if ttl == -1 then
    --key 存在，但是没有过期时间
    return -1 --返回给go
elseif ttl == -2 or ttl < 540 then
    -- key不存在或者有效时间小于9分钟，则可以发送验证码
    redis.call("set", key, val)
    redis.call("expire", key, 600) -- 设置600s的有效期
    redis.call("set", cntKey, 3) -- 设置验证次数
    redis.call("expire", cntKey, 600)
else
    -- 发送太频繁
    return -2
end
