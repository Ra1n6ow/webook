-- 验证码的key
local key = KEYS[1]
-- 验证码的计数器key
local cntKey = key..":cnt"
-- 你准备的存储的验证码
local val = ARGV[1]

-- 获取验证码的过期时间
local ttl = tonumber(redis.call("ttl", key))
-- key 存在，但是没有过期时间
if ttl == -1 then
    return -2
-- key 不存在，或者过期时间小于 540 秒
elseif ttl == -2 or ttl < 540 then
    --    可以发验证码
    redis.call("set", key, val)
    -- 600 秒
    redis.call("expire", key, 600)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
    return 0
else
    -- 发送太频繁
    return -1
end