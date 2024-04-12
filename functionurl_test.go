package lambroll_test

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/fujiwara/lambroll"
	"github.com/go-test/deep"
)

var permissionsTestCases = []struct {
	subject                string
	statementJSON          []byte
	expectedPrincipal      *string
	expectedPrincipalOrgID *string
	expectedSourceArn      *string
}{
	{
		subject: "AuthType NONE",
		statementJSON: []byte(`{
			"Action": "lambda:InvokeFunctionUrl",
			"Condition": {
				"StringEquals": {
					"lambda:FunctionUrlAuthType": "NONE"
				}
			},
			"Effect": "Allow",
			"Principal": "*",
			"Resource": "arn:aws:lambda:ap-northeast-1:123456789012:function:hello",
			"Sid": "lambroll-8f4ec83e623a309d9ca15db9276da30b2129be9c"
		}`),
		expectedPrincipal:      aws.String("*"),
		expectedPrincipalOrgID: nil,
	},
	{
		subject: "AuthType AWS_IAM with Principal OrgID",
		statementJSON: []byte(`{
			"Sid": "lambroll-622ed5c2bb0714ef0af1929fcea568e4ba0c4dbe",
			"Effect": "Allow",
			"Principal": "*",
			"Action": "lambda:InvokeFunctionUrl",
			"Resource": "arn:aws:lambda:ap-northeast-1:1234567890:function:hello",
			"Condition": {
				"StringEquals": {
					"lambda:FunctionUrlAuthType": "AWS_IAM",
					"aws:PrincipalOrgID": "o-xxxxxxxxxx"
				}
			}
		}`),
		expectedPrincipal:      aws.String("*"),
		expectedPrincipalOrgID: aws.String("o-xxxxxxxxxx"),
	},
	{
		subject: "AuthType AWS_IAM with Principal",
		statementJSON: []byte(`{
			"Action": "lambda:InvokeFunctionUrl",
			"Condition": {
				"StringEquals": {
					"lambda:FunctionUrlAuthType": "AWS_IAM"
				}
			},
			"Effect": "Allow",
			"Principal": {
				"AWS": "arn:aws:iam::123456789012:root"
			},
			"Resource": "arn:aws:lambda:ap-northeast-1:123456789012:function:hello",
			"Sid": "lambroll-3b135eca4b14335775cda9f947966093a57d270f"
		}`),
		expectedPrincipal:      aws.String("123456789012"),
		expectedPrincipalOrgID: nil,
	},
	{
		subject: "AuthType AWS_IAM with Principal CF OAC",
		statementJSON: []byte(`{
			"Action": "lambda:InvokeFunctionUrl",
			"Condition": {
				"ArnLike": {
					"aws:SourceArn": "arn:aws:cloudfront::123456789012:distribution/ABCDEFG12345678"
				}
			},
			"Effect": "Allow",
			"Principal": {
				"Service": "cloudfront.amazonaws.com"
			},
			"Resource": "arn:aws:lambda:ap-northeast-1:123456789012:function:hello",
			"Sid": "lambroll-3b135eca4b14335775cda9f947966093a57d270f"
		}`),
		expectedPrincipal:      aws.String("cloudfront.amazonaws.com"),
		expectedPrincipalOrgID: nil,
		expectedSourceArn:      aws.String("arn:aws:cloudfront::123456789012:distribution/ABCDEFG12345678"),
	},
}

func TestParseStatement(t *testing.T) {
	for _, c := range permissionsTestCases {
		st := &lambroll.PolicyStatement{}
		if err := json.Unmarshal(c.statementJSON, st); err != nil {
			t.Errorf("%s failed to unmarshal json: %s", c.subject, err)
		}
		if diff := deep.Equal(c.expectedPrincipal, st.PrincipalString()); diff != nil {
			t.Errorf("%s PrincipalString diff %s", c.subject, diff)
		}
		if diff := deep.Equal(c.expectedPrincipalOrgID, st.PrincipalOrgID()); diff != nil {
			t.Errorf("%s PrincipalOrgID diff %s", c.subject, diff)
		}
		if diff := deep.Equal(c.expectedSourceArn, st.SourceArn()); diff != nil {
			t.Errorf("%s SourceArn diff %s", c.subject, diff)
		}
	}
}
