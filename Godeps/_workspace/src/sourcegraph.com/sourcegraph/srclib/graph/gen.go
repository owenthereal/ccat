package graph

//go:generate protoc --proto_path=/usr/include:$HOME/src:$HOME/src/github.com/gogo/protobuf/protobuf/google/protobuf:../ann:. --gogo_out=. def.proto doc.proto output.proto ref.proto
//go:generate sed -i "s/^import ann .*$//" output.pb.go
//go:generate sed -i "s/sourcegraph_com_sourcegraph_srclib_ann/ann/g" output.pb.go
//go:generate sed -i "s/Data \\[\\]byte/Data json.RawMessage/g" def.pb.go
//go:generate sed -i "s/^package graph$/package graph;import \"encoding\\/json\"/" def.pb.go
