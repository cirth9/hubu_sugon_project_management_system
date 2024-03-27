package errs

import (
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	common "test.com/project-common"
)

func GrpcError(err *BError) error {
	return status.Error(codes.Code(err.Code), err.Msg)
}

func ParseGrpcError(err error) (common.BusinessCode, string) {
	fromError, _ := status.FromError(err)
	return common.BusinessCode(fromError.Code()), fromError.Message()
}

func ToBError(err error) *BError {
	fromError, _ := status.FromError(err)
	return NewError(ErrorCode(fromError.Code()), fromError.Message())
}
