package ceph

import (
	"gopkg.in/amz.v1/aws"
	"gopkg.in/amz.v1/s3"
)

var cephConn *s3.S3

func GetcephConnection() *s3.S3 {
	if cephConn != nil {
		return cephConn
	}
	// 1、初始化 ceph 的一些信息
	auth := aws.Auth{
		AccessKey: "your key",
		SecretKey: "your key",
	}
	curRegion := aws.Region{
		Name:                 "default",
		EC2Endpoint:          "http://127.0.0.1:9080",
		S3Endpoint:           "http://127.0.0.1:9080",
		S3BucketEndpoint:     "",
		S3LocationConstraint: false,
		S3LowercaseBucket:    false,
		Sign:                 aws.SignV2,
	}

	// 2、创建 S3 类型的连接
	cephConn = s3.New(auth, curRegion)
	return cephConn
}

// 获取指定的 bucket 对象
func GetCephBucket(bucket string) *s3.Bucket {
	conn := GetcephConnection()
	return conn.Bucket(bucket)
}

// 上传文件到 ceph
func Putobject(bucket string, path string, data []byte) error {
	return GetCephBucket(bucket).Put(path, data, "object-stream", s3.PublicRead)
}
