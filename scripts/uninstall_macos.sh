#!/bin/bash

# Define argument flags
while getopts 'n:' flag; do
  case "${flag}" in
    n) FILE_NAME=$OPTARG ;;
  esac
done

BASH_COMPLETION_DIR=$(brew --prefix)/etc/bash_completion.d
CURRENT_DIR=$PWD

# Find the GOPATH
GO_BIN_DIR=${HOME}/go/bin
if [[ -z "${GOPATH}" ]]; then
  GO_BIN_DIR=${HOME}/go/bin
else
  GO_BIN_DIR=${GOPATH}/bin
fi

# Use default file name if not set
if [ -z $FILE_NAME ]
then
  FILE_NAME=incognito-cli
fi

# Check if the binary file has been installed before. If yes, remove it.
if test -f "${GO_BIN_DIR}/${FILE_NAME}";
then
  echo "Removing the binary named '${FILE_NAME}' from '${GO_BIN_DIR}'..."
  rm ${GO_BIN_DIR}/${FILE_NAME}
  echo "Binary removed."
fi

# Check if auto completion is set
if test -f "${BASH_COMPLETION_DIR}/${FILE_NAME}";
  then
    echo "Removing the old bash completion file..."
    sudo rm ${BASH_COMPLETION_DIR}/${FILE_NAME}
    echo "Old bash completion file removed."
  fi
