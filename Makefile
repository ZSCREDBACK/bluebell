# Makefile 适用于Mac OS 和 Linux

# 用来定义伪目标;不创建目标文件,而是去执行这个目标下面的命令
.PHONY: all build run gotool clean help

# 定义变量,这里是用作文件名
BINARY="bluebell"

# 执行all(make)等于执行make gotool + make build
all: gotool build

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARY}

run:
	# @表示不把实际的命令输出到终端
	@go run ./

gotool:
	# 格式换代码
	go fmt ./
	# 检查代码
	go vet ./

clean:
	# 支持使用shell命令
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

help:
	# 打印一些提示信息
	@echo "make - 格式化 Go 代码, 并编译生成二进制文件"
	@echo "make build - 编译 Go 代码, 生成二进制文件"
	@echo "make run - 直接运行 Go 代码"
	@echo "make clean - 移除二进制文件和 vim swap files"
	@echo "make gotool - 运行 Go 工具 'fmt' and 'vet'"