package main

import (
	"os"
	"os/exec"
	"testing"
)

const apigeecli = "./apigeecli"

var token = os.Getenv("APIGEE_TOKEN")
var env = os.Getenv("APIGEE_ENV")
var org = os.Getenv("APIGEE_ORG")

func TestMain(t *testing.T) {
	cmd := exec.Command(apigeecli)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

// orgs test
func TestListOrgs(t *testing.T) {
	cmd := exec.Command(apigeecli, "orgs", "list", "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetOrg(t *testing.T) {
	cmd := exec.Command(apigeecli, "orgs", "get", "-o", org, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetMart(t *testing.T) {
	mart := os.Getenv("MART")
	cmd := exec.Command(apigeecli, "orgs", "setmart", "-o", org, "-m", mart, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

// env tests
func TestListEnvs(t *testing.T) {
	cmd := exec.Command(apigeecli, "envs", "list", "-o", org, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetEnv(t *testing.T) {
	cmd := exec.Command(apigeecli, "envs", "get", "-o", org, "-e", env, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

// developers test
func TestCreateDeveloper(t *testing.T) {
    email := "test@example.com"
    first := "frstname"
    last := "lastname"
    user := "username"
    
    cmd := exec.Command(apigeecli, "developers", "create", "-o", org, "-n", email, "-f", first, "-s", last, "-u", user, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetDeveloper(t *testing.T) {
    email := "test@example.com"
    
    cmd := exec.Command(apigeecli, "developers", "get", "-o", org, "-n", email, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestListDeveloper(t *testing.T) {
    
    cmd := exec.Command(apigeecli, "developers", "list", "-o", org, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestListExpandDeveloper(t *testing.T) {
    expand := "true"
    
    cmd := exec.Command(apigeecli, "developers", "list", "-o", org, "-x", expand, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteDeveloper(t *testing.T) {
    email := "test@example.com"
    
    cmd := exec.Command(apigeecli, "developers", "delete", "-o", org, "-n", email, "-t", token)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

// kvm test

