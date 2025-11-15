package h3c

import (
	"strings"
	"testing"
)

type mockDriver struct {
	ReadSideEffect func() string
	CmdCalls       *string
	PatternCalls   *string
	PromptRegex    *string
	GenericCalls   *string
}

func (c mockDriver) Connect() error {
	return nil
}
func (c mockDriver) Disconnect() {
	*c.GenericCalls = "disconnect"
}
func (c mockDriver) SendCommand(cmd string, expectPattern string) (string, error) {
	if c.CmdCalls != nil {
		*c.CmdCalls += cmd + ", "
	} else {
		c.CmdCalls = &cmd
	}
	return c.ReadUntil(expectPattern)
}
func (c mockDriver) SendCommandsSet(cmds []string, expectPattern string) (string, error) {
	for _, cmd := range cmds {
		_, err := c.SendCommand(cmd, expectPattern)
		if err != nil {
			panic(err)
		}
	}
	return c.ReadUntil(expectPattern)
}
func (c mockDriver) FindDevicePrompt(regex string, pattern string) (string, error) {
	*c.PromptRegex = regex
	return c.ReadUntil(pattern)
}
func (c mockDriver) ReadUntil(pattern string) (string, error) {
	*c.PatternCalls += pattern + ", "
	return c.ReadSideEffect(), nil
}
func (c mockDriver) SetTimeout(timeout uint8) {}

// 测试 H3C Comware Connect
func TestComwareDevice_Connect(t *testing.T) {
	mockD := mockDriver{}
	var cmdCalls, patternCalls, promptRegexCall string
	mockD.CmdCalls = &cmdCalls
	mockD.PatternCalls = &patternCalls
	mockD.PromptRegex = &promptRegexCall

	callsCount := 0
	mockD.ReadSideEffect = func() string {
		callsCount++
		switch callsCount {
		case 1:
			return "<Sysname>"
		case 2:
			return "screen-length disable"
		default:
			return ""
		}
	}

	base := ComwareDevice{mockD, "", ""}
	if err := base.Connect(); err != nil {
		t.Fatal(err)
	}

	if base.Prompt != "<Sysname>" {
		t.Error("Driver.FindDevicePrompt was not called or prompt mismatch")
	}

	// 验证 Pattern 调用顺序
	expected := ">|], <Sysname>, "
	if patternCalls != expected {
		t.Errorf("wrong Comware Pattern calls, Expected: (%s) Got: (%s)", expected, patternCalls)
	}

	// 验证 CLI 初始化命令
	expected = "screen-length disable, "
	if cmdCalls != expected {
		t.Errorf("wrong Comware commands calls, Expected: (%s) Got: (%s)", expected, cmdCalls)
	}

	// 验证提示符匹配正则
	expected = `<(\S+)>|\[(\S+)\]`
	if promptRegexCall != expected {
		t.Errorf("wrong Comware prompt regex calls, Expected: (%s) Got: (%s)", expected, promptRegexCall)
	}
}

// 断开连接测试
func TestComwareDevice_Disconnect(t *testing.T) {
	mockD := mockDriver{}
	var genericCalls string
	mockD.GenericCalls = &genericCalls

	base := ComwareDevice{mockD, "", ""}

	base.Disconnect()

	if genericCalls != "disconnect" {
		t.Error("Driver.Disconnect() was not called")
	}
}

// 发送单个命令测试
func TestComwareDevice_SendCommand(t *testing.T) {
	mockD := mockDriver{}
	var cmdCalls, patternCalls, promptRegexCall string
	mockD.CmdCalls = &cmdCalls
	mockD.PatternCalls = &patternCalls
	mockD.PromptRegex = &promptRegexCall
	mockD.ReadSideEffect = func() string {
		return "display interface brief\n" +
			"Interface 1: UP\n" +
			"Interface 2: DOWN\n" +
			"<Sysname>"
	}

	base := ComwareDevice{mockD, "comware", "<Sysname>"}
	result, _ := base.SendCommand("display interface brief")

	if !strings.Contains(result, "Interface 1: UP") {
		t.Error("wrong result returned, missing Interface 1: UP")
	}
	if !strings.Contains(result, "Interface 2: DOWN") {
		t.Error("wrong result returned, missing Interface 2: DOWN")
	}

	expected := "display interface brief, "
	if cmdCalls != expected {
		t.Errorf("wrong commands calls, Expected: (%s) Got: (%s)", expected, cmdCalls)
	}
}

// 发送配置集测试
func TestComwareDevice_SendConfigSet(t *testing.T) {
	mockD := mockDriver{}
	var cmdCalls, patternCalls, promptRegexCall string
	mockD.CmdCalls = &cmdCalls
	mockD.PatternCalls = &patternCalls
	mockD.PromptRegex = &promptRegexCall
	mockD.ReadSideEffect = func() string {
		return "[system]"
	}

	base := ComwareDevice{mockD, "comware", "<Sysname>"}
	cmds := []string{
		"interface GigabitEthernet 1/0/1",
		"description Uplink to Core",
	}
	_, err := base.SendConfigSet(cmds)
	if err != nil {
		t.Fatal(err)
	}

	expected := "system-view, interface GigabitEthernet 1/0/1, description Uplink to Core, quit, "
	if cmdCalls != expected {
		t.Errorf("wrong commands calls, Expected: (%s) Got: (%s)", expected, cmdCalls)
	}
}
