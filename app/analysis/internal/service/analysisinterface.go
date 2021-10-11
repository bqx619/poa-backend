package service

import (
	"context"
	"fmt"
	"poa-service/app/analysis/internal/biz"

	pb "poa-service/api/analysis/v1"

	"github.com/go-kratos/kratos/v2/log"
)

type AnalysisInterfaceService struct {
	pb.UnimplementedAnalysisInterfaceServer

	ana *biz.AnalysisUseCase
	log *log.Helper
}

func NewAnalysisInterfaceService(ana *biz.AnalysisUseCase, logger log.Logger) *AnalysisInterfaceService {
	return &AnalysisInterfaceService{
		ana: ana,
		log: log.NewHelper(log.With(logger, "module", "service/analysis")),
	}
}

func (s *AnalysisInterfaceService) RunAnalysis(ctx context.Context, req *pb.RunAnalysisRequest) (*pb.RunAnalysisReply, error) {
	go s.ana.Detail(context.TODO())
	return &pb.RunAnalysisReply{}, nil
}
func (s *AnalysisInterfaceService) Top10Category(ctx context.Context, req *pb.Top10CategoryRequest) (*pb.Top10Reply, error) {
	category, err := s.ana.GetTop10Category(ctx)
	if err != nil {
		return nil, err
	}
	top10 := make([]*pb.Top10, 0)
	for _, top10Category := range category {
		top10 = append(top10, &pb.Top10{
			Name:  top10Category.Category,
			Count: int64(top10Category.Count),
		})
	}
	return &pb.Top10Reply{Top10: top10}, nil
}
func (s *AnalysisInterfaceService) Top10Area(ctx context.Context, req *pb.Top10AreaRequest) (*pb.Top10Reply, error) {
	area, err := s.ana.GetTop10Area(ctx)
	if err != nil {
		return nil, err
	}
	top10 := make([]*pb.Top10, 0)
	for _, top10Area := range area {
		top10 = append(top10, &pb.Top10{
			Name:  top10Area.Area,
			Count: int64(top10Area.Count),
		})
	}
	return &pb.Top10Reply{Top10: top10}, nil
}
func (s *AnalysisInterfaceService) Top10Keywords(ctx context.Context, req *pb.Top10KeywordsRequest) (*pb.Top10Reply, error) {
	keywords, err := s.ana.GetTop10Keywords(ctx)
	if err != nil {
		return nil, err
	}
	top10 := make([]*pb.Top10, 0)
	for _, keyword := range keywords {
		top10 = append(top10, &pb.Top10{
			Name:  keyword.Keyword,
			Count: int64(keyword.Count),
		})
	}
	return &pb.Top10Reply{Top10: top10}, nil
}
func (s *AnalysisInterfaceService) EmotionDistribute(ctx context.Context, req *pb.EmotionDistributeRequest) (*pb.EmotionDistributeReply, error) {
	distribute, err := s.ana.GetEmotionDistribute(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.EmotionDistributeReply{
		Positive: int64(distribute.Positive),
		Negative: int64(distribute.Negative),
		Neutral:  int64(distribute.Neutral),
	}, nil
}

func (s *AnalysisInterfaceService) GetAnalysis(ctx context.Context, req *pb.GetAnalysisRequest) (rsp *pb.GetAnalysisReply, err error) {
	analysis, err := s.ana.GetAnalysis(ctx, int(req.GetId()))
	if err != nil {
		return nil, err
	}
	return &pb.GetAnalysisReply{Info: &pb.AnalysisInfo{
		Ts:       int64(analysis.OPData.Ts),
		Uid:      analysis.OPData.UID,
		OrigH:    analysis.OPData.OrigH,
		RespH:    analysis.OPData.RespH,
		Domain:   analysis.OPData.Domain,
		Proto:    analysis.OPData.Proto,
		Area:     analysis.Area,
		Keywords: analysis.Keywords,
		Category: analysis.Category,
		Emotion:  int32(analysis.Emotion),
		OpId:     int64(analysis.OPData.ID),
	}}, nil
}

func (s *AnalysisInterfaceService) ListAnalysis(ctx context.Context, req *pb.ListAnalysisRequest) (rsp *pb.ListAnalysisReply, err error) {
	analysis, count, err := s.ana.ListAnalysis(ctx, req.Sort, req.SortType, int(req.PageSize), int(req.Current), req.Domain, req.Uid)
	if err != nil {
		return nil, err
	}
	infos := make([]*pb.AnalysisInfo, 0)
	for _, a := range analysis {

		infos = append(infos, &pb.AnalysisInfo{
			Ts:       int64(a.OPData.Ts),
			Uid:      a.OPData.UID,
			OrigH:    a.OPData.OrigH,
			RespH:    a.OPData.RespH,
			Domain:   a.OPData.Domain,
			Proto:    a.OPData.Proto,
			Area:     a.Area,
			Keywords: a.Keywords,
			Category: a.Category,
			Emotion:  int32(a.Emotion),
			OpId:     int64(a.OPData.ID),
		})
	}
	return &pb.ListAnalysisReply{Infos: infos, Total: int64(count)}, nil
}

func (s *AnalysisInterfaceService) EmotionNotice(ctx context.Context, req *pb.EmotionNoticeRequest) (*pb.EmotionNoticeReply, error) {
	distribute, err := s.ana.GetEmotionDistribute(ctx)
	if err != nil {
		return nil, err
	}
	if (int64(distribute.Positive) + int64(distribute.Neutral)) > int64(distribute.Negative) {
		return &pb.EmotionNoticeReply{
			Status: "success",
			Msg:    "当前系统无倾向预警, 舆情正常!",
		}, nil
	}
	return &pb.EmotionNoticeReply{
		Status: "fail",
		Msg:    fmt.Sprintf("请注意, 当前负面舆情数量为%d, 超过预警阈值, 请及时引导", int64(distribute.Negative)),
	}, nil
}

