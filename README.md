# AWS Lambda GSuite Directory User Query

This little Golang microservice queries GSuite Directory with Domain-Wide Delegation of Authority: https://developers.google.com/admin-sdk/directory/v1/guides/delegation

## Steps to perform in Google 

1. Create an app in http://console.cloud.google.com to attach a service account
1. Enable Admin SDK api in your new application (https://console.cloud.google.com/apis/dashboard)
1. Create a service account https://console.developers.google.com/iam-admin/serviceaccounts in your newly created application
1. Delegate domain-wide authority to your service account with scope "https://www.googleapis.com/auth/admin.directory.user.readonly" at https://admin.google.com/ManageOauthClients

## Steps to perform in AWS

1. https://docs.aws.amazon.com/lambda/latest/dg/lambda-go-how-to-create-deployment-package.html
1. Build the query-user.go executable
1. Zip it
1. Create your lambda function 
1. Setup environment variables:
	- **gsuite_credentials**: the value is the JSON downloaded in the step 2 from Google console (service account)
	- **impersonatedUser**: user to impersonate

