package h3c

import (
	"errors"

	"github.com/Ali-aqrabawi/gomiko/pkg/driver"
)

// ComwareDevice 表示H3C Comware系统设备
type ComwareDevice struct {
	Driver     driver.IDriver // 底层驱动（处理SSH通信）
	DeviceType string         // 设备类型标识（如 "h3c_comware"）
	Prompt     string         // 设备提示符（如 <H3C> 或 [H3C]）
}

// Connect 建立与设备的连接并初始化会话
func (d *ComwareDevice) Connect() error {

	if err := d.Driver.Connect(); err != nil {
		return err
	}

	prompt, err := d.Driver.FindDevicePrompt(`<(\S+)>|\[(\S+)\]`, ">|]")
	if err != nil {
		return err
	}
	d.Prompt = prompt

	return d.sessionPreparation()
}

// Disconnect 断开与设备的连接
func (d *ComwareDevice) Disconnect() {
	d.Driver.Disconnect()
}

// SendCommand 发送单条命令并返回结果
func (d *ComwareDevice) SendCommand(cmd string) (string, error) {
	return d.Driver.SendCommand(cmd, d.Prompt)
}

// SendConfigSet 批量发送配置命令（自动处理配置模式切换）
func (d *ComwareDevice) SendConfigSet(cmds []string) (string, error) {

	configEnterOutput, err := d.Driver.SendCommand("system-view", "\\[.*\\]")
	if err != nil {
		return "", errors.New("进入配置模式失败: " + err.Error())
	}

	cmdsOutput, err := d.Driver.SendCommandsSet(cmds, "\\[.*\\]")
	if err != nil {
		return "", errors.New("执行配置命令失败: " + err.Error())
	}

	exitOutput, err := d.Driver.SendCommand("quit", "<.*>")
	if err != nil {
		return "", errors.New("退出配置模式失败: " + err.Error())
	}

	return configEnterOutput + cmdsOutput + exitOutput, nil
}

// SetTimeout 设置命令超时时间
func (d *ComwareDevice) SetTimeout(timeout uint8) {
	d.Driver.SetTimeout(timeout)
}

// sessionPreparation 初始化会话（禁用分页等）
func (d *ComwareDevice) sessionPreparation() error {
	_, err := d.Driver.SendCommand("screen-length disable", d.Prompt)
	if err != nil {
		return errors.New("禁用分页失败: " + err.Error())
	}
	return nil
}
