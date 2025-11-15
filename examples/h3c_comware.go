package examples

import (
	"fmt"
	"log"

	gomiko "github.com/Ali-aqrabawi/gomiko/pkg"
)

func mainH3C() {
	exampleH3CBasic()
}

// 基础示例：连接H3C交换机并执行命令
func exampleH3CBasic() {
	// 创建H3C设备实例（设备类型为h3c_comware）
	device, err := gomiko.NewDevice(
		"192.168.1.2", // IP地址
		"admin",       // 用户名
		"password",    // 密码
		"h3c",         // 设备类型
		22,            // SSH端口
	)
	if err != nil {
		log.Fatalf("创建设备失败: %v", err)
	}

	// 建立连接
	if err := device.Connect(); err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer device.Disconnect()

	// 发送查看命令
	vlanOutput, err := device.SendCommand("display vlan")
	if err != nil {
		log.Fatalf("发送命令失败: %v", err)
	}
	fmt.Println("VLAN信息:")
	fmt.Println(vlanOutput)

	// 发送配置命令
	configCmds := []string{
		"vlan 100",
		"name h3c-test-vlan",
	}
	configOutput, err := device.SendConfigSet(configCmds)
	if err != nil {
		log.Fatalf("发送配置命令失败: %v", err)
	}
	fmt.Println("配置结果:")
	fmt.Println(configOutput)
}
