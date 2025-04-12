package dicontainer

import (
	"vnc-summarizer/adapters/databases/s3"
	interfaces "vnc-summarizer/core/interfaces/s3"
)

func GetAwsS3() interfaces.AwsS3 {
	return s3.NewAwsS3()
}
