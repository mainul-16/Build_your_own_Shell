#!/usr/bin/env sh

while true; do
  printf "$ "

  # Read entire line into $REPLY
  read -r || exit 0

  # Skip empty lines
  [ -z "$REPLY" ] && continue

  # Split words from REPLY
  set -- $REPLY

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
