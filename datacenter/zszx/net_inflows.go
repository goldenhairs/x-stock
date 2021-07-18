// 获取个股指定时间段内资金净流入数据

package zszx

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"go.uber.org/zap"
)

// NetInflow 资金净流入详情
type NetInflow struct {
	// 交易日期
	TrdDt string `json:"TrdDt"`
	// 当日股价
	ClsPrc string `json:"ClsPrc"`
	// 主力净流入（万元）
	MainMnyNetIn string `json:"MainMnyNetIn"`
	// 超大单净流入（万元）
	HugeNetIn string `json:"HugeNetIn"`
	// 大单净流入（万元）
	BigNetIn string `json:"BigNetIn"`
	// 中单净流入（万元）
	MidNetIn string `json:"MidNetIn"`
	// 小单净流入（万元）
	SmallNetIn  string `json:"SmallNetIn"`
	TTLMnyNetIn string `json:"TtlMnyNetIn"`
}

// NetInflowList 净流入详情列表
type NetInflowList []NetInflow

// SumMainNetIn 主力净流入列表求和
func (n NetInflowList) SumMainNetIn(ctx context.Context) float64 {
	var netFlowin float64 = 0.0
	for _, i := range n {
		mainNetIn, err := strconv.ParseFloat(i.MainMnyNetIn, 64)
		if err != nil {
			logging.Errorf(ctx, "Parse MainMnyNetIn:%v to Float error:%v", i.MainMnyNetIn, err)
		}
		netFlowin += mainNetIn
	}
	return netFlowin
}

// RespMainMoneyNetInflows QueryMainMoneyNetInflows 返回json结构
type RespMainMoneyNetInflows struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Code    int           `json:"code"`
	Data    NetInflowList `json:"data"`
}

// QueryMainMoneyNetInflows 查询主力资金净流入数据
func (z Zszx) QueryMainMoneyNetInflows(ctx context.Context, secuCode, startDate, endDate string) (NetInflowList, error) {
	apiurl := "https://zszx.cmschina.com/pcnews/f10/stkcnmnyflow"
	stockCodeAndMarket := strings.Split(secuCode, ".")
	if len(stockCodeAndMarket) != 2 {
		return nil, errors.New("invalid secuCode:" + secuCode)
	}
	stockCode := stockCodeAndMarket[0]
	market := stockCodeAndMarket[1]
	marketCode := "0"
	if strings.ToUpper(market) == "SH" {
		marketCode = "1"
	}
	params := map[string]string{
		"dateStart": startDate,
		"dateEnd":   endDate,
		"ecode":     marketCode,
		"scode":     stockCode,
	}
	logging.Debug(ctx, "Zszx QueryMainMoneyNetInflows "+apiurl+" begin", zap.Any("params", params))
	beginTime := time.Now()
	apiurl, err := goutils.NewHTTPGetURLWithQueryString(ctx, apiurl, params)
	if err != nil {
		return nil, err
	}
	resp := RespMainMoneyNetInflows{}
	err = goutils.HTTPGET(ctx, z.HTTPClient, apiurl, &resp)
	latency := time.Now().Sub(beginTime).Milliseconds()
	logging.Debug(
		ctx,
		"Zszx QueryMainMoneyNetInflows "+apiurl+" end",
		zap.Int64("latency(ms)", latency),
		zap.Any("resp", resp),
	)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("%s %#v", secuCode, resp)
	}
	return resp.Data, nil
}
