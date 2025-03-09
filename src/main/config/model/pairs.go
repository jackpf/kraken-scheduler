package model

import (
	"encoding/json"
	"fmt"
)

type Asset struct {
	Name           string
	NormalisedName string
	Symbol         string
	IsFiat         bool
}

type Pair struct {
	// What are we buying
	First Asset
	// What are we buying it with
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
	ADA  = Asset{"ADA", "ADA", "₳", false}
	ETH  = Asset{"ETH", "XETH", "Ξ", false}
	XETH = Asset{"XETH", "XETH", "Ξ", false}
	XBT  = Asset{"XBT", "XXBT", "₿", false}
	XXBT = Asset{"XXBT", "XXBT", "₿", false}
	BCH  = Asset{"BCH", "BCH", "BCH", false}
	DASH = Asset{"DASH", "DASH", "DASH", false}
	EOS  = Asset{"EOS", "EOS", "EOS", false}
	GNO  = Asset{"GNO", "GNO", "GNO", false}
	LINK = Asset{"LINK", "LINK", "LINK", false}
	QTUM = Asset{"QTUM", "QTUM", "QTUM", false}
	AAVE = Asset{"AAVE", "AAVE", "AAVE", false}
	USDT = Asset{"USDT", "USDT", "USDT", false}
	XETC = Asset{"XETC", "XETC", "ETC", false}
	XICN = Asset{"XICN", "XICN", "ICN", false}
	XLTC = Asset{"XLTC", "XLTC", "Ł", false}
	XMLN = Asset{"XMLN", "XMLN", "MLN", false}
	XREP = Asset{"XREP", "XREP", "REP", false}
	XTZ  = Asset{"XTZ", "XTZ", "TZ", false}
	XXDG = Asset{"XXDG", "XXDG", "DOGE", false}
	XXLM = Asset{"XXLM", "XXLM", "XLM", false}
	XXMR = Asset{"XXMR", "XXMR", "XMR", false}
	XXRP = Asset{"XXRP", "XXRP", "XRP", false}
	XZEC = Asset{"XZEC", "XZEC", "ZEC", false}

	EUR  = Asset{"EUR", "ZEUR", "€", true}
	ZEUR = Asset{"ZEUR", "ZEUR", "€", true}
	USD  = Asset{"USD", "ZUSD", "$", true}
	ZUSD = Asset{"ZUSD", "ZUSD", "$", true}
	CAD  = Asset{"CAD", "ZCAD", "$", true}
	ZCAD = Asset{"ZCAD", "ZCAD", "$", true}
	GBP  = Asset{"GBP", "ZGBP", "£", true}
	ZGBP = Asset{"ZGBP", "ZGBP", "£", true}
	ZJPY = Asset{"ZJPY", "ZJPY", "¥", true}
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
