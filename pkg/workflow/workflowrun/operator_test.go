package workflowrun

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/caicloud/cyclone/pkg/apis/cyclone/v1alpha1"
	"github.com/caicloud/cyclone/pkg/util/k8s/fake"
)

func TestOverallStatus(t *testing.T) {
	client := fake.NewSimpleClientset()
	recorder := new(MockedRecorder)
	recorder.On("Event", mock.Anything).Return()
	wf := &v1alpha1.Workflow{}
	wfr := &v1alpha1.WorkflowRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
	}
	o := &operator{
		client:   client,
		recorder: recorder,
		wf:       wf,
		wfr:      wfr,
	}
	overall, _ := o.OverallStatus()
	assert.Equal(t, v1alpha1.StatusPending, overall.Phase)

	wf = &v1alpha1.Workflow{
		Spec: v1alpha1.WorkflowSpec{
			Stages: []v1alpha1.StageItem{
				{
					Name: "A",
				},
				{
					Name: "B",
				},
			},
		},
	}
	wfr = &v1alpha1.WorkflowRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
		Status: v1alpha1.WorkflowRunStatus{
			Stages: map[string]*v1alpha1.StageStatus{
				"A": {
					Status: v1alpha1.Status{Phase: v1alpha1.StatusSucceeded},
				},
				"B": {
					Status: v1alpha1.Status{Phase: v1alpha1.StatusPending},
				},
			},
		},
	}
	o = &operator{
		client:   client,
		recorder: recorder,
		wf:       wf,
		wfr:      wfr,
	}
	overall, _ = o.OverallStatus()
	assert.Equal(t, v1alpha1.StatusRunning, overall.Phase)

	wfr = &v1alpha1.WorkflowRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
		Status: v1alpha1.WorkflowRunStatus{
			Stages: map[string]*v1alpha1.StageStatus{
				"A": {
					Status: v1alpha1.Status{Phase: v1alpha1.StatusSucceeded},
				},
				"B": {
					Status: v1alpha1.Status{Phase: v1alpha1.StatusSucceeded},
				},
			},
		},
	}
	o = &operator{
		client:   client,
		recorder: recorder,
		wf:       wf,
		wfr:      wfr,
	}
	overall, _ = o.OverallStatus()
	assert.Equal(t, v1alpha1.StatusSucceeded, overall.Phase)

	wfr = &v1alpha1.WorkflowRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
		Status: v1alpha1.WorkflowRunStatus{
			Stages: map[string]*v1alpha1.StageStatus{
				"A": {
					Status: v1alpha1.Status{Phase: v1alpha1.StatusSucceeded},
				},
				"B": {
					Status: v1alpha1.Status{Phase: v1alpha1.StatusFailed},
				},
			},
		},
	}
	o = &operator{
		client:   client,
		recorder: recorder,
		wf:       wf,
		wfr:      wfr,
	}
	overall, _ = o.OverallStatus()
	assert.Equal(t, v1alpha1.StatusFailed, overall.Phase)

	wfr = &v1alpha1.WorkflowRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
		Status: v1alpha1.WorkflowRunStatus{
			Stages: map[string]*v1alpha1.StageStatus{
				"A": {
					Status: v1alpha1.Status{Phase: v1alpha1.StatusFailed},
				},
				"B": {
					Status: v1alpha1.Status{Phase: v1alpha1.StatusRunning},
				},
			},
		},
	}
	o = &operator{
		client:   client,
		recorder: recorder,
		wf:       wf,
		wfr:      wfr,
	}
	overall, _ = o.OverallStatus()
	assert.Equal(t, v1alpha1.StatusRunning, overall.Phase)
}
