package pkg

import (
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	
	nlp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/nlp/v20190408"
)

type TencentNLP struct {
	client *nlp.Client

}

func NewTencentNLP(secretId, secretKey string) *TencentNLP {
	credential := common.NewCredential(
		secretId,
		secretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "nlp.tencentcloudapi.com"
	client, _ := nlp.NewClient(credential, "ap-guangzhou", cpf)
	return &TencentNLP{
		client: client,
	}
}

func (t *TencentNLP) GetCategory(text string) (string, error) {
	request := nlp.NewTextClassificationRequest()

	request.Text = common.StringPtr(text)
	request.Flag = common.Uint64Ptr(1)

	response, err := t.client.TextClassification(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return "", err
	}
	if err != nil {
		panic(err)
	}
	if len(response.Response.Classes) == 0 {
		return "", fmt.Errorf("error resp classes")
	}
	return *response.Response.Classes[0].FirstClassName, nil
}

func (t *TencentNLP) GetEmotion(text string) (EmotionType, error) {
	request := nlp.NewSentimentAnalysisRequest()

	request.Text = common.StringPtr(text)
	request.Mode = common.StringPtr("3class")
	response, err := t.client.SentimentAnalysis(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return Emotion_Negative, err
	}
	if err != nil {
		panic(err)
	}
	switch *response.Response.Sentiment {
	case "positive":
		return Emotion_Positive, nil
	case "negative":
		return Emotion_Negative, nil
	case "neutral":
		return Emotion_Neutral, nil
	default:
		return Emotion_Negative, fmt.Errorf("error resp args")
	}
}

func (t *TencentNLP) GetKeywords(text string) ([]string, error) {
	req := nlp.NewKeywordsExtractionRequest()
	req.Text = common.StringPtr(text)
	response, err := t.client.KeywordsExtraction(req)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return nil, err
	}
	ret := make([]string, 0)
	for _, keyword := range response.Response.Keywords {
		ret = append(ret, *keyword.Word)
	}
	return ret, nil
}