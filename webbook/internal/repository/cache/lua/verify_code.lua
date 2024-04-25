-- 存储在redis中的key ：phone_code:biz:phone
local key = KEYS[1]
-- 要验证的的验证码
local expectedCode = ARGV[1]
-- 验证次数
local cntKey = key .. ":cnt"
local cnt = tonumber(redis.call("get", cntKey))
-- 没有验证次数了，验证码失效
if cntKey <= 0 then
    return -1
end
-- 否则比较存储的验证码
local code = redis.call("get",key)
if code == expectedCode then
    -- 验证成功，将验证码置为无效，并返回0
    redis.call("set",cntKey,0)
    return 0
else
    -- 验证码不相等
    redis.call("decr",cntKey)
    return -2
end
