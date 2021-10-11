package pkg

import (
	"fmt"
	"testing"
)

func TestAnalysis(t *testing.T) {
	nlp := NewBaiduNLPRepo("jmax2aeel6QbtjLz41KgXcFR", "UFAnh9swfOhexrs8juRd6c51hUXFeF3h")
	nlp.GetCategory("欧洲冠军联赛", "欧洲冠军联赛是欧洲足球协会联盟主办的年度足球比赛，代表欧洲俱乐部足球最高荣誉和水平，被认为是全世界最高素质、最具影响力以及最高水平的俱乐部赛事，亦是世界上奖金最高的足球赛事和体育赛事之一。")
	emotion, err := nlp.GetEmotion("苹果是一家伟大的公司")
	if err != nil {
		return
	}
	fmt.Println(emotion)
	keywords, err := nlp.GetKeywords("iphone手机出现“白苹果”原因及解决办法，用苹果手机的可以看下", "如果下面的方法还是没有解决你的问题建议来我们门店看下成都市锦江区红星路三段99号银石广场24层01室。在通电的情况下掉进清水，这种情况一不需要拆机处理。尽快断电。用力甩干，但别把机器甩掉，主意要把屏幕内的水甩出来。如果屏幕残留有水滴，干后会有痕迹。^H3 放在台灯，射灯等轻微热源下让水分慢慢散去。")
	if err != nil {
		return
	}
	fmt.Println(keywords)
}