package s3access

import (
	"bytes"
	"context"
	"errors"
	common_conf "github.com/kmahyyg/my-gin-components/common-conf"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"net/url"
	"path"
	"time"
)

const (
	TEST_CREATE_FILENAME = "testp/B8C63255-60EE-412A-BC39-C99594B56BE4.txt"
	TEST_CREATE_FILEDATA = "abcdefghijklmnopqrstuvwxyz"
)

var (
	ErrClientNotBuilt = errors.New("s3 client is not built or no config provided")
)

type S3ClientFactory struct {
	isBuilt      bool
	s3Client     *minio.Client
	s3Config     *common_conf.S3Config
	uriSchema    string
	uriPrefixStr string
	uriPrefix    *url.URL
}

func (s3cf *S3ClientFactory) BuildS3ClientFactory(s3conf *common_conf.S3Config) {
	if s3cf.isBuilt {
		return
	}
	if s3conf.UseTLS {
		s3cf.uriSchema = "https"
	} else {
		s3cf.uriSchema = "http"
	}
	s3cf.s3Config = s3conf
	var err error = nil
	s3cf.s3Client, err = minio.New(s3conf.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3conf.AKID, s3conf.AKSK, ""),
		Secure: s3conf.UseTLS,
		Region: s3conf.Region,
	})
	if err != nil {
		panic(err)
	}
	if len(s3conf.ReverseProxyEndPoint) != 0 {
		// reverse proxy enabled
		s3cf.uriPrefixStr = s3conf.ReverseProxyEndPoint
	}
	s3cf.isBuilt = true
	s3cf.uriPrefix = &url.URL{}
}

func (s3cf *S3ClientFactory) GetS3ClientInstance() (*minio.Client, error) {
	if s3cf.isBuilt && s3cf.s3Client != nil {
		return s3cf.s3Client, nil
	}
	if !s3cf.isBuilt {
		return nil, ErrClientNotBuilt
	}
	return nil, nil
}

func (s3cf *S3ClientFactory) ResetS3ClientConfig() {
	s3cf.isBuilt = false
	s3cf.s3Client = nil
	s3cf.uriSchema = "http"
	s3cf.uriPrefixStr = ""
	s3cf.uriPrefix = nil
}

func (s3cf *S3ClientFactory) UploadFileWithPath(fdpath string, uploadPath string) (string, error) {
	if !s3cf.isBuilt {
		return "", ErrClientNotBuilt
	}
	ctx := context.Background()
	upinfo, err := s3cf.s3Client.FPutObject(ctx, s3cf.s3Config.Bucket, uploadPath, fdpath,
		minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return "", err
	}
	return upinfo.Key, nil
}

func (s3cf *S3ClientFactory) UploadFileWithBinary(filename string, filedata []byte) (string, error) {
	if !s3cf.isBuilt {
		return "", ErrClientNotBuilt
	}
	ctx := context.Background()
	bf := bytes.NewReader(filedata)
	// upload test file
	upinfo, err := s3cf.s3Client.PutObject(ctx, s3cf.s3Config.Bucket, filename, bf, -1, minio.PutObjectOptions{
		ContentType: "application/octet-stream", // use application/octet-stream for binary
	})
	if err != nil {
		// UploadInfo.Key will include subpath of file, due to K-V based OSS implementation.
		return upinfo.Key, err
	}
	return "", nil
}

func (s3cf *S3ClientFactory) DownloadFile(fdpath string) ([]byte, error) {
	//TODO
	// call: downloadToTemp, then rename to UUIDv4-based, move to specific restricted folder
	return nil, nil
}

func (s3cf *S3ClientFactory) downloadToTemp(fdpath string) (string, error) {
	//TODO
	return "", nil
}

// Download2PresignedURL returns presigned url for file download using minio client
func (s3cf *S3ClientFactory) Download2PresignedURL(fdpath string) (string, error) {
	// check if inited
	if !s3cf.isBuilt {
		return "", ErrClientNotBuilt
	}
	// get file download url
	ctx := context.Background()
	reqParam := make(url.Values)
	reqParam.Set("utm_source", "mgc")
	presignURL, err := s3cf.s3Client.PresignedGetObject(ctx, s3cf.s3Config.Bucket, fdpath,
		time.Hour*1, reqParam)
	if err != nil {
		return "", err
	}
	// rewrite to reverse proxy url
	return s3cf.replaceUri2ReverseProxy(presignURL), nil
}

// replaceUri2ReverseProxy: replace uri to reverse proxy uri
func (s3cf *S3ClientFactory) replaceUri2ReverseProxy(oripath *url.URL) string {
	if s3cf.uriPrefixStr == "" {
		return oripath.String()
	}
	// uriPrefixStr parse from string to *url.Values
	var err error
	s3cf.uriPrefix, err = url.Parse(s3cf.uriPrefixStr)
	if err != nil {
		// prefix uri not illegal
		panic(err)
	}
	// replace schema
	oripath.Scheme = s3cf.uriSchema
	// replace host
	oripath.Host = s3cf.uriPrefix.Host
	// concat path
	oripath.Path = path.Join(s3cf.uriPrefix.Path, oripath.Path)
	// convert to string
	return oripath.String()
}

// initBucket will create bucket if not exist and check writable
func (s3cf *S3ClientFactory) initBucket() error {
	ctx := context.Background()
	// check if client working
	if !s3cf.isBuilt {
		return ErrClientNotBuilt
	}
	locationStr := "us-east-1"
	// check if bucket exists
	exists, errExists := s3cf.s3Client.BucketExists(ctx, s3cf.s3Config.Bucket)
	if errExists != nil {
		return errExists
	}
	if !exists {
		// create bucket
		err := s3cf.s3Client.MakeBucket(ctx, s3cf.s3Config.Bucket,
			minio.MakeBucketOptions{Region: locationStr})
		if err != nil {
			return err
		}
	}
	// upload
	_, err := s3cf.UploadFileWithBinary(TEST_CREATE_FILENAME, []byte(TEST_CREATE_FILEDATA))
	if err != nil {
		return err
	}
	return nil
}
