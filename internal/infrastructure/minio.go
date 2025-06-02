package infrastructure

// import (
// 	"bytes"
// 	"context"
// 	"io"
// 	"log"
// 	"net/http"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	"github.com/minio/minio-go/v7"
// 	"github.com/minio/minio-go/v7/pkg/credentials"
// 	"github.com/spf13/viper"
// 	"gitlab.ai-vlab.com/cygate/common/pkg/config/common"
// 	cipher "gitlab.ai-vlab.com/cygate/common/pkg/config/encrypt"
// )

// type Upload struct {
// 	Url      string `json:"url"`
// 	FileName string `json:"fileName"`
// }

// var client *minio.Client
// var returnUrl string

// func InitMinioClient() {
// 	var err error
// 	client, err = minio.New(viper.GetString("minio.url"), &minio.Options{
// 		Creds:  credentials.NewStaticV4(cipher.DecryptEnv("minio.access-key"), cipher.DecryptEnv("minio.secret-key"), ""),
// 		Secure: viper.GetBool("minio.secure"),
// 	})
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	returnUrl = viper.GetString("minio.return-url")
// }

// func UploadFile(c *gin.Context) (*common.Upload, *common.Error) {
// 	file, fileHeader, err := c.Request.FormFile("file")
// 	if err != nil {
// 		return nil, common.FileEmpty
// 	}
// 	bucket := c.Request.FormValue("bucket")
// 	rootPath := c.Request.FormValue("rootPath")
// 	defer file.Close()
// 	fileBytes, err := io.ReadAll(file)
// 	ctx := context.Background()
// 	if len(bucket) == 0 {
// 		bucket = viper.GetString("minio.default-bucket")
// 	}

// 	err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{Region: viper.GetString("minio.location")})
// 	if err != nil {
// 		// Check to see if we already own this bucket (which happens if you run this twice)
// 		_, errBucketExists := client.BucketExists(ctx, bucket)
// 		if errBucketExists != nil {
// 			log.Println(err)
// 			return nil, common.CreateBucketFailed
// 		}
// 	}

// 	dt := time.Now()
// 	objectName := dt.Format("20060102150405") + "_" + uuid.New().String()
// 	if len(rootPath) > 0 {
// 		objectName = rootPath + "/" + objectName
// 	}

// 	contentType := http.DetectContentType(fileBytes)

// 	log.Printf("start upload file %s", objectName)
// 	// Upload the zip file with FPutObject
// 	info, err := client.PutObject(ctx, bucket, objectName, bytes.NewReader(fileBytes), fileHeader.Size, minio.PutObjectOptions{
// 		ContentType: contentType,
// 		UserMetadata: map[string]string{
// 			"filename": fileHeader.Filename,
// 		}})
// 	if err != nil {
// 		log.Println(err)
// 		return nil, common.UploadFileFailed
// 	}
// 	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
// 	return &common.Upload{Url: returnUrl + "/" + bucket + "/" + objectName, FileName: fileHeader.Filename}, nil
// }
