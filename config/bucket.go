/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package config

import (
	"github.com/spf13/viper"
)

//Bucket represents the basic configuration for a bucket
type Bucket struct {
	MaxCPU     int64 `mapstructure:"bucketMaxCPU"`
	MaxMemory  int64 `mapstructure:"bucketMaxMemory"`
	MaxStorage int64 `mapstructure:"bucketMaxStorage"`

	MinCPU     int64 `mapstructure:"bucketMinCPU"`
	MinMemory  int64 `mapstructure:"bucketMinMemory"`
	MinStorage int64 `mapstructure:"bucketMinStorage"`

	UnitCPU     int64 `mapstructure:"bucketUnitCPU"`
	UnitMemory  int64 `mapstructure:"bucketUnitMemory"`
	UnitStorage int64 `mapstructure:"bucketUnitStorage"`

	MaxBuckets int64 `mapstructure:"maxBuckets"`
}

//NewBucket generates a Bucket configuration from the given viper
//  Configuration
func NewBucket(v *viper.Viper) (Bucket, error) {
	out := Bucket{}
	return out, v.Unmarshal(&out)
}

func setBucketBindings(v *viper.Viper) error {
	err := v.BindEnv("bucketMaxCPU", "BUCKET_MAX_CPU")
	if err != nil {
		return err
	}

	err = v.BindEnv("bucketMaxMemory", "BUCKET_MAX_MEMORY")
	if err != nil {
		return err
	}

	err = v.BindEnv("bucketMaxStorage", "BUCKET_MAX_STORAGE")
	if err != nil {
		return err
	}

	err = v.BindEnv("bucketMinCPU", "BUCKET_MIN_CPU")
	if err != nil {
		return err
	}

	err = v.BindEnv("bucketMinMemory", "BUCKET_MIN_MEMORY")
	if err != nil {
		return err
	}

	err = v.BindEnv("bucketMinStorage", "BUCKET_MIN_STORAGE")
	if err != nil {
		return err
	}

	err = v.BindEnv("bucketUnitCPU", "BUCKET_UNIT_CPU")
	if err != nil {
		return err
	}

	err = v.BindEnv("bucketUnitMemory", "BUCKET_UNIT_MEMORY")
	if err != nil {
		return err
	}

	err = v.BindEnv("bucketUnitStorage", "BUCKET_UNIT_STORAGE")
	if err != nil {
		return err
	}

	return v.BindEnv("maxBuckets", "MAX_BUCKETS")
}

func setBucketDefaults(v *viper.Viper) {
	v.SetDefault("bucketMaxCPU", 64)
	v.SetDefault("bucketMaxMemory", 128*1024)  // 128 GiB in MiB
	v.SetDefault("bucketMaxStorage", 375*1024) // 375 GiB in MiB
	v.SetDefault("bucketMinCPU", 1)
	v.SetDefault("bucketMinMemory", 1*1024)   // 1 GiB in MiB
	v.SetDefault("bucketMinStorage", 10*1024) // 10 GiB in MiB
	v.SetDefault("bucketUnitCPU", 1)
	v.SetDefault("bucketUnitMemory", 128)       // 128MiB
	v.SetDefault("bucketUnitStorage", 375*1024) // 375 GiB in MiB
	v.SetDefault("maxBuckets", 3000)            // Max 3000 instances
}
