# rebed
[![codecov](https://codecov.io/gh/soypat/rebed/branch/main/graph/badge.svg)](https://codecov.io/gh/soypat/rebed)
[![Build Status](https://travis-ci.org/soypat/rebed.svg?branch=main)](https://travis-ci.org/soypat/rebed)
[![Go Report Card](https://goreportcard.com/badge/github.com/soypat/rebed)](https://goreportcard.com/report/github.com/soypat/rebed)
[![go.dev reference](https://pkg.go.dev/badge/github.com/soypat/rebed)](https://pkg.go.dev/github.com/soypat/rebed)

Recreate embedded filesystems from embed.FS type in current working directory. 

Expose the files you've embedded in your binary so users can see and/or tinker with them. See [where is this useful](#where-is-this-useful) for an application example.

## Four actions available:

```go

//go:embed someFS/*
var bdFS embed.FS

// Just replicate folder Structure
rebed.Tree(bdFS)

// Make empty files
rebed.Touch(bdFS)

// Recreate entire FS
rebed.Create(bdFS)

// Recreate FS without modifying existing files
rebed.Patch(bdFS)
```

### Where is this useful?
You could theoretically embed your web development assets folder and deploy it. The binary could run from the working directory folder generated by `rebed` and users could modify the assets and change the website's look (to do this run `rebed.Patch` to not overwrite modified files). If asset files are lost or the site breaks, running `rebed.Create` replaces all files with original ones.
