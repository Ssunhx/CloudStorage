package oss

import (
	"cloudstorage/config"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var ossCli *oss.Client

func Client() *oss.Client {
	if ossCli != nil {
		return ossCli
	}
	ossCli, err := oss.New(config.OSSEndpoint, config.OSSAccesskeyID, config.OSSAccessKeySecret)
	if err != nil {
		return nil
	}
	return ossCli
}

// 获取 Bucket 存储空间
func Bucket() *oss.Bucket {
	cli := Client()
	if cli != nil {
		bucket, err := cli.Bucket(config.OSSBucket)
		if err != nil {
			return nil
		}
		return bucket
	}
	return nil
}

// oss 临时授权下载 URL
func DownloadURL(objName string) string {
	signedURL, err := Bucket().SignURL(objName, oss.HTTPGet, 3600)
	if err != nil {
		return ""
	}
	return signedURL
}

// oss 设定bucket生命周期
func BuildLifecycleRule(bucketName string) {
	reluTest1 := oss.BuildLifecycleRuleByDays("rule1", "test/", true, 30)
	rules := []oss.LifecycleRule{reluTest1}

	Client().SetBucketLifecycle(bucketName, rules)
}
