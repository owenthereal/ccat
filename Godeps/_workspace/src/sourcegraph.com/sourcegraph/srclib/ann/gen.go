package ann

//go:generate protoc --proto_path=/usr/include:$HOME/src:$HOME/src/github.com/gogo/protobuf/protobuf/google/protobuf:. --gogo_out=. ann.proto
//go:generate sed -i "s/Data \\[\\]byte/Data json.RawMessage/g" ann.pb.go
//go:generate sed -i "s/^package ann$/package ann;import \"encoding\\/json\"/" ann.pb.go
