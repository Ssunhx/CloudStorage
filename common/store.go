package common

// 存储类型 表示文件存到哪里
type StoreType int

const (
	_ StoreType = iota
	// 节点本地
	StoreLocal
	// ceph 集群
	StoreCeph
	// 阿里 OSS
	StoreOSS
	// 混合 ceph 和 OSS
	StoreMix
	// 所有的都存储
	StoreAll
)
