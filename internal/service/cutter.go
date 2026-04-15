package service

import (
	"github.com/samantonio28/meowcut/pkg/utils"
)

type LinkMeowCutter struct{}

func NewLinkMeowCutter() *LinkMeowCutter {
	return &LinkMeowCutter{}
}

// Cut реализует domain.LinkCutter.
func (c *LinkMeowCutter) Cut(originalURL string) (string, error) {
	return utils.GenerateShortID()
}