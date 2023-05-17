package proto

//go:generate sh -c "mkdir -p ../pkg"
//go:generate protoc --go_out=../pkg --go_opt=paths=source_relative order/order_kafka_message.proto
