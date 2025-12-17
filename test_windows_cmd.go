// +build ignore

package main

import (
	"fmt"
	"os/exec"
)

func main() {
	// 测试不同的命令构造方式
	command := `dir "C:\Program Files (x86)"`
	
	fmt.Println("原始命令:", command)
	fmt.Println()
	
	// 方式1: 直接传递
	fmt.Println("方式1: cmd /C command")
	cmd1 := exec.Command("cmd", "/C", command)
	fmt.Printf("  Args: %#v\n", cmd1.Args)
	fmt.Printf("  String: %s\n\n", cmd1.String())
	
	// 方式2: 使用 /S /C "command"
	fmt.Println("方式2: cmd /S /C \"command\"")
	wrappedCommand := `"` + command + `"`
	cmd2 := exec.Command("cmd", "/S", "/C", wrappedCommand)
	fmt.Printf("  Args: %#v\n", cmd2.Args)
	fmt.Printf("  String: %s\n\n", cmd2.String())
	
	// 方式3: 使用 /c 和完整引号包裹
	fmt.Println("方式3: cmd /c \"command\"")
	cmd3 := exec.Command("cmd", "/c", `"`+command+`"`)
	fmt.Printf("  Args: %#v\n", cmd3.Args)
	fmt.Printf("  String: %s\n\n", cmd3.String())
}
