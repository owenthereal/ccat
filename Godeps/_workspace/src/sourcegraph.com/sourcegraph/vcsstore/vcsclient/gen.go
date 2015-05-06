package vcsclient

//go:generate protoc -I../../../.. -I../../../../github.com/gogo/protobuf/protobuf -I. --gogo_out=. vcsclient.proto
//go:generate sed -i "s#\tTreeEntryType_#\t#g" vcsclient.pb.go
