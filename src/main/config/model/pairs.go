package model

import (
	"encoding/json"
	"fmt"
)

type Asset struct {
	Name   string
	Symbol string
}

type Pair struct {
	First  Asset
	Second Asset
}

func (p *Pair) Name() string {
	return p.First.Name + p.Second.Name
}

func (p *Pair) UnmarshalJSON(data []byte) error {
	var pairStr string
	if err := json.Unmarshal(data, &pairStr); err != nil {
		return err
	}

	pair, ok := Pairs[pairStr]
	if !ok {
		return fmt.Errorf("pair %s is not valid", pairStr)
	}

	*p = pair
	return nil
}

var (
	ADA  = Asset{"ADA", "₳"}
	ETH  = Asset{"ETH", "Ξ"}
	XETH = Asset{"XETH", "Ξ"}
	XBT  = Asset{"XBT", "₿"}
	XXBT = Asset{"XXBT", "₿"}
	BCH  = Asset{"BCH", "BCH"}
	DASH = Asset{"DASH", "DASH"}
	EOS  = Asset{"EOS", "EOS"}
	GNO  = Asset{"GNO", "GNO"}
	LINK = Asset{"LINK", "LINK"}
	QTUM = Asset{"QTUM", "QTUM"}
	AAVE = Asset{"AAVE", "AAVE"}
	USDT = Asset{"USDT", "USDT"}
	XETC = Asset{"XETC", "ETC"}
	XICN = Asset{"XICN", "ICN"}
	XLTC = Asset{"XLTC", "Ł"}
	XMLN = Asset{"XMLN", "MLN"}
	XREP = Asset{"XREP", "REP"}
	XTZ  = Asset{"XTZ", "TZ"}
	XXDG = Asset{"XXDG", "DOGE"}
	XXLM = Asset{"XXLM", "XLM"}
	XXMR = Asset{"XXMR", "XMR"}
	XXRP = Asset{"XXRP", "XRP"}
	XZEC = Asset{"XZEC", "ZEC"}

	EUR  = Asset{"EUR", "€"}
	ZEUR = Asset{"ZEUR", "€"}
	USD  = Asset{"USD", "$"}
	ZUSD = Asset{"ZUSD", "$"}
	CAD  = Asset{"CAD", "$"}
	ZCAD = Asset{"ZCAD", "$"}
	GBP  = Asset{"GBP", "£"}
	ZGBP = Asset{"GBP", "£"}
	ZJPY = Asset{"ZJPY", "¥"}
)

var Pairs = map[string]Pair{
	"ADACAD":   {ADA, CAD},
	"ADAETH":   {ADA, ETH},
	"ADAEUR":   {ADA, EUR},
	"ADAUSD":   {ADA, USD},
	"ADAXBT":   {ADA, XBT},
	"AAVEUSD":  {AAVE, USD},
	"BCHEUR":   {BCH, EUR},
	"BCHUSD":   {BCH, USD},
	"BCHXBT":   {BCH, XBT},
	"DASHEUR":  {DASH, EUR},
	"DASHUSD":  {DASH, USD},
	"DASHXBT":  {DASH, XBT},
	"EOSETH":   {EOS, ETH},
	"EOSEUR":   {EOS, EUR},
	"EOSUSD":   {EOS, USD},
	"EOSXBT":   {EOS, XBT},
	"GNOETH":   {GNO, ETH},
	"GNOEUR":   {GNO, EUR},
	"GNOUSD":   {GNO, USD},
	"GNOXBT":   {GNO, XBT},
	"LINKUSD":  {LINK, USD},
	"LINKXBT":  {LINK, XBT},
	"QTUMCAD":  {QTUM, CAD},
	"QTUMETH":  {QTUM, ETH},
	"QTUMEUR":  {QTUM, EUR},
	"QTUMUSD":  {QTUM, USD},
	"QTUMXBT":  {QTUM, XBT},
	"USDTZUSD": {USDT, ZUSD},
	"XBTUSDT":  {XBT, USDT},
	"XETCXETH": {XETC, XETH},
	"XETCXXBT": {XETC, XXBT},
	"XETCZEUR": {XETC, ZEUR},
	"XETCZUSD": {XETC, ZUSD},
	"XETHXXBT": {XETH, XXBT},
	"XETHZCAD": {XETH, ZCAD},
	"XETHZEUR": {XETH, ZEUR},
	"XETHZGBP": {XETH, ZGBP},
	"XETHZJPY": {XETH, ZJPY},
	"XETHZUSD": {XETH, ZUSD},
	"XICNXETH": {XICN, ETH},
	"XICNXXBT": {XICN, XXBT},
	"XLTCXXBT": {XLTC, XXBT},
	"XLTCZEUR": {XLTC, ZEUR},
	"XLTCZUSD": {XLTC, ZUSD},
	"XMLNXETH": {XMLN, XETH},
	"XMLNXXBT": {XMLN, XXBT},
	"XREPXETH": {XREP, XETH},
	"XREPXXBT": {XREP, XXBT},
	"XREPZEUR": {XREP, ZEUR},
	"XREPZUSD": {XREP, ZUSD},
	"XTZCAD":   {XTZ, CAD},
	"XTZETH":   {XTZ, ETH},
	"XTZEUR":   {XTZ, EUR},
	"XTZUSD":   {XTZ, USD},
	"XTZXBT":   {XTZ, XBT},
	"XXBTZCAD": {XXBT, ZCAD},
	"XXBTZEUR": {XXBT, ZEUR},
	"XXBTZGBP": {XXBT, ZGBP},
	"XXBTZJPY": {XXBT, ZJPY},
	"XXBTZUSD": {XXBT, ZUSD},
	"XXDGXXBT": {XXDG, XXBT},
	"XXLMXXBT": {XXLM, XXBT},
	"XXLMZEUR": {XXLM, ZEUR},
	"XXLMZUSD": {XXLM, ZUSD},
	"XXMRXXBT": {XXMR, XXBT},
	"XXMRZEUR": {XXMR, ZEUR},
	"XXMRZUSD": {XXMR, ZUSD},
	"XXRPXXBT": {XXRP, XXBT},
	"XXRPZCAD": {XXRP, ZCAD},
	"XXRPZEUR": {XXRP, ZEUR},
	"XXRPZJPY": {XXRP, ZJPY},
	"XXRPZUSD": {XXRP, ZUSD},
	"XZECXXBT": {XZEC, XXBT},
	"XZECZEUR": {XZEC, ZEUR},
	"XZECZUSD": {XZEC, ZUSD},
}
