@test "emacs installation" {
  run ./insco emacs

  [ "$status" -eq 0 ]
}
