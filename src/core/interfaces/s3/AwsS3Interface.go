package s3

type AwsS3 interface {
	SavePropositionImage(propositionCode int, image []byte) (string, error)
}
