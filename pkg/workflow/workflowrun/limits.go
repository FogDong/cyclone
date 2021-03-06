package workflowrun

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/caicloud/cyclone/pkg/apis/cyclone/v1alpha1"
	"github.com/caicloud/cyclone/pkg/k8s/clientset"
	"github.com/caicloud/cyclone/pkg/workflow/controller"
)

// LimitedQueues manages WorkflowRun queue for each Workflow. Queue for each Workflow is limited to
// a given maximum size, if new WorkflowRun created, the oldest one would be removed.
type LimitedQueues struct {
	// Maximum queue size, it indicates maximum number of WorkflowRuns to retain for each Workflow.
	MaxQueueSize int
	// Workflow queue map. It use Workflow name and namespace as the key, and manage Queue for each
	// Workflow.
	Queues map[string]*LimitedSortedQueue
	// k8s client used to clean old WorkflowRun
	Client clientset.Interface
}

// NewLimitedQueues creates a limited queues for WorkflowRuns, and start auto scan.
func NewLimitedQueues(client clientset.Interface, maxSize int) *LimitedQueues {
	log.WithField("max", maxSize).Info("Create limited queues")
	queues := &LimitedQueues{
		MaxQueueSize: maxSize,
		Queues:       make(map[string]*LimitedSortedQueue),
		Client:       client,
	}
	go queues.AutoScan()
	return queues
}

func key(wfr *v1alpha1.WorkflowRun) string {
	return fmt.Sprintf("%s/%s", wfr.Spec.WorkflowRef.Namespace, wfr.Spec.WorkflowRef.Name)
}

// Refresh refreshes the WorkflowRun in the queue, the refresh time would be updated.
func (w *LimitedQueues) Refresh(wfr *v1alpha1.WorkflowRun) {
	q, ok := w.Queues[key(wfr)]
	if !ok {
		log.WithField("key", key(wfr)).Warn("Queue not exist")
		return
	}

	q.Refresh(wfr)
}

// AddOrRefresh adds a WorkflowRun to its corresponding queue, if the queue size exceed the maximum size, the
// oldest one would be deleted. And if the WorkflowRun already exists in the queue, its 'refresh' time field
// would be refreshed.
func (w *LimitedQueues) AddOrRefresh(wfr *v1alpha1.WorkflowRun) {
	q, ok := w.Queues[key(wfr)]
	if !ok {
		q = NewQueue(key(wfr), w.MaxQueueSize)
		w.Queues[key(wfr)] = q
	}

	// PushOrRefresh push the WorkflowRun to the queue. If it's already existed in the queue, its refresh
	// time would be updated to now.
	q.PushOrRefresh(wfr)

	for q.size > w.MaxQueueSize {
		log.WithField("max", w.MaxQueueSize).Debug("Max WorkflowRun exceeded, delete the oldest one")
		old := q.Pop()
		err := w.Client.CycloneV1alpha1().WorkflowRuns(old.namespace).Delete(context.TODO(), old.wfr, metav1.DeleteOptions{})
		if err != nil && !errors.IsNotFound(err) {
			log.WithField("wfr", old.wfr).Error("Delete old WorkflowRun error: ", err)
		} else {
			log.WithField("wfr", old.wfr).Info("Old WorkflowRun deleted")
		}
	}
}

// AutoScan scans all WorkflowRuns in the queues regularly, remove abnormal ones with old enough
// refresh time.
func (w *LimitedQueues) AutoScan() {
	ticker := time.NewTicker(time.Hour)
	for range ticker.C {
		for _, q := range w.Queues {
			scanQueue(q)
		}
	}
}

// scanQueue scans all WorkflowRun in the queue, check their refresh time with current time, if refresh
// time is old enough, it means WorkflowRun is actually deleted in k8s, but somehow Workflow Controller
// didn't know (this seldom happen), in this case, remove the WorkflowRun from the queue.
func scanQueue(q *LimitedSortedQueue) {
	q.lock.Lock()
	defer q.lock.Unlock()

	h := q.head
	for h.next != nil {
		// If the node's refresh time is old enough compared to the resync time
		// (5 minutes by default) of WorkflowRun Controller, it means the WorkflowRun
		// is actually removed from etcd somehow, so we will remove it also here.
		if h.next.refresh.Add(controller.Config.ResyncPeriodSeconds * time.Second * 2).Before(time.Now()) {
			log.WithField("wfr", h.next.wfr).Info("remove wfr with outdated refresh time from queue")
			h.next = h.next.next
			q.size--
			continue
		}

		h = h.next
	}
}

// LimitedSortedQueue is a sorted fixed length queue implemented with single linked list.
// Note that each queue would have a sentinel node to assist the implementation, it's a
// dummy node, and won't be counted in the queue size. So an empty queue would have head
// pointed to dummy node, with queue size 0.
type LimitedSortedQueue struct {
	// Key of the Workflow, it's generated by namespace and workflow name
	key string
	// Lock to for concurrency control
	lock sync.Mutex
	// Maximum queue size
	max int
	// Current size of the queue
	size int
	// Head of the queue
	head *Node
}

// NewQueue creates a limited sorted queue.
func NewQueue(key string, max int) *LimitedSortedQueue {
	dummy := &Node{}
	return &LimitedSortedQueue{
		key:  key,
		max:  max,
		size: 0,
		head: dummy,
	}
}

// Node represents a WorkflowRun in the queue. The 'next' link to next node in the queue, and the
// 'refresh' stands the last time this node is refreshed.
//
// 'refresh' here is used to deal with a rarely occurred case: when one WorkflowRun got deleted in
// etcd, but workflow controller didn't get notified. Workflow controller would perform resync with
// etcd regularly, (5 minutes by default), during resync every WorkflowRun in the queue would be
// refreshed, it one WorkflowRun is deleted in etcd abnormally, it wouldn't get refreshed in the queue,
// so we can judge by the refresh time for this case.
//
// When we found a node that hasn't be refreshed for a long time (for example, twice the resync period),
// then we remove this node from the queue.
type Node struct {
	next *Node
	// Name of the WorkflowRun
	wfr string
	// Namespace of the WorkflowRun
	namespace string
	// Time when the WorkflowRun is created
	created int64
	// Time when the node is refreshed
	refresh time.Time
}

// PushOrRefresh pushes a WorkflowRun object to the queue, it will be inserted in the right place to keep
// the queue sorted by creation time.
// If the object already existed in the queue, its refresh time would be updated.
func (q *LimitedSortedQueue) PushOrRefresh(wfr *v1alpha1.WorkflowRun) {
	q.lock.Lock()
	defer q.lock.Unlock()

	node := &Node{
		wfr:       wfr.Name,
		namespace: wfr.Namespace,
		created:   wfr.ObjectMeta.CreationTimestamp.Time.Unix(),
		refresh:   time.Now(),
	}

	if q.Refresh(wfr) {
		return
	}

	p := q.head
	for p.next != nil && p.next.created < node.created {
		p = p.next
	}

	node.next = p.next
	p.next = node
	q.size++
}

// Refresh updates refresh time of WorkflowRun in the queue, if the WorkflowRun found in the queue
// and update successfully, return true, otherwise return false.
func (q *LimitedSortedQueue) Refresh(wfr *v1alpha1.WorkflowRun) bool {
	p := q.head
	for p.next != nil && (p.next.namespace != wfr.Namespace || p.next.wfr != wfr.Name) {
		p = p.next
	}

	if p.next != nil {
		log.WithField("queue", q.key).WithField("wfr", wfr.Name).Debug("Update refresh time")
		p.next.refresh = time.Now()
		return true
	}

	return false
}

// Pop pops up a WorkflowRun object from the queue, it's the oldest one that will be popped.
func (q *LimitedSortedQueue) Pop() *Node {
	if q.size <= 0 {
		return nil
	}

	q.lock.Lock()
	defer q.lock.Unlock()

	n := q.head.next
	q.head.next = q.head.next.next
	q.size--
	return n
}
