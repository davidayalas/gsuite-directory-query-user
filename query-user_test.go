package main

import (
        "testing"
        "os"
        "fmt"
    )

func TestGetUser(t *testing.T){

    if GoogleClientErr != nil{
    	t.Errorf("%s", GoogleClientErr)
    }

    _, err := GetUser(os.Getenv("impersonatedUser"))
    if err != nil{
        t.Errorf("%s", err)
    }
}

func TestGetUserKO(t *testing.T){

    if GoogleClientErr != nil{
        t.Errorf("%s", GoogleClientErr)
    }

    _, err := GetUser(fmt.Sprintf("xyzsasdadasd%s", os.Getenv("impersonatedUser")))

    if err == nil{
        t.Errorf("%s", err)
    }
}

func init(){
    GC, GCE := CreateClient()
    GoogleClient = GC
    GoogleClientErr = GCE	
}