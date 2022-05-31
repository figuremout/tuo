influx query 'from(bucket:"init-bucket") |> range(start:-10m)'
influx query 'from(bucket:"init-bucket") |> range(start:-30m) |> filter(fn: (r) => r.cpu == "cpu0" and r.type == "static")'
from(bucket:"init-bucket") |> range(start:-30m) |> filter(fn: (r) => r.cpu == "cpu0") |> last()
from(bucket:"init-bucket") |> range(start:-1d) |> filter(fn: (r) => r.cpu == "cpu0" and r._field == "time_user") |> keep(columns: ["_field", "_value", "_time"])