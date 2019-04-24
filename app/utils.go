// Copyright © 2019 Timo Runge <me@timorunge.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
//    this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors
//    may be used to endorse or promote products derived from this software
//    without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/viper"
)

const (
	httpSchema  = "http"
	httpsSchema = "https"
	ossEndpoint = "%s://oss-%s.aliyuncs.com"
)

// OssEndpoint returns the endpoint which will be used for the client
// connection.
func OssEndpoint(region string, disableTLS bool) string {
	scheme := httpsSchema
	if disableTLS {
		scheme = httpSchema
	}
	return fmt.Sprintf(ossEndpoint, scheme, region)
}

// LocalFileInfo denotes required information for the local files.
type LocalFileInfo struct {
	AbsDir  string
	AbsPath string
	Exists  bool
	IsDir   bool
	Name    string
	Size    int64
}

// LocalFile returns LocalFileInfo for on OSS object.
func LocalFile(ossObject oss.ObjectProperties) (*LocalFileInfo, error) {
	absPath := filepath.Clean(fmt.Sprintf("%s/%s", viper.GetString("destinationDir"), ossObject.Key))
	fileExists := false
	fileSize := int64(0)
	isDir := false

	file, err := os.Stat(absPath)
	if err == nil {
		fileExists = true
		fileSize = file.Size()
		isDir = file.IsDir()
	}

	return &LocalFileInfo{
		AbsDir:  filepath.Dir(absPath),
		AbsPath: absPath,
		Exists:  fileExists,
		IsDir:   isDir,
		Name:    ossObject.Key,
		Size:    fileSize,
	}, err
}
