package biz

import (
	"context"
	"poa-service/app/analysis/internal/data/ent"
	"poa-service/app/analysis/internal/data/ent/opinion"
	"poa-service/app/analysis/pkg"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

type Top10Category struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
}

type Top10Area struct {
	Area  string `json:"area"`
	Count int    `json:"count"`
}

type Top10Keyword struct {
	Keyword string `json:"keyword"`
	Count   int    `json:"count"`
}

type EmotionDistribute struct {
	Positive int `json:"position"`
	Neutral  int `json:"neutral"`
	Negative int `json:"negative"`
}

type AnalysisInfo struct {
	OPData   *ent.Opinion `json:"op_data,omitempty"`
	Area     string       `json:"area,omitempty" json:"area,omitempty"`
	Category string       `json:"category,omitempty" json:"category,omitempty"`
	Emotion  int8         `json:"emotion,omitempty" json:"emotion,omitempty"`
	Keywords []string     `json:"keywords,omitempty" json:"keywords,omitempty"`
}

type AnalysisRepo interface {
	GetOpinion(ctx context.Context) ([]*ent.Opinion, error)
	UpdateAnalysis(ctx context.Context, id int, status opinion.AnalysisStatus) error
	CreateArea(ctx context.Context, opId int, area string) error
	CreateCategoryRecord(ctx context.Context, opId int, category string) error
	CreateKeywordsRecord(ctx context.Context, opId int, keywords []string) error
	CreateEmotionRecord(ctx context.Context, opId int, emotion int8) error
	GetTop10Category(ctx context.Context) (top10 []*Top10Category, err error)
	GetTop10Area(ctx context.Context) (top10 []*Top10Area, err error)
	GetTop10Keywords(ctx context.Context) (top10 []*Top10Keyword, err error)
	GetEmotionDistribute(ctx context.Context) (dis *EmotionDistribute, err error)
	GetOpinionByID(ctx context.Context, opID int) (*ent.Opinion, error)
	ListOpinion(ctx context.Context, sort string, sortType string, pageSize int, current int, domain, uid string) ([]*ent.Opinion, int, error)
	GetAnalysisInfo(ctx context.Context, opID int) (*AnalysisInfo, error)
}

type AnalysisUseCase struct {
	repo AnalysisRepo
	ext  pkg.ExternalAnalysisRepo
	log  *log.Helper
}

func NewAnalysisUseCase(analysisRepo AnalysisRepo, extAnalysisRepo pkg.ExternalAnalysisRepo, logger log.Logger) *AnalysisUseCase {
	return &AnalysisUseCase{
		repo: analysisRepo,
		ext:  extAnalysisRepo,
		log:  log.NewHelper(logger),
	}
}

func (ana *AnalysisUseCase) GetAnalysis(ctx context.Context, opId int) (*AnalysisInfo, error) {
	op, err := ana.repo.GetOpinionByID(ctx, opId)
	if err != nil {
		return nil, err
	}
	info, err := ana.repo.GetAnalysisInfo(ctx, opId)
	if err != nil {
		return nil, err
	}
	info.OPData = op
	return info, nil
}

func (ana *AnalysisUseCase) ListAnalysis(ctx context.Context, sort string, sortType string, pageSize int, current int, domain, uid string) ([]*AnalysisInfo, int, error) {
	ops, count, err := ana.repo.ListOpinion(ctx, sort, sortType, pageSize, current, domain, uid)
	if err != nil {
		return nil, 0, err
	}
	infos := make([]*AnalysisInfo, 0)
	for _, op := range ops {
		info, err := ana.repo.GetAnalysisInfo(ctx, op.ID)
		if err != nil {
			return nil, 0, err
		}
		info.OPData = op
		infos = append(infos, info)
	}

	return infos, count, nil
}

func (ana *AnalysisUseCase) GetTop10Category(ctx context.Context) (top10 []*Top10Category, err error) {
	return ana.repo.GetTop10Category(ctx)
}

func (ana *AnalysisUseCase) GetTop10Area(ctx context.Context) (top10 []*Top10Area, err error) {
	return ana.repo.GetTop10Area(ctx)
}

func (ana *AnalysisUseCase) GetTop10Keywords(ctx context.Context) (top10 []*Top10Keyword, err error) {
	return ana.repo.GetTop10Keywords(ctx)
}

func (ana *AnalysisUseCase) GetEmotionDistribute(ctx context.Context) (dis *EmotionDistribute, err error) {
	return ana.repo.GetEmotionDistribute(ctx)
}

func (ana *AnalysisUseCase) Detail(ctx context.Context) {

	ops, err := ana.repo.GetOpinion(context.TODO())
	if err != nil {
		ana.log.Errorf("get op err=%+v", err)
		return
	}
	for _, op := range ops {
		category, err := ana.ext.GetCategory(ctx, op.Title, op.Content)
		if err != nil {
			ana.log.Errorf("get category err=%+v", err)
			continue
		}
		emotion, err := ana.ext.GetEmotion(ctx, op.Content)
		if err != nil {
			ana.log.Errorf("get emotion err=%+v", err)
			continue
		}
		keywords, err := ana.ext.GetKeywords(ctx, op.Title, op.Content)
		if err != nil {
			ana.log.Errorf("get keywords err=%+v", err)
			continue
		}

		err = ana.repo.CreateArea(context.TODO(), op.ID, strings.Split(op.OrigH, ".")[2])
		if err != nil {
			ana.log.Errorf("create area err=%+v", err)
			continue
		}
		err = ana.repo.CreateCategoryRecord(context.TODO(), op.ID, category)
		if err != nil {
			ana.log.Errorf("create category err=%+v", err)
			continue
		}
		err = ana.repo.CreateEmotionRecord(context.TODO(), op.ID, int8(emotion))
		if err != nil {
			ana.log.Errorf("create emotion err=%+v", err)
			return
		}
		err = ana.repo.CreateKeywordsRecord(context.TODO(), op.ID, keywords)
		if err != nil {
			ana.log.Errorf("create keywords err=%+v", err)
			return
		}

		err = ana.repo.UpdateAnalysis(context.TODO(), op.ID, "SUCCESS")
		if err != nil {
			ana.log.Errorf("update analysis err=%+v", err)
			return
		}
	}
}
