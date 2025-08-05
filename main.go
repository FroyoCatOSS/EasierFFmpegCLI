package main

import (
	"fmt"

	ParseConfigJson "github.com/FroyoCatOSS/EasierFFmpegCLI/ParseConfigJson"
)

func main() {
	fmt.Println("Easier FFmpeg CLI is Now Running...")
	ParseConfigJson.ParseJson("configs/DownloadConfig.json")

}
