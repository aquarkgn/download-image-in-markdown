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
	markdownPath = flag.String("markdownPath", "", "包含markdown文件的目录")
	targetPath   = "source/image" // 图片保存的目录
)

func main() {
	log.SetFlags(0)
	defer log.Println("INFO: 下载完成.")

	// 检索CLI标志和执行环境
	flag.Parse()

	if *markdownPath == "" {
		*markdownPath, _ = os.Getwd()
	}

	targetPath = path.Join(*markdownPath, targetPath)

	// 循环处理目录下的所有markdown文件
	err := filepath.Walk(*markdownPath, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".md" {
			log.Printf("INFO: Processing %s\n", path)
			err := processFile(path, targetPath)
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
func processFile(filePath string, destDir string) error {
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		// 创建目录
		err := os.MkdirAll(destDir, 0755)
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
	for _, match := range matches {
		imageURL := match[1]

		// 获取文件名，用于保存到本地目录
		imageName := filepath.Base(imageURL)
		destFilePath := filepath.Join(destDir, imageName)

		// 下载图片到本地
		err := downloadImage(imageURL, destFilePath)
		if err != nil {
			fmt.Printf("下载图片失败：%s\n", err.Error())
			continue
		}

		// 替换文档中的图片地址
		newImageURL := filepath.Join("/source/image", imageName)
		//// 将原始的图片URL进行转义，用于正则表达式替换
		escapedImageURL := regexp.QuoteMeta(imageURL)
		// 替换文档的内容
		content = []byte(strings.ReplaceAll(string(content), escapedImageURL, newImageURL))
	}

	// 将修改后的内容写回文件
	err = ioutil.WriteFile(filePath, content, 0755)
	if err != nil {
		return fmt.Errorf("写回文件失败：%s", err.Error())
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
