#!/bin/bash

set -e

GO_VERSION="1.2"
GO_DOWNLOAD_URL="https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz"

curl -LO "${GO_DOWNLOAD_URL}"
tar -C /usr/local -xzf "go${GO_VERSION}.linux-amd64.tar.gz"
rm "go${GO_VERSION}.linux-amd64.tar.gz"

echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
source ~/.bashrc

echo "Go ${GO_VERSION} installed successfully"
