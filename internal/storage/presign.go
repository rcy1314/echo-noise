package storage

import (
    "context"
    "time"

    "github.com/aws/aws-sdk-go-v2/aws"
    awscfg "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
    "github.com/lin-snow/ech0/internal/models"
)

func getAWSConfigFromSite(cfg models.SiteConfig) (aws.Config, error) {
    cr := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(cfg.StorageAccessKey, cfg.StorageSecretKey, ""))
    region := cfg.StorageRegion
    if region == "" {
        region = "auto"
    }
    base := awscfg.WithRegion(region)
    base2 := awscfg.WithCredentialsProvider(cr)
    return awscfg.LoadDefaultConfig(context.Background(), base, base2, awscfg.WithEndpointResolverWithOptions(
        aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
            if cfg.StorageEndpoint != "" {
                return aws.Endpoint{
                    URL:               cfg.StorageEndpoint,
                    SigningRegion:     region,
                    HostnameImmutable: true,
                }, nil
            }
            return aws.Endpoint{}, &aws.EndpointNotFoundError{}
        }),
    ))
}

func getS3Client(site models.SiteConfig) (*s3.Client, error) {
    awsCfg, err := getAWSConfigFromSite(site)
    if err != nil {
        return nil, err
    }
    cli := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
        if site.StorageUsePathStyle {
            o.UsePathStyle = true
        }
    })
    return cli, nil
}

func PresignUpload(site models.SiteConfig, bucket, key string, expires time.Duration, contentType string) (string, error) {
    cli, err := getS3Client(site)
    if err != nil {
        return "", err
    }
    presigner := s3.NewPresignClient(cli, func(po *s3.PresignOptions) {
        po.Expires = expires
    })
    input := &s3.PutObjectInput{
        Bucket:      &bucket,
        Key:         &key,
        ContentType: &contentType,
        ACL:         s3types.ObjectCannedACLPrivate,
    }
    req, err := presigner.PresignPutObject(context.Background(), input)
    if err != nil {
        return "", err
    }
    return req.URL, nil
}

func PresignDownload(site models.SiteConfig, bucket, key string, expires time.Duration) (string, error) {
    cli, err := getS3Client(site)
    if err != nil {
        return "", err
    }
    presigner := s3.NewPresignClient(cli, func(po *s3.PresignOptions) {
        po.Expires = expires
    })
    input := &s3.GetObjectInput{
        Bucket: &bucket,
        Key:    &key,
    }
    req, err := presigner.PresignGetObject(context.Background(), input)
    if err != nil {
        return "", err
    }
    return req.URL, nil
}
