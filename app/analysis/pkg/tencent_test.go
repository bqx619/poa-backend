package pkg

import (
	"fmt"
	"testing"
)

func TestTencentNLP(t *testing.T) {
	nlp := NewTencentNLP("AKIDPaVAZhY9TWHlzxe4E7OO8bniMkrgEWXH", "UifKjrjJcePs4zd06uRihPq1vYuznKZX")
	category, err := nlp.GetCategory("欧洲冠军联赛是欧洲足球协会联盟主办的年度足球比赛，代表欧洲俱乐部足球最高荣誉和水平，被认为是全世界最高素质、最具影响力以及最高水平的俱乐部赛事，亦是世界上奖金最高的足球赛事和体育赛事之一。")
	if err != nil {
		return
	}
	fmt.Println(category)
	emotion, err := nlp.GetEmotion("苹果是一家伟大的公司")
	if err != nil {
		return
	}
	fmt.Println(emotion)
	keywords, err := nlp.GetKeywords("苹果是一家伟大的公司")
	if err != nil {
		return
	}
	fmt.Println(keywords)
}