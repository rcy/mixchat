package rndcolor

import (
	"fmt"
	"hash/crc32"
)

// Return a random tailwind color from a checksum of string
func FromString0(item string) string {
	return "text-white"
}
func FromString(item string) string {
	// TODO: access full set of tailwind colors
	colors := []string{
		"red",
		"orange",
		"yellow",
		"green",
		"teal",
		"blue",
		"indigo",
		"purple",
		"pink",
	}
	shades := []string{
		"400",
		"500",
		"600",
	}

	sum := int(crc32.ChecksumIEEE([]byte(item)))

	return fmt.Sprintf("text-%s-%s", colors[sum%len(colors)], shades[sum%len(shades)])
}
