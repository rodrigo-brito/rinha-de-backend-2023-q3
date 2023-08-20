//go:generate go run go.uber.org/mock/mockgen -package=mocks -source=$FILE -destination=../../testdata/mocks/grpc.go . Storage
package grpc
