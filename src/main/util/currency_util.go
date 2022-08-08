package util

import (
	"fmt"
	"github.com/jackpf/kraken-scheduler/src/main/config/model"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func FormatAsset(asset model.Asset, amount float64) string {
	if asset.IsFiat {
		return fmt.Sprintf("%s%s", asset.Symbol, FormatCurrency(amount))
	} else {
		return fmt.Sprintf("%s%s", FormatFloat(amount, 8), asset.Symbol)
	}
}

func FormatCurrency(amount float64) string {
	return FormatFloat(amount, 2)
}

func FormatFloat(f float64, precision int) string {
	withPrecision, err := strconv.ParseFloat(fmt.Sprintf("%."+strconv.Itoa(precision)+"f", f), 64)
	if err != nil {
		log.Panic(err.Error())
	}

	return strconv.FormatFloat(withPrecision, 'f', -1, 64)
}
