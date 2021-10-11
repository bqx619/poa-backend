package pkg

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/parnurzeal/gorequest"
)

type BaiduNLP struct {
	client         *gorequest.SuperAgent
	accessTokenURL string
	categoryURL    string
	emotionURL     string
	keywordsURL string
	appID          string
	secretKey      string

	AccessToken string
	ExpiresIn   int64
}

type BaiduNLPAK struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func NewBaiduNLPRepo(appID, secretKey string) *BaiduNLP {
	ins := &BaiduNLP{
		client:         gorequest.New(),
		accessTokenURL: "https://aip.baidubce.com/oauth/2.0/token",
		categoryURL:    "https://aip.baidubce.com/rpc/2.0/nlp/v1/topic",
		emotionURL:     "https://aip.baidubce.com/rpc/2.0/nlp/v1/sentiment_classify",
		keywordsURL: "https://aip.baidubce.com/rpc/2.0/nlp/v1/keyword",
		appID:          appID,
		secretKey:      secretKey,
	}
	ins.getAccessToken()
	return ins
}

func (b *BaiduNLP) getAccessToken() {
	if time.Now().Unix() < b.ExpiresIn {
		return
	}
	var ret BaiduNLPAK
	_, _, errors := b.client.Post(b.accessTokenURL).
		Query(fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", b.appID, b.secretKey)).
		EndStruct(&ret)
	if errors != nil {
		return
	}
	b.AccessToken = ret.AccessToken
	b.ExpiresIn = ret.ExpiresIn
	return
}

func (b *BaiduNLP) GetCategory(title string, content string) (string, error) {
	b.getAccessToken()
	var params = map[string]string{}
	params["title"] = title
	params["content"] = content
	sendBody, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
	}
	type BaiduCategoryResp struct {
		LogId int64 `json:"log_id,omitempty"`
		Item  struct {
			Lv1TagList []struct {
				Score float32 `json:"score,omitempty"`
				Tag   string  `json:"tag,omitempty"`
			} `json:"lv1_tag_list,omitempty"`
		} `json:"item"`
	}
	sendData := string(sendBody)
	var resp BaiduCategoryResp
	_, _, errors := b.client.Post(b.categoryURL).
		Set("Content-Type", "application/json").
		Query(fmt.Sprintf("charset=UTF-8&access_token=%s", b.AccessToken)).
		Send(sendData).EndStruct(&resp)
	if errors != nil {
		return "", errors[0]
	}
	if len(resp.Item.Lv1TagList) > 0 {
		return resp.Item.Lv1TagList[0].Tag, nil
	}
	return "", fmt.Errorf("no lv1 tag")
}

func (b *BaiduNLP) GetEmotion(text string) (EmotionType, error) {
	b.getAccessToken()
	var params = map[string]string{}
	params["text"] = text
	sendBody, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
	}
	type BaiduEmotionResp struct {
		LogId int64 `json:"log_id"`
		Items []struct {
			Sentiment int `json:"sentiment"`
		} `json:"items"`
	}
	sendData := string(sendBody)
	var resp BaiduEmotionResp
	_, _, errors := b.client.Post(b.emotionURL).
		Set("Content-Type", "application/json").
		Query(fmt.Sprintf("charset=UTF-8&access_token=%s", b.AccessToken)).
		Send(sendData).EndStruct(&resp)
	if errors != nil {
		return Emotion_Negative, errors[0]
	}
	if len(resp.Items) == 0 {
		return Emotion_Negative, fmt.Errorf("error args")
	}
	switch resp.Items[0].Sentiment {
	case 0:
		return Emotion_Negative, nil
	case 1:
		return Emotion_Neutral, nil
	case 2:
		return Emotion_Positive, nil
	default:
		return Emotion_Negative, fmt.Errorf("error args")
	}
}

func (b *BaiduNLP) GetKeywords(title, content string) ([]string, error) {
	b.getAccessToken()
	var params = map[string]string{}
	params["title"] = title
	params["content"] = content
	sendBody, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
	}
	type BaiduKeywordsResp struct {
		LogId int64 `json:"log_id,omitempty"`
		Items  []struct {
			Score float32 `json:"score,omitempty"`
			Tag   string  `json:"tag,omitempty"`
		} `json:"items"`
	}
	sendData := string(sendBody)
	var resp BaiduKeywordsResp
	_, _, errors := b.client.Post(b.keywordsURL).
		Set("Content-Type", "application/json").
		Query(fmt.Sprintf("charset=UTF-8&access_token=%s", b.AccessToken)).
		Send(sendData).EndStruct(&resp)
	if errors != nil {
		return nil, errors[0]
	}
	ks := make([]string, 0)
	for _, item := range resp.Items {
		ks = append(ks, item.Tag)
	}
	return ks, nil
}