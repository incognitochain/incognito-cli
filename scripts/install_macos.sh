#!/bin/bash

# Define argument flags
AUTO_COMPLETE='false'
while getopts 'an:' flag; do
  case "${flag}" in
    n) FILE_NAME=$OPTARG ;;
    a) AUTO_COMPLETE='true' ;;
  esac
done

BASH_COMPLETION_DIR=$(brew --prefix)/bash_completion.d
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
echo "Installing the binary named '${FILE_NAME}' into '${GO_BIN_DIR}'..."

# Check if the binary file has been installed before. If yes, remove it.
if test -f "${GO_BIN_DIR}/${FILE_NAME}";
then
  echo "Removing the previously-installed binary..."
  rm ${GO_BIN_DIR}/${FILE_NAME}
  echo "Binary removed."
fi

# Install the binary
go build -o $GO_BIN_DIR/$FILE_NAME
echo "Binary installed."

# Check if auto completion is set
if [ $AUTO_COMPLETE = "true" ]
then
  if test -f "${BASH_COMPLETION_DIR}/${FILE_NAME}";
  then
    echo "Removing the old bash completion file..."
    sudo rm ${BASH_COMPLETION_DIR}/${FILE_NAME}
    echo "Old bash completion file removed."
  fi
  echo "Enabling bash completion..."
  sudo cp $CURRENT_DIR/scripts/autocomplete/bash_autocomplete $BASH_COMPLETION_DIR/$FILE_NAME
  source $BASH_COMPLETION_DIR/$FILE_NAME
  echo "Bash completion enabled."
fi
