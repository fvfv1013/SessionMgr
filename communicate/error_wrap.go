package communicate

import (
	"errors"
	"fmt"
	"sessionmgr"
	"sessionmgr/proto/pkg/error_pb"
)

// ErrorWrap ErrorWrap将Go标准库的error类型转换为error_pb.Error类型
func ErrorWrap(err error) *error_pb.Error {
	var errorProto error_pb.Error
	errorProto.Message = err.Error()

	switch {
	case errors.Is(err, sessionmgr.ErrID):
		errorProto.ErrId = &error_pb.ErrID{
			ID: 0, // 这里假设你可能需要根据具体情况设置合适的ID值，目前先设为0
		}
	case errors.Is(err, sessionmgr.ErrCall):
		errorProto.ErrCall = &error_pb.ErrCall{}
	case errors.Is(err, sessionmgr.ErrLost):
		errorProto.ErrLost = &error_pb.ErrLost{
			ID: 0, // 同样假设先设为0，可根据实际需求调整
		}
	case errors.Is(err, sessionmgr.ErrWait):
		errorProto.ErrWait = &error_pb.ErrWait{}
	case errors.Is(err, sessionmgr.ErrSdp):
		errorProto.ErrSdp = &error_pb.ErrSdp{}
	}

	return &errorProto
}

func ErrorUnwrap(errorProto *error_pb.Error) error {
	if errorProto == nil {
		return nil
	}

	if errorProto.ErrId != nil {
		return sessionmgr.ErrID
	}

	if errorProto.ErrCall != nil {
		return sessionmgr.ErrCall
	}

	if errorProto.ErrLost != nil {
		return sessionmgr.ErrLost
	}

	if errorProto.ErrWait != nil {
		return sessionmgr.ErrWait
	}

	if errorProto.ErrSdp != nil {
		return sessionmgr.ErrSdp
	}

	return fmt.Errorf("unknown error: %s", errorProto.Message)
}
