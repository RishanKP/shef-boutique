package aws

import (
 "github.com/aws/aws-sdk-go/aws"
 "github.com/aws/aws-sdk-go/aws/credentials"
 "github.com/aws/aws-sdk-go/aws/session"
 "os"
)

func ConnectAws() *session.Session {
 AccessKeyID := os.Getenv("AWS_ACCESS_KEY")
 SecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
 MyRegion := os.Getenv("AWS_REGION")


 sess, err := session.NewSession(
  &aws.Config{
   Region: aws.String(MyRegion),
   Credentials: credentials.NewStaticCredentials(
    AccessKeyID,
    SecretAccessKey,
    "", // a token will be created when the session it's used.
   ),
  })

  if err != nil {
   panic(err)
  }

  return sess
}
