package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

func main() {
	var (
		keyID     = flag.String("keyId", "arn:aws:kms:us-east-1:681289614897:key/9f1814c0-2781-4e65-b3f9-5ed8b749dda1", "A unique identifier for the customer master key")
		plaintext = flag.String("plaintext", "", "Data to be encrypted")
	)
	flag.Parse()

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("failed to load config, " + err.Error())
	}

	svc := kms.New(cfg)
	input := &kms.EncryptInput{
		KeyId:     aws.String(*keyID),
		Plaintext: []byte(*plaintext),
	}

	req := svc.EncryptRequest(input)
	result, err := req.Send()
	if err != nil {
		log.Fatal(err)
	}
	ciphertext := base64.StdEncoding.EncodeToString(result.CiphertextBlob)
	fmt.Println(ciphertext)
}
