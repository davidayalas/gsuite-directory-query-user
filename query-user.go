package main

import (
    "fmt"
    "context"
    "encoding/json"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-lambda-go/events"
    "log"
    "os"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/admin/directory/v1"
)

type errorString struct {
    s string
}

func (e *errorString) Error() string {
  return e.s
}

type Usuari struct {
    Email    string `json:"email"`
    FullName string `json:"fullname"`
}

type Event struct {
        Email string `json:"email"`
}

var GoogleClient *admin.Service
var GoogleClientErr error 

func CreateClient() (*admin.Service, error){
    credentials := os.Getenv("gsuite_credentials")
    impersonatedUser := os.Getenv("impersonatedUser")

    if credentials == "" {
        return nil, &errorString{"{\"590\":\"No credentials provided\"}"}
    }

    if impersonatedUser == "" {
        return nil, &errorString{"{\"591\":\"No admin impersonated user provided\"}"}
    }

    config, err := google.JWTConfigFromJSON([]byte(credentials), admin.AdminDirectoryUserReadonlyScope)

    if err != nil {
        log.Fatalf("Could not create config object => {%s}", err)
    }

    config.Subject = impersonatedUser
    client, err := admin.New(config.Client(context.Background()))

    if err != nil {
        log.Fatalf("Could not create directory service client => {%s}", err)
    }

    return client, err
}

func GetUser(usuari string) (string, error){
    if GoogleClientErr == nil{
        r, err := GoogleClient.Users.Get(usuari).Do()

        if err != nil {
            log.Printf("Unable to retrieve user in domain.", err)
            return "", &errorString{"{\"code\":\"404\",\"msg\":\"user not found\"}"}
        }

        User := new(Usuari)
        User.FullName = r.Name.FullName
        User.Email = r.PrimaryEmail
        UserData, _ := json.Marshal(User)
        return string(UserData), nil
    }

    return "", &errorString{"{\"580\":\"Unable to create GSuite Client\"}"}
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    if request.QueryStringParameters["email"] == "" {
        return events.APIGatewayProxyResponse{Body: "{\"error\" : {\"code\":\"400\",\"msg\":\"bad request, querystring param email not found\"}}", StatusCode: 400}, nil
    } 

    user, err := GetUser(request.QueryStringParameters["email"])

    if err!=nil{
        return events.APIGatewayProxyResponse{Body: fmt.Sprintf("{\"error\" : %s}", err.Error()), StatusCode: 400}, nil
    }

    return events.APIGatewayProxyResponse{Body: user, StatusCode: 200}, nil
}

func main() {
    GC, GCE := CreateClient()
    GoogleClient = GC
    GoogleClientErr = GCE
    
    lambda.Start(HandleRequest)
}