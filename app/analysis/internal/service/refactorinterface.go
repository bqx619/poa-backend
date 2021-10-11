package service

import (
	"context"
	"errors"
	"poa-service/api/analysis/v1"
	"poa-service/app/analysis/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type RefactorInterfaceService struct {
	v1.UnimplementedRefactorInterfaceServer

	ref *biz.RefactorUseCase
	log *log.Helper
}

func NewRefactorInterfaceService(ref *biz.RefactorUseCase, logger log.Logger) *RefactorInterfaceService {
	return &RefactorInterfaceService{
		ref: ref,
		log: log.NewHelper(log.With(logger, "module", "service/refactor")),
	}
}

func (s *RefactorInterfaceService) RecapRefactor(ctx context.Context, req *v1.RecapRefactorRequest) (*v1.RecapRefactorReply, error) {
	status, err := s.ref.Status(ctx)
	if err != nil {
		return &v1.RecapRefactorReply{Status: v1.Status_STOP}, err
	}
	if len(status) != 4 {
		return &v1.RecapRefactorReply{Status: v1.Status_STOP}, errors.New("status data error")
	}
	if int32(status[1]) == 0 {
		return &v1.RecapRefactorReply{
			Status:  v1.Status_HANGING,
			Running: int32(status[1]),
			Success: int32(status[3]),
			Failed:  int32(status[2]),
			Waiting: int32(status[0]),
		}, nil
	}
	return &v1.RecapRefactorReply{
		Status:  v1.Status_RUNNING,
		Running: int32(status[1]),
		Success: int32(status[3]),
		Failed:  int32(status[2]),
		Waiting: int32(status[0]),
	}, nil
}

func (s *RefactorInterfaceService) GetRefactor(ctx context.Context, req *v1.GetRefactorRequest) (*v1.GetRefactorReply, error) {
	zeekLog, err := s.ref.Get(ctx, int(req.GetId()))
	if err != nil {
		return nil, err
	}
	return &v1.GetRefactorReply{
		Info: &v1.RefactorRefactorInfo{
			Id:         int64(zeekLog.ID),
			Dir:        zeekLog.Dir,
			Day:        int64(zeekLog.LogDay),
			Type:       zeekLog.Type.String(),
			BeginHour:  int32(zeekLog.BeginHour),
			EndHour:    int32(zeekLog.EndHour),
			Md5:        zeekLog.Md5,
			CreateTime: zeekLog.CreateTime,
			UpdateTime: zeekLog.UpdateTime,
		},
	}, nil
}

func (s *RefactorInterfaceService) ListRefactor(ctx context.Context, req *v1.ListRefactorRequest) (*v1.ListRefactorReply, error) {
	zeekLogs, count, err := s.ref.List(ctx, req.Sort, req.SortType, int(req.PageSize), int(req.Current), int(req.Id), req.Dir, req.Day)
	if err != nil {
		return nil, err
	}
	infos := make([]*v1.RefactorRefactorInfo, 0)
	for i := 0; i < len(zeekLogs); i++ {
		infos = append(infos, &v1.RefactorRefactorInfo{
			Id:         int64(zeekLogs[i].ID),
			Dir:        zeekLogs[i].Dir,
			Day:        int64(zeekLogs[i].LogDay),
			Type:       zeekLogs[i].Type.String(),
			BeginHour:  int32(zeekLogs[i].BeginHour),
			EndHour:    int32(zeekLogs[i].EndHour),
			Md5:        zeekLogs[i].Md5,
			CreateTime: zeekLogs[i].CreateTime,
			UpdateTime: zeekLogs[i].UpdateTime,
		})
	}
	return &v1.ListRefactorReply{
		Infos: infos,
		Total: int64(count),
	}, nil
}

func (s *RefactorInterfaceService) RunRefactor(ctx context.Context, req *v1.RunRefactorRequest) (*v1.RunRefactorReply, error) {
	s.log.Info("exec run refactor ...")
	go func() {
		err := s.ref.Run(context.TODO())
		//err := s.ref.SyncRun(ctx)
		if err != nil {
			s.log.Error(err)
			return
		}
	}()
	return &v1.RunRefactorReply{}, nil
}

