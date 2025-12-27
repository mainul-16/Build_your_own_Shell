#!/usr/bin/env bash

while true; do
  printf "$ "
  read -r cmd args || exit 0

  if [ "$cmd" = "exit" ]; then
    exit 0

  elif [ "$cmd" = "echo" ]; then
    printf "%s\n" "$args"

  elif [ -n "$cmd" ]; then
    printf "%s: command not found\n" "$cmd"
  fi
done
