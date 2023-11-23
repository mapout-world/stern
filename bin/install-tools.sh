set -e

PROTOC_VERSION=25.1
PROTOC_ZIP=protoc-${PROTOC_VERSION}-osx-x86_64.zip
curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/${PROTOC_ZIP}
unzip -q -o $PROTOC_ZIP -d $HOME/.protobuf
rm -f $PROTOC_ZIP
echo "✅ protoc"

go install $(go list -e -f '{{ join .Imports "\n" }}' tools/tools.go)
echo "✅ tools.go"