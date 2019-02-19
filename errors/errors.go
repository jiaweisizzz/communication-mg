package except

import "errors"

var (
	//ChannelNotFound 通道不存在
	ChannelNotFound = errors.New("通道不存在")
	//ReadDirError 读取目录失败
	ReadDirError = errors.New("读取目录失败")
)
