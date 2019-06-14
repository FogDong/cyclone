/*
Copyright 2019 caicloud authors. All rights reserved.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/caicloud/cyclone/pkg/apis/cyclone/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeWorkflowTriggers implements WorkflowTriggerInterface
type FakeWorkflowTriggers struct {
	Fake *FakeCycloneV1alpha1
	ns   string
}

var workflowtriggersResource = schema.GroupVersionResource{Group: "cyclone.dev", Version: "v1alpha1", Resource: "workflowtriggers"}

var workflowtriggersKind = schema.GroupVersionKind{Group: "cyclone.dev", Version: "v1alpha1", Kind: "WorkflowTrigger"}

// Get takes name of the workflowTrigger, and returns the corresponding workflowTrigger object, and an error if there is any.
func (c *FakeWorkflowTriggers) Get(name string, options v1.GetOptions) (result *v1alpha1.WorkflowTrigger, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(workflowtriggersResource, c.ns, name), &v1alpha1.WorkflowTrigger{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.WorkflowTrigger), err
}

// List takes label and field selectors, and returns the list of WorkflowTriggers that match those selectors.
func (c *FakeWorkflowTriggers) List(opts v1.ListOptions) (result *v1alpha1.WorkflowTriggerList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(workflowtriggersResource, workflowtriggersKind, c.ns, opts), &v1alpha1.WorkflowTriggerList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.WorkflowTriggerList{}
	for _, item := range obj.(*v1alpha1.WorkflowTriggerList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested workflowTriggers.
func (c *FakeWorkflowTriggers) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(workflowtriggersResource, c.ns, opts))

}

// Create takes the representation of a workflowTrigger and creates it.  Returns the server's representation of the workflowTrigger, and an error, if there is any.
func (c *FakeWorkflowTriggers) Create(workflowTrigger *v1alpha1.WorkflowTrigger) (result *v1alpha1.WorkflowTrigger, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(workflowtriggersResource, c.ns, workflowTrigger), &v1alpha1.WorkflowTrigger{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.WorkflowTrigger), err
}

// Update takes the representation of a workflowTrigger and updates it. Returns the server's representation of the workflowTrigger, and an error, if there is any.
func (c *FakeWorkflowTriggers) Update(workflowTrigger *v1alpha1.WorkflowTrigger) (result *v1alpha1.WorkflowTrigger, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(workflowtriggersResource, c.ns, workflowTrigger), &v1alpha1.WorkflowTrigger{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.WorkflowTrigger), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeWorkflowTriggers) UpdateStatus(workflowTrigger *v1alpha1.WorkflowTrigger) (*v1alpha1.WorkflowTrigger, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(workflowtriggersResource, "status", c.ns, workflowTrigger), &v1alpha1.WorkflowTrigger{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.WorkflowTrigger), err
}

// Delete takes name of the workflowTrigger and deletes it. Returns an error if one occurs.
func (c *FakeWorkflowTriggers) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(workflowtriggersResource, c.ns, name), &v1alpha1.WorkflowTrigger{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeWorkflowTriggers) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(workflowtriggersResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.WorkflowTriggerList{})
	return err
}

// Patch applies the patch and returns the patched workflowTrigger.
func (c *FakeWorkflowTriggers) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.WorkflowTrigger, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(workflowtriggersResource, c.ns, name, pt, data, subresources...), &v1alpha1.WorkflowTrigger{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.WorkflowTrigger), err
}
