local key = KEYS[1]

local cntKey = key..":cnt"
local inputCode = ARGV[1]
local code = redis.call("get",key)
local cnt = tonumber(redis.call("get",cntKey))
if cnt == nil or cnt <= 0 then
    return -1
end

if inputCode == code then
    redis.call('set',cntKey,0)
    return 0
else
    redis.call('decr',cntKey)
    return -2
end