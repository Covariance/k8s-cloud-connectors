// Copyright (c) 2021 Yandex LLC. All rights reserved.
// Author: Martynov Pavel <covariance@yandex-team.ru>

package controller

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"

	connectorsv1 "k8s-connectors/connector/yos/api/v1"
	"k8s-connectors/pkg/awsutils"
)

func (r *yandexObjectStorageReconciler) allocateResource(
	ctx context.Context, log logr.Logger, object *connectorsv1.YandexObjectStorage, key, secret string,
) error {
	log.V(1).Info("started")

	lst, err := r.adapter.List(ctx, key, secret)
	if err != nil {
		return fmt.Errorf("unable to list resources: %v", err)
	}
	for _, bucket := range lst {
		if *bucket.Name == object.Name {
			log.V(1).Info("bucket found")
			return nil
		}
	}

	err = r.adapter.Create(ctx, key, secret, object.Spec.Name)
	if err != nil {
		return fmt.Errorf("unable to create resource: %v", err)
	}
	log.Info("successful")
	return nil
}

func (r *yandexObjectStorageReconciler) deallocateResource(
	ctx context.Context, log logr.Logger, object *connectorsv1.YandexObjectStorage, key, secret string,
) error {
	log.V(1).Info("started")

	err := r.adapter.Delete(ctx, key, secret, object.Spec.Name)
	if err != nil {
		if awsutils.CheckS3DoesNotExist(err) {
			log.Info("already deleted")
			return nil
		}
		return fmt.Errorf("unable to delete resource: %v", err)
	}

	log.Info("successful")
	return nil
}
