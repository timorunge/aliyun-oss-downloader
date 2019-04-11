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
	"os"
	"path/filepath"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	golimit "github.com/lenaten/go-limit"
	"github.com/spf13/viper"
)

const (
	statusExcluded string = "EXCLUDED"
	statusFailed   string = "FAILED"
	statusGet      string = "GET"
	statusSkip     string = "SKIP"
	statusUnknown  string = "UNKNOWN"
	statusUpdate   string = "UPDATE"
)

func downloadLogger(ossBucket *oss.Bucket, ossObject oss.ObjectProperties, objectStatus string) {
	InfoLog.Printf("%s %s %s %d", ossBucket.BucketName, objectStatus, ossObject.Key, ossObject.Size)
}

func downloadOssObject(semaphore *golimit.GoLimit, ossBucket *oss.Bucket, ossObject oss.ObjectProperties, objectStatus string) {
	localFile, _ := LocalFile(ossObject)
	absDir := filepath.Dir(localFile.AbsPath)
	_, err := os.Stat(absDir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(absDir, 0755); err != nil {
			ErrorLog.Printf("Cannot create directory %s: %v", absDir, err)
		}
	}

	if err = ossBucket.GetObjectToFile(ossObject.Key, localFile.AbsPath); err != nil {
		objectStatus = statusFailed
		ErrorLog.Printf("Cannot download %s to %s: %v", ossObject.Key, localFile.AbsPath, err)
	}

	defer semaphore.Done()
	downloadLogger(ossBucket, ossObject, objectStatus)
}

// Download is starting the process of downloading all objects from an
// OSS Bucket.
func Download() {
	createDestinationDir := viper.GetBool("createDestinationDir")
	destinationDir := viper.GetString("destinationDir")
	if _, err := os.Stat(destinationDir); err != nil {
		if createDestinationDir == true {
			if err := os.MkdirAll(destinationDir, 0755); err != nil {
				ErrorLog.Printf("Cannot create destination directory %s: %v", destinationDir, err)
				os.Exit(1)
			}
		} else {
			ErrorLog.Printf("Local destination directory %s is not existing: %v", destinationDir, err)
			os.Exit(1)
		}
	}

	accessKeyID := viper.GetString("accessKeyID")
	accessKeySecret := viper.GetString("accessKeySecret")
	region := viper.GetString("region")
	client, err := oss.New(OssEndpoint(region), accessKeyID, accessKeySecret)
	if err != nil {
		ErrorLog.Printf("Cannot connect to Aliyun OSS API: %v", err)
		os.Exit(1)
	}

	bucket := viper.GetString("bucket")
	ossBucket, err := client.Bucket(bucket)
	if err != nil {
		ErrorLog.Printf("Cannot get the specified bucket %s instance: %v", bucket, err)
		os.Exit(1)
	}

	exclude := viper.GetStringSlice("exclude")
	marker := viper.GetString("marker")
	maxKeys := viper.GetInt("maxKeys")
	prefix := viper.GetString("prefix")
	threads := viper.GetInt("threads")
	semaphore := golimit.New(threads)
	for {
		resp, err := ossBucket.ListObjects(oss.Prefix(prefix), oss.Marker(marker), oss.MaxKeys(maxKeys))
		if err != nil {
			ErrorLog.Printf("Cannot list objects in bucket %s: %v", bucket, err)
		}
		for _, ossObject := range resp.Objects {
			excludeObject := false
			if len(exclude) > 0 {
				for _, str := range exclude {
					excludeObject = strings.Contains(ossObject.Key, str)
					if excludeObject {
						downloadLogger(ossBucket, ossObject, statusExcluded)
						break
					}
				}
			}
			if excludeObject == false && ossObject.Size > 0 {
				localFile, _ := LocalFile(ossObject)
				switch localFile.Exists {
				case true:
					if ossObject.Size != localFile.Size {
						semaphore.Add(1)
						go downloadOssObject(semaphore, ossBucket, ossObject, statusUpdate)
					} else {
						downloadLogger(ossBucket, ossObject, statusSkip)
					}
				case false:
					semaphore.Add(1)
					go downloadOssObject(semaphore, ossBucket, ossObject, statusGet)
				default:
					downloadLogger(ossBucket, ossObject, statusUnknown)
				}
			}
		}
		if resp.IsTruncated {
			marker = resp.NextMarker
		} else {
			break
		}
	}
	semaphore.Wait()
}
