package man

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func Showman() {
	clearCommand := exec.Command("clear")
	clearCommand.Stdout = os.Stdout
	clearCommand.Run()

	filePath := "man.txt"
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("无法打开文件：%s\n", filePath)
		return
	}

	fmt.Println(string(content))
	fmt.Println("按任意键返回……")
	fmt.Scanln()
	fmt.Scanln()
}
