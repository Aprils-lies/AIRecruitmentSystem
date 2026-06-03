package oss

import (
	"fmt"
	"io"
	"time"

	"logic-grpc-service/config"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var client *oss.Client

func Init() error {
	cfg := config.GetOSSConfig()

	var err error
	client, err = oss.New(
		cfg.Endpoint,
		cfg.AccessKeyID,
		cfg.AccessKeySecret,
	)
	if err != nil {
		return fmt.Errorf("创建OSS客户端失败: %w", err)
	}

	return nil
}

func GetClient() *oss.Client {
	return client
}

func GetBucket() (*oss.Bucket, error) {
	cfg := config.GetOSSConfig()
	return client.Bucket(cfg.BucketName)
}

func GenerateUploadSignURL(objectKey string, expireSeconds int64, contentType string) (string, error) {
	bucket, err := GetBucket()
	if err != nil {
		return "", err
	}

	options := []oss.Option{
		oss.ContentType(contentType),
	}

	signedURL, err := bucket.SignURL(objectKey, oss.HTTPPut, expireSeconds, options...)
	if err != nil {
		return "", fmt.Errorf("生成上传签名URL失败: %w", err)
	}

	return signedURL, nil
}

func GenerateDownloadSignURL(objectKey string, expireSeconds int64) (string, error) {
	bucket, err := GetBucket()
	if err != nil {
		return "", err
	}

	signedURL, err := bucket.SignURL(objectKey, oss.HTTPGet, expireSeconds)
	if err != nil {
		return "", fmt.Errorf("生成下载签名URL失败: %w", err)
	}

	return signedURL, nil
}

func GetObjectInfo(objectKey string) (map[string][]string, error) {
	bucket, err := GetBucket()
	if err != nil {
		return nil, err
	}

	props, err := bucket.GetObjectDetailedMeta(objectKey)
	if err != nil {
		return nil, fmt.Errorf("获取对象信息失败: %w", err)
	}

	return props, nil
}

func CheckObjectExists(objectKey string) (bool, error) {
	bucket, err := GetBucket()
	if err != nil {
		return false, err
	}

	exists, err := bucket.IsObjectExist(objectKey)
	if err != nil {
		return false, fmt.Errorf("检查对象是否存在失败: %w", err)
	}

	return exists, nil
}

func GetObjectBytes(objectKey string, length int64) ([]byte, error) {
	bucket, err := GetBucket()
	if err != nil {
		return nil, err
	}

	body, err := bucket.GetObject(objectKey, oss.Range(0, length-1))
	if err != nil {
		return nil, fmt.Errorf("读取OSS文件头失败: %w", err)
	}
	defer body.Close()

	return io.ReadAll(body)
}

func GenerateObjectKey(candidateID int64, fileName string) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("resumes/%d/%d_%s", candidateID, timestamp, fileName)
}

func DeleteObject(objectKey string) error {
	bucket, err := GetBucket()
	if err != nil {
		return err
	}

	err = bucket.DeleteObject(objectKey)
	if err != nil {
		return fmt.Errorf("删除对象失败: %w", err)
	}

	return nil
}
