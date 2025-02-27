package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"

	"golang.org/x/mod/semver"
)

func SemverNormalize(version string) string {
	if strings.HasPrefix(version, "v") {
		return version
	}

	return "v" + version
}

func SemverCompare(version1, version2 string) int {
	return semver.Compare(SemverNormalize(version1), SemverNormalize(version2))
}

func SemverIsStable(version string) bool {
	if strings.Contains(version, "alpha") {
		return false
	}
	if strings.Contains(version, "beta") {
		return false
	}
	if strings.Contains(version, "rc") {
		return false
	}
	if strings.Contains(version, "pre") {
		return false
	}
	if strings.Contains(version, "dev") {
		return false
	}
	if strings.Contains(version, "snapshot") {
		return false
	}
	return true
}

func Md5sum(s string) string {
	hash := md5.New()
	hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}
