package sonyflake

import (
	"github.com/sony/sonyflake"
	"supreme-flamego/core/logx"
)

var flake *sonyflake.Sonyflake

func init() {
	flake = sonyflake.NewSonyflake(sonyflake.Settings{})
}

func GenSonyFlakeId() (int64, error) {
	id, err := flake.NextID()
	if err != nil {
		logx.NameSpace("sonyFlakeId").Warn("flake NextID failed: ", err)
		return 0, err
	}
	return int64(id), nil
}
