#!/bin/bash

SOURCE_FILE_1="/psqlcopy/pg_hba.conf"
TARGET_FILE_1="/var/lib/postgresql/data/pg_hba.conf"
SOURCE_FILE_2="/psqlcopy/postgresql.conf"
TARGET_FILE_2="/var/lib/postgresql/data/postgresql.conf"

if [ -L "$TARGET_FILE_1" ]; then
  echo "Symlink already exists, no action needed."
else
  if [ -e "$TARGET_FILE_1" ]; then
    rm "$TARGET_FILE_1"
    echo "Existing file pg_hba.conf deleted."
  fi

  ln -s "$SOURCE_FILE_1" "$TARGET_FILE_1"
  echo "Symlink for pg_hba.conf created."
fi

if [ -L "$TARGET_FILE_2" ]; then
  echo "Symlink already exists, no action needed."
else
  if [ -e "$TARGET_FILE_2" ]; then
    rm "$TARGET_FILE_2"
    echo "Existing file postgresql.conf deleted."
  fi

  ln -s "$SOURCE_FILE_2" "$TARGET_FILE_2"
  echo "Symlink for postgresql.conf created."
fi