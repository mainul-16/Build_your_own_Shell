#!/usr/bin/env sh

while true; do
  printf "$ "
  read -r line || exit 0

  # exit builtin
  if [ "$line" = "exit" ]; then
    exit 0
  fi

  # echo builtin: print everything after the first space
  case "$line" in
    echo\ *)
      printf "%s\n" "${line#* }"
      ;;
    echo)
      printf "\n"
      ;;
    *)
      [ -z "$line" ] && continue
      printf "%s: command not found\n" "${line%% *}"
      ;;
  esac
done
