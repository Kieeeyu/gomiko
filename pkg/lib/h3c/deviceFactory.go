package h3c

import (
	"errors"

	"github.com/Ali-aqrabawi/gomiko/pkg/connections"
	"github.com/Ali-aqrabawi/gomiko/pkg/driver"
	"github.com/Ali-aqrabawi/gomiko/pkg/types"
)

// NewDevice 根据连接和设备类型创建H3C设备实例
func NewDevice(connection connections.Connection, deviceType string) (types.Device, error) {
	if deviceType != "h3c" {
		return nil, errors.New("不支持的H3C设备类型: " + deviceType)
	}

	devDriver := driver.NewDriver(connection, "\n")
	return &ComwareDevice{
		Driver:     devDriver,
		DeviceType: deviceType,
		Prompt:     "",
	}, nil
}
