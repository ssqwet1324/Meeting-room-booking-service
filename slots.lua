wrk.method = "GET"

local token = "<Bearer token>"

request = function()
  wrk.headers["Authorization"] = "Bearer " .. token
  return wrk.format(
    nil,
    "/rooms/{roomID(UUID)}/slots/list?date={yyy-MM-dd}"
  )
end
