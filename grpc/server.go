package grpc

import "context"

type UserServer struct {
	UnimplementedUserServiceServer // 当API新加方法却不一定实现时组合
}

func (us UserServer) GetById(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error) {
	return &GetByIdResp{User: &User{Id: 111, Name: "AAA"}}, nil
}
