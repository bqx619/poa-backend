package pkg

import (
	"context"
	"fmt"
)

type EmotionType int8

const (
	Emotion_Positive EmotionType = 1
	Emotion_Neutral EmotionType = 0
	Emotion_Negative EmotionType = -1
)

type ExternalAnalysisRepo interface {
	GetCategory(ctx context.Context, title, content string) (string, error)
	GetEmotion(ctx context.Context, content string) (EmotionType, error)
	GetKeywords(ctx context.Context, title, content string) ([]string, error)
}

type AnalysisIns struct {
	baidu *BaiduNLP
	tencent *TencentNLP
}

func NewAnalysisIns(baiduAPPID, baidUSecretKey, tencentSecretID, tencentSecretKey string) ExternalAnalysisRepo {
	return &AnalysisIns{baidu: NewBaiduNLPRepo(baiduAPPID, baidUSecretKey), tencent: NewTencentNLP(tencentSecretID, tencentSecretKey)}
}

func (a *AnalysisIns) GetCategory(ctx context.Context, title, content string) (string, error) {
	if len(content) > 9000 {
		content = content[:9000]
	}
	baiduCategory, err := a.baidu.GetCategory(title, content)
	if err != nil {
		tencentCategory, err2 := a.tencent.GetCategory(title + ", " + content)
		if err2 != nil {
			return "", fmt.Errorf("baidu=%+v, tencent=%+v", err, err2)
		}
		return tencentCategory, nil
	}
	return baiduCategory, nil
}

func (a *AnalysisIns) GetEmotion(ctx context.Context, content string) (EmotionType, error) {
	if len(content) > 190 {
		content = content[:190]
	}
	baiduEmotion, err := a.baidu.GetEmotion(content)
	if err != nil {
		tencentEmotion, err2 := a.tencent.GetEmotion(content)
		if err2 != nil {
			return 0, fmt.Errorf("baidu=%+v, tencent=%+v", err, err2)
		}
		return tencentEmotion, nil
	}
	return baiduEmotion, nil
}

func (a *AnalysisIns) GetKeywords(ctx context.Context, title, content string) ([]string, error) {
	if len(content) > 9000 {
		content = content[:9000]
	}
	baiduKeywords, err := a.baidu.GetKeywords(title, content)
	if err != nil {
		tencentKeywords, err2 := a.tencent.GetKeywords(title + ", " + content)
		if err2 != nil {
			return nil, fmt.Errorf("baidu=%+v, tencent=%+v", err, err2)
		}
		return tencentKeywords, nil
	}
	return baiduKeywords, nil
}

