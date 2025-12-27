#!/usr/bin/env sh

while true; do
  printf "$ "
  read -r line || exit 0

  # exit builtin
  if [ "$line" = "exit" ]; then
    exit 0
  fi

  # echo builtin (strip leading "echo ")
  case "$line" in
    echo\ *)
      printf "%s\n" "${line#echo }"
      ;;
    *)
      # ignore empty line
      [ -z "$line" ] && continue
      printf "%s: command not found\n" "$line"
      ;;
  esac
done
