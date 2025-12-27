#!/bin/sh

while true; do
  printf "$ "
  read -r line || exit 0

  # Ignore empty input
  [ -z "$line" ] && continue

  # Split input into words
  set -- $line

  cmd="$1"
  shift

  if [ "$cmd" = "exit" ]; then
    exit 0

  elif [ "$cmd" = "echo" ]; then
    printf "%s\n" "$*"

  else
    printf "%s: command not found\n" "$cmd"
  fi
done
