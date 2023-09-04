package uuid

import (
	"fmt"
	"time"

	"github.com/sony/sonyflake"
)

var Sonyflake *sonyflake.Sonyflake

func InitSonyflake(nodeId uint16) {
	s, err := sonyflake.New(sonyflake.Settings{
		StartTime: time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
		MachineID: func() (uint16, error) {
			return nodeId, nil
		},
	})
	if err != nil {
		panic(fmt.Errorf("error creating sonyflake: %w", err))
	}
	Sonyflake = s
}
