package main

import (
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGolangTestTask(t *testing.T) {
	gin.SetMode("test")
	RegisterFailHandler(Fail)
	RunSpecs(t, "GolangTestTask Suite")
}
