#!/usr/bin/env sh

while true; do
  printf "$ "
  read -r line || exit 0

  # Remove carriage return if present (CRLF-safe)
  line=$(printf "%s" "$line" | tr -d '\r')

  # exit builtin
  if [ "$line" = "exit" ]; then
    exit 0
  fi

  # echo builtin (match anything starting with "echo")
  case "$line" in
    echo\ *)
      printf "%s\n" "${line#echo }"
      ;;
    echo)
      printf "\n"
      ;;
    *)
      [ -z "$line" ] && continue
      printf "%s: command not found\n" "$line"
      ;;
  esac
done
