# download-image-in-markdown
下载markdown中的图片到本地目录中


## 使用方法

```bash
# 下载命令包
go get -u github.com/aquarkgn/download-image-in-markdown@latest

# 安装命令包
go install github.com/aquarkgn/download-image-in-markdown@latest
```

## 执行命令
```bash
download-image-in-markdown -markdownPath=/path/to/markdown -rewrite=y
```

## 参数说明：
download-image-in-markdown
- markdownPath：markdown文件所在目录，默认为当前目录
- imagePath : 默认：source/image，参数markdownPath的相对地址
- rewrite：是否覆盖文档中图片地址，默认:n , 可选项：n/y 不覆盖/覆盖

## 示例下载图片

![image](https://th.bing.com/th/id/OIP.vVsxOjwiBfvojJ_IIqeTEAHaR7?w=144&h=349&c=7&r=0&o=5&pid=1.7) ![image](https://th.bing.com/th/id/OIP.vVsxOjwiBfvojJ_IIqeTEAHaR7?w=144&h=349&c=7&r=0&o=5&pid=1.7) ![image](https://th.bing.com/th/id/OIP.RxL0OCAKQqcmKM0u9_Y7FQHaR_?w=144&h=350&c=7&r=0&o=5&pid=1.7)

