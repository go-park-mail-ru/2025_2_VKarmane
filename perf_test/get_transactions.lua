wrk.method = "GET"

local MIN_ID = 80000
local MAX_ID = 150000

local counter = 0
local MAX_REQUESTS = 25000

wrk.headers["Content-Type"] = "application/json"
wrk.headers["Cookie"] =
  "auth_token=; csrf_token=>
wrk.headers["X-CSRF-Token"] =
  ""

request = function()
  if counter >= MAX_REQUESTS then
    return nil 
  end
  counter = counter + 1
  local id = math.random(MIN_ID, MAX_ID)
  local path = "/api/v1/operations/account/17/operation/" .. id
  return wrk.format(nil, path)
end

response = function(status, headers, body)
  if status ~= 200 then
    print("STATUS:", status)
  end
end

