package notifications

import (
	"fmt"
	configmodel "github.com/jackpf/kraken-scheduler/src/main/config/model"
)

func NewErrorNotification(schedule configmodel.Schedule, err error) ErrorNotification {
	return ErrorNotification{schedule: schedule, err: err}
}

type ErrorNotification struct {
	schedule configmodel.Schedule
	err      error
}

func (n ErrorNotification) Subject() string {
	return fmt.Sprintf("kraken-scheduler: error")
}

func (n ErrorNotification) Body() string {
	return fmt.Sprintf(`An error occured attempting to order %f %s.

Error: %s`,
		n.schedule.Amount,
		n.schedule.Pair.Name(),
		n.err.Error())

}
