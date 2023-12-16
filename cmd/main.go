package main

import (
	"github.com/yash492/statusy/pkg/config"
	"github.com/yash492/statusy/pkg/domain"
	"github.com/yash492/statusy/pkg/store"
)

func main() {
	config.New()
	store.New()
	domain.New()
}
