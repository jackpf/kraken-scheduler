package model

import (
	"encoding/json"
	"fmt"
)

type Asset struct {
	Name   string
	Symbol string
	IsFiat bool
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
	ADA  = Asset{"ADA", "₳", false}
	ETH  = Asset{"ETH", "Ξ", false}
	XETH = Asset{"XETH", "Ξ", false}
	XBT  = Asset{"XBT", "₿", false}
	XXBT = Asset{"XXBT", "₿", false}
	BCH  = Asset{"BCH", "BCH", false}
	DASH = Asset{"DASH", "DASH", false}
	EOS  = Asset{"EOS", "EOS", false}
	GNO  = Asset{"GNO", "GNO", false}
	LINK = Asset{"LINK", "LINK", false}
	QTUM = Asset{"QTUM", "QTUM", false}
	AAVE = Asset{"AAVE", "AAVE", false}
	USDT = Asset{"USDT", "USDT", false}
	XETC = Asset{"XETC", "ETC", false}
	XICN = Asset{"XICN", "ICN", false}
	XLTC = Asset{"XLTC", "Ł", false}
	XMLN = Asset{"XMLN", "MLN", false}
	XREP = Asset{"XREP", "REP", false}
	XTZ  = Asset{"XTZ", "TZ", false}
	XXDG = Asset{"XXDG", "DOGE", false}
	XXLM = Asset{"XXLM", "XLM", false}
	XXMR = Asset{"XXMR", "XMR", false}
	XXRP = Asset{"XXRP", "XRP", false}
	XZEC = Asset{"XZEC", "ZEC", false}

	EUR  = Asset{"EUR", "€", true}
	ZEUR = Asset{"ZEUR", "€", true}
	USD  = Asset{"USD", "$", true}
	ZUSD = Asset{"ZUSD", "$", true}
	CAD  = Asset{"CAD", "$", true}
	ZCAD = Asset{"ZCAD", "$", true}
	GBP  = Asset{"GBP", "£", true}
	ZGBP = Asset{"GBP", "£", true}
	ZJPY = Asset{"ZJPY", "¥", true}
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
