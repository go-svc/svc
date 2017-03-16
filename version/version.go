// Package version 能夠協助你標註服務的版本號碼。
package version

import "fmt"

// Version 呈現了一個版本的詳細資訊。
type Version struct {
	// Major 是一個不符合相容性的 API 重大改變版本號。
	Major int
	// Minor 是符合相容性的新功能版本號。
	Minor int
	// Patch 是符合相容性的 Bug 修正版本號。
	Patch int
	// Branch 是版本的開發分歧，如 `dev` 或 `stable`。
	Branch string
}

// Define 定義了一個版本資訊。
func Define(major, minor, patch int, branch string) Version {
	return Version{
		major, minor, patch, branch,
	}
}

// String 會將版本資訊轉換成類似 `1.0.0+dev` 的字串。
func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d+%s", v.Major, v.Minor, v.Patch, v.Branch)
}
