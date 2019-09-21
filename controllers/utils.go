package controllers

import (
	"bytes"
	"context"
	"io/ioutil"
	"os/exec"
	"strings"
)

const finalizerName = "finalizer.teleport-plus.gravitational.com"

const (
	ConditionRegistered = "Registered"
	ConditionFailed     = "Failed"
)

func execTctl(ctx context.Context, args ...string) ([]byte, []byte, error) {
	var stdout, stderr bytes.Buffer
	cmdArgs := []string{"-c", "/etc/teleport/teleport.yaml"}
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.CommandContext(ctx, "/tctl", cmdArgs...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.Bytes(), stderr.Bytes(), err
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}

func ownNamespace() (string, error) {
	data, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		return "", err
	}
	ns := strings.TrimSpace(string(data))
	if len(ns) == 0 {
		return "", err
	}
	return ns, nil
}
