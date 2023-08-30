/*
  Copyright The Kubernetes Authors.

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/srl-labs/clabernetes/apis/topology/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeKnes implements KneInterface
type FakeKnes struct {
	Fake *FakeTopologyV1alpha1
	ns   string
}

var knesResource = v1alpha1.SchemeGroupVersion.WithResource("knes")

var knesKind = v1alpha1.SchemeGroupVersion.WithKind("Kne")

// Get takes name of the kne, and returns the corresponding kne object, and an error if there is any.
func (c *FakeKnes) Get(
	ctx context.Context,
	name string,
	options v1.GetOptions,
) (result *v1alpha1.Kne, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(knesResource, c.ns, name), &v1alpha1.Kne{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Kne), err
}

// List takes label and field selectors, and returns the list of Knes that match those selectors.
func (c *FakeKnes) List(
	ctx context.Context,
	opts v1.ListOptions,
) (result *v1alpha1.KneList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(knesResource, knesKind, c.ns, opts), &v1alpha1.KneList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.KneList{ListMeta: obj.(*v1alpha1.KneList).ListMeta}
	for _, item := range obj.(*v1alpha1.KneList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested knes.
func (c *FakeKnes) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(knesResource, c.ns, opts))

}

// Create takes the representation of a kne and creates it.  Returns the server's representation of the kne, and an error, if there is any.
func (c *FakeKnes) Create(
	ctx context.Context,
	kne *v1alpha1.Kne,
	opts v1.CreateOptions,
) (result *v1alpha1.Kne, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(knesResource, c.ns, kne), &v1alpha1.Kne{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Kne), err
}

// Update takes the representation of a kne and updates it. Returns the server's representation of the kne, and an error, if there is any.
func (c *FakeKnes) Update(
	ctx context.Context,
	kne *v1alpha1.Kne,
	opts v1.UpdateOptions,
) (result *v1alpha1.Kne, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(knesResource, c.ns, kne), &v1alpha1.Kne{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Kne), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeKnes) UpdateStatus(
	ctx context.Context,
	kne *v1alpha1.Kne,
	opts v1.UpdateOptions,
) (*v1alpha1.Kne, error) {
	obj, err := c.Fake.
		Invokes(
			testing.NewUpdateSubresourceAction(knesResource, "status", c.ns, kne),
			&v1alpha1.Kne{},
		)

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Kne), err
}

// Delete takes name of the kne and deletes it. Returns an error if one occurs.
func (c *FakeKnes) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(knesResource, c.ns, name, opts), &v1alpha1.Kne{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeKnes) DeleteCollection(
	ctx context.Context,
	opts v1.DeleteOptions,
	listOpts v1.ListOptions,
) error {
	action := testing.NewDeleteCollectionAction(knesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.KneList{})
	return err
}

// Patch applies the patch and returns the patched kne.
func (c *FakeKnes) Patch(
	ctx context.Context,
	name string,
	pt types.PatchType,
	data []byte,
	opts v1.PatchOptions,
	subresources ...string,
) (result *v1alpha1.Kne, err error) {
	obj, err := c.Fake.
		Invokes(
			testing.NewPatchSubresourceAction(knesResource, c.ns, name, pt, data, subresources...),
			&v1alpha1.Kne{},
		)

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Kne), err
}
