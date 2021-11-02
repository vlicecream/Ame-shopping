package sentinel

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"go.uber.org/zap"
)

func Init() error {
	// 先初始化sentinel
	err := sentinel.InitDefault()
	if err != nil {
		zap.L().Error("goods.sentinel.Init sentinel.InitDefault failed", zap.Error(err))
		return err
	}

	// 配置规则
	_, err = flow.LoadRules([]*flow.Rule{
		// 限流
		{
			Resource:               "goodsApi",
			TokenCalculateStrategy: flow.Direct, // 直接使用规则中的Threshold表示当前统计周期内的最大Token数量
			ControlBehavior:        flow.Reject, // 直接拒绝
			Threshold:              3,
			StatIntervalInMs:       6000,
		},
		// 预热与冷启动
		{
			Resource:               "goodsApi-limiting",
			TokenCalculateStrategy: flow.WarmUp, // 冷启动策略
			ControlBehavior:        flow.Reject, // 直接拒绝
			Threshold:              10000,
			WarmUpPeriodSec:        1000,
		},
		// 匀速通过(比如6秒内有三个请求，每2秒放一个请求)
		{
			Resource:               "some-test",
			TokenCalculateStrategy: flow.WarmUp, // 冷启动策略
			ControlBehavior:        flow.Throttling, // 均速通过
			Threshold:              10000,
			WarmUpPeriodSec:        1000,
		},
	})

	if err != nil {
		zap.L().Error("goods.sentinel.Init flow.LoadRules failed", zap.Error(err))
		return err
	}
	return nil
}