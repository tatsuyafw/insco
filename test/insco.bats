@test "insco needs a target as an argument" {
  run ./insco

  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "Usage:" ]
}
