protoc --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --go-grpc_opt=require_unimplemented_servers=false --go_out=.. --go-grpc_out=..  --proto_path=. *.proto

# protoc --go_out=plugins=grpc:.  *.proto