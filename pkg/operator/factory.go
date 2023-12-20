// Copyright 2023 The prometheus-operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package operator

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/ptr"
)

const (
	managedByOperatorLabel      = "managed-by"
	managedByOperatorLabelValue = "prometheus-operator"
)

type ObjectOption func(metav1.Object)

type Owner interface {
	metav1.ObjectMetaAccessor
	schema.ObjectKind
}

func WithOwner(owner Owner) ObjectOption {
	return func(o metav1.Object) {
		o.SetOwnerReferences(
			append(
				o.GetOwnerReferences(),
				metav1.OwnerReference{
					APIVersion:         owner.GroupVersionKind().GroupVersion().String(),
					BlockOwnerDeletion: ptr.To(true),
					Controller:         ptr.To(true),
					Kind:               owner.GroupVersionKind().Kind,
					Name:               owner.GetObjectMeta().GetName(),
					UID:                owner.GetObjectMeta().GetUID(),
				},
			),
		)
	}
}

func WithName(name string) ObjectOption {
	return func(o metav1.Object) {
		o.SetName(name)
	}
}

func WithLabels(labels map[string]string) ObjectOption {
	return func(o metav1.Object) {
		l := Map{}
		l = l.Merge(labels)
		l = l.Merge(o.GetLabels())

		o.SetLabels(l)
	}
}

func WithAnnotations(annotations map[string]string) ObjectOption {
	return func(o metav1.Object) {
		a := Map{}
		a = a.Merge(annotations)
		a = a.Merge(o.GetAnnotations())

		o.SetAnnotations(a)
	}
}

func UpdateObject(o metav1.Object, opts ...ObjectOption) {
	WithLabels(map[string]string{managedByOperatorLabel: managedByOperatorLabelValue})(o)

	for _, opt := range opts {
		opt(o)
	}
}
