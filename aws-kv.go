package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
)

const (
	colorGreen = "\033[32m"
	colorRed   = "\033[31m"
	colorReset = "\033[0m"
)

type CredentialInfo struct {
	AccessKey string
	Valid     bool
	Details   string
}

func validateAndFetchDetails(accessKey, secretKey string) *CredentialInfo {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	if err != nil {
		return &CredentialInfo{
			AccessKey: accessKey,
			Valid:     false,
		}
	}

	stsSvc := sts.New(sess)
	callerIdentity, err := stsSvc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return &CredentialInfo{
			AccessKey: accessKey,
			Valid:     false,
		}
	}

	iamSvc := iam.New(sess)
	authDetails, err := iamSvc.GetAccountAuthorizationDetails(&iam.GetAccountAuthorizationDetailsInput{})
	if err != nil {
		return &CredentialInfo{
			AccessKey: accessKey,
			Valid:     true,
			Details: fmt.Sprintf(
				"Caller Identity: %s, Account: %s, ARN: %s",
				*callerIdentity.UserId,
				*callerIdentity.Account,
				*callerIdentity.Arn,
			),
		}
	}

	return &CredentialInfo{
		AccessKey: accessKey,
		Valid:     true,
		Details: fmt.Sprintf(
			"Caller Identity: %s, Account: %s, ARN: %s\nIAM Details: %v",
			*callerIdentity.UserId,
			*callerIdentity.Account,
			*callerIdentity.Arn,
			authDetails,
		),
	}
}

func processCredentials(creds []string, wg *sync.WaitGroup, results chan<- *CredentialInfo) {
	defer wg.Done()

	for _, cred := range creds {
		parts := strings.Split(cred, ":")
		if len(parts) != 2 {
			results <- &CredentialInfo{
				AccessKey: cred,
				Valid:     false,
			}
			continue
		}

		accessKey, secretKey := parts[0], parts[1]
		results <- validateAndFetchDetails(accessKey, secretKey)
	}
}

func main() {
	keyFlag := flag.String("k", "", "AWS_KEY:secret pair to validate")
	fileFlag := flag.String("f", "", "File containing AWS_KEY:secret pairs (one per line)")
	threadsFlag := flag.Int("t", 1, "Number of goroutines to use")
	debugFlag := flag.Bool("d", false, "Enable debug mode to show invalid credentials")
	flag.Parse()

	var credentials []string

	// Collect credentials from the `-k` flag
	if *keyFlag != "" {
		credentials = append(credentials, *keyFlag)
	}

	// Collect credentials from the `-f` flag
	if *fileFlag != "" {
		file, err := os.Open(*fileFlag)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			credentials = append(credentials, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}
	}

	// Split credentials into chunks for goroutines
	chunkSize := (len(credentials) + *threadsFlag - 1) / *threadsFlag
	var wg sync.WaitGroup
	results := make(chan *CredentialInfo, len(credentials))

	for i := 0; i < len(credentials); i += chunkSize {
		end := i + chunkSize
		if end > len(credentials) {
			end = len(credentials)
		}

		wg.Add(1)
		go processCredentials(credentials[i:end], &wg, results)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(results)

	// Counters for valid and invalid credentials
	validCount := 0
	invalidCount := 0

	// Print results based on the flags
	for result := range results {
		if result.Valid {
			validCount++
			fmt.Printf("%sValid:%s %s\nDetails:\n%s\n\n", colorGreen, colorReset, result.AccessKey, result.Details)
		} else {
			invalidCount++
			if *debugFlag {
				fmt.Printf("%sInvalid:%s %s\n", colorRed, colorReset, result.AccessKey)
			}
		}
	}

	// Print summary
	fmt.Printf("%sNumber of valid credentials:%s %d\n", colorGreen, colorReset, validCount)
	if *debugFlag {
		fmt.Printf("%sNumber of invalid credentials:%s %d\n", colorRed, colorReset, invalidCount)
	}
}
