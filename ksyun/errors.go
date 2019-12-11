package ksyun

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"strings"
	"time"
)

const (
	NotFound = "Notfound"
)

type ProviderError struct {
	errorCode string
	message   string
}

func (e *ProviderError) Error() string {
	return fmt.Sprintf("[ERROR] Terraform Ksyun Provider Error: Code: %s Message: %s", e.errorCode, e.message)
}

func (err *ProviderError) ErrorCode() string {
	return err.errorCode
}

func (err *ProviderError) Message() string {
	return err.message
}

func newNotFoundError(str string) error {
	return &ProviderError{
		errorCode: NotFound,
		message:   str,
	}
}

func getNotFoundMessage(product, id string) string {
	return fmt.Sprintf("the specified %s %s is not found", product, id)
}

func isNotFoundError(err error) bool {
	if e, ok := err.(*ProviderError); ok &&
		(e.ErrorCode() == NotFound || strings.Contains(strings.ToLower(e.Message()), NotFound)) {
		return true
	}

	return false
}
func notFoundError(err error) bool {
	errMessage := strings.ToLower(err.Error())
	if strings.Contains(errMessage, "notfound") ||
		strings.Contains(errMessage, "invalid") ||
		strings.Contains(errMessage, "not_found") {
		//strings.Contains(errMessage,"notfound"){
		return true
	}
	return false
}
func inUseError(err error) bool {
	errMessage := strings.ToLower(err.Error())
	if strings.Contains(errMessage, "inuse") ||
		strings.Contains(errMessage, "in use") ||
		strings.Contains(errMessage, "used") {
		return true
	}
	return false
}

func retryOnAwsCode(code string, f func() (interface{}, error)) (interface{}, error) {
	var resp interface{}
	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		var err error
		resp, err = f()
		if err != nil {
			awsErr, ok := err.(awserr.Error)
			if ok && awsErr.Code() == code {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	return resp, err
}

func isAWSErr(err error, code string, message string) bool {
	var awsErr awserr.Error
	if errors.As(err, &awsErr) {
		return awsErr.Code() == code && strings.Contains(awsErr.Message(), message)
	}
	return false
}
