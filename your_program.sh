#!/usr/bin/env bash

while true; do
  printf "$ "
  read -r line || exit 0

  # Skip empty lines
  [ -z "$line" ] && continue

  # Split input safely into an array
  read -r -a tokens <<< "$line"

  cmd="${tokens[0]}"
  args=("${tokens[@]:1}")

  if [ "$cmd" = "exit" ]; then
    exit 0

  elif [ "$cmd" = "echo" ]; then
    printf "%s\n" "${args[*]}"

  else
    printf "%s: command not found\n" "$cmd"
  fi
done
