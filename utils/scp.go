package utils

import (
	"fmt"
	"os/exec"
)

// ScpCopy 使用 scp 命令拷贝文件到远程服务器
func ScpCopy(localPath, remotePath, host, user, password string) error {
	// 构建 scp 命令
	cmd := exec.Command("scp",
		"-o", "StrictHostKeyChecking=no", // 跳过主机密钥检查
		"-o", "UserKnownHostsFile=/dev/null", // 不保存主机密钥
		localPath,
		fmt.Sprintf("%s@%s:%s", user, host, remotePath),
	)

	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("scp failed: %v, output: %s", err, string(output))
	}

	return nil
}

// 使用示例
func main() {
	err := ScpCopy(
		"/local/file.txt",
		"/remote/file.txt",
		"192.168.1.100",
		"root",
		"password123",
	)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("File copied successfully")
	}
}
