package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

func main() {
	var (
		ciphertext = flag.String("ciphertext", "", "Ciphertext to be decrypted")
	)
	flag.Parse()

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("failed to load config, " + err.Error())
	}

	ciphertextBlob, err := base64.StdEncoding.DecodeString(*ciphertext)
	if err != nil {
		log.Fatal(err)
	}

	svc := kms.New(cfg)
	input := &kms.DecryptInput{
		CiphertextBlob: []byte(ciphertextBlob),
	}

	req := svc.DecryptRequest(input)
	result, err := req.Send()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(result.Plaintext))
}
