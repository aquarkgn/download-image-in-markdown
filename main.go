// 下载指定目录下的所有markdown文件内的图片到本地source/images目录下
package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	markdownPath = flag.String("markdownPath", "", "markdown文件目录,默认为执行命令的当前目录")
	imagePath    = flag.String("imagePath", "source/image", "图片下载目录,相对于参数:markdownPath地址")
	rewrite      = flag.String("rewrite", "n", "y:覆盖原文件,n:不覆盖原文件,默认为n")
)

func main() {
	log.SetFlags(0)
	defer log.Println("INFO: 下载完成.")

	// 检索CLI标志和执行环境
	flag.Parse()

	if *markdownPath == "" {
		*markdownPath, _ = os.Getwd()
	}

	downloadImageDir := path.Join(*markdownPath, *imagePath)

	// 循环处理目录下的所有markdown文件
	err := filepath.Walk(*markdownPath, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".md" {
			log.Printf("INFO: Processing %s\n", path)
			err := processFile(path, downloadImageDir)
			if err != nil {
				log.Printf("ERROR: %s\n", err)
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("ERROR: %s\n", err)
	}
}

// processFile 处理指定的markdown文件
func processFile(filePath string, downloadImageDir string) error {
	if _, err := os.Stat(downloadImageDir); os.IsNotExist(err) {
		// 创建目录
		err := os.MkdirAll(downloadImageDir, 0755)
		if err != nil {
			return fmt.Errorf("创建目录失败：%s", err.Error())
		}
	}

	// 读取文件内容
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取文件失败：%s", err.Error())
	}

	// 正则表达式匹配图片地址
	re := regexp.MustCompile(`!\[.*?\]\((.+?)\)`)
	// 查找所有匹配项
	matches := re.FindAllStringSubmatch(string(content), -1)

	// 下载并替换图片地址
OuterLoop:
	for key, match := range matches {
		imageURL := match[1]

		// 判断imageURL是否为网络地址
		if !strings.HasPrefix(imageURL, "http") {
			continue
		}

		// 如果图片地址已经处理过，则跳过
		for i := key - 1; i >= 0; i-- {
			if matches[i][1] == imageURL {
				continue OuterLoop
			}
		}

		// 获取文件名，用于保存到本地目录
		imageName := fmt.Sprintf("%s_%d", strings.TrimSuffix(path.Base(filePath), filepath.Ext(filePath)), key)

		// 判断常见图片格式，如果不是常见图片格式则默认为jpeg格式
		if strings.HasSuffix(imageURL, ".png") {
			imageName += ".png"
		} else if strings.HasSuffix(imageURL, ".jpg") {
			imageName += ".jpg"
		} else if strings.HasSuffix(imageURL, ".jpeg") {
			imageName += ".jpeg"
		} else {
			imageName += ".jpeg"
		}
		imageFilePath := filepath.Join(downloadImageDir, imageName)

		// 下载图片到本地
		err := downloadImage(imageURL, imageFilePath)
		if err != nil {
			fmt.Printf("下载图片失败：%s\n", err.Error())
			continue
		}

		// 替换文档中的图片地址
		newImageURL := filepath.Join(*imagePath, imageName)
		// 替换文档内容中的图片地址
		content = []byte(strings.ReplaceAll(string(content), imageURL, newImageURL))
	}

	// 将修改后的内容写回文件
	if *rewrite == "y" {
		err = ioutil.WriteFile(filePath, content, 0755)
		if err != nil {
			return fmt.Errorf("写回文件失败：%s", err.Error())
		}
	}

	return nil
}

// 下载图片
func downloadImage(url string, filePath string) error {
	// 创建一个自定义的 http.Client
	client := &http.Client{}

	// 构建一个 GET 请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// 设置请求头参数，模拟 Chrome 浏览器的 User-Agent 和 Referer
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载图片失败，状态码：%d", resp.StatusCode)
	}

	// 创建文件
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 将响应体的内容写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("图片已下载到：%s\n", filePath)
	return nil
}
