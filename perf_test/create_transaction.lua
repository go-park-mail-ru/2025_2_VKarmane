wrk.method = "POST"
wrk.path = "/api/v1/operations/account/17"

wrk.headers["Content-Type"] = "application/json"
wrk.headers["Cookie"] =
  "auth_token=; csrf_token=>
wrk.headers["X-CSRF-Token"] =
  ">

wrk.body = [[
{
  "account_id": 17,
  "category_id": 0,
  "date": "2025-12-18T00:00:00.000Z",
  "description": "",
  "name": "Доход",
  "sum": 1,
  "type": "income"
}
]]

response = function(status, headers, body)
  if status ~= 201 and status ~= 200 then
    print("STATUS:", status)
  end
end

