package catalog

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/yueli-official/foundation/go/identifier"

	"github.com/yueli-official/nav/api/internal/dao"
	"github.com/yueli-official/nav/api/internal/model"
	"github.com/yueli-official/nav/api/internal/naverr"
)

const (
	CheckJobRunning   = "running"
	CheckJobCompleted = "completed"
	CheckJobFailed    = "failed"

	retainedCheckJobs = 20
)

type CheckJob struct {
	ID         string
	Scope      string
	Status     string
	Total      int
	Completed  int
	StartedAt  time.Time
	FinishedAt *time.Time
	Error      string
}

type checkJobRegistry struct {
	mu       sync.Mutex
	jobs     map[string]*CheckJob
	order    []string
	activeID string
	now      func() time.Time
}

func newCheckJobRegistry() *checkJobRegistry {
	return &checkJobRegistry{jobs: map[string]*CheckJob{}, now: time.Now}
}

func (registry *checkJobRegistry) active() (CheckJob, bool) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	return registry.activeLocked()
}

func (registry *checkJobRegistry) start(scope string, total int) (CheckJob, bool) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	if active, ok := registry.activeLocked(); ok {
		return active, true
	}
	job := &CheckJob{
		ID: identifier.MustNew().String(), Scope: scope, Status: CheckJobRunning,
		Total: total, StartedAt: registry.now().UTC(),
	}
	registry.jobs[job.ID] = job
	registry.order = append(registry.order, job.ID)
	registry.activeID = job.ID
	registry.pruneLocked()
	return cloneCheckJob(job), false
}

func (registry *checkJobRegistry) activeLocked() (CheckJob, bool) {
	if registry.activeID == "" {
		return CheckJob{}, false
	}
	job := registry.jobs[registry.activeID]
	if job == nil || job.Status != CheckJobRunning {
		registry.activeID = ""
		return CheckJob{}, false
	}
	return cloneCheckJob(job), true
}

func (registry *checkJobRegistry) progress(id string) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	job := registry.jobs[id]
	if job == nil || job.Status != CheckJobRunning || job.Completed >= job.Total {
		return
	}
	job.Completed++
}

func (registry *checkJobRegistry) finish(id string, runErr error) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	job := registry.jobs[id]
	if job == nil || job.Status != CheckJobRunning {
		return
	}
	finishedAt := registry.now().UTC()
	job.FinishedAt = &finishedAt
	if runErr != nil {
		job.Status = CheckJobFailed
		job.Error = "检查任务未能完成，请稍后重试。"
	} else {
		job.Status = CheckJobCompleted
		job.Completed = job.Total
	}
	if registry.activeID == id {
		registry.activeID = ""
	}
}

func (registry *checkJobRegistry) get(id string) (CheckJob, bool) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	job := registry.jobs[strings.TrimSpace(id)]
	if job == nil {
		return CheckJob{}, false
	}
	return cloneCheckJob(job), true
}

func (registry *checkJobRegistry) pruneLocked() {
	for len(registry.order) > retainedCheckJobs {
		id := registry.order[0]
		registry.order = registry.order[1:]
		if id != registry.activeID {
			delete(registry.jobs, id)
		}
	}
}

func cloneCheckJob(job *CheckJob) CheckJob {
	if job == nil {
		return CheckJob{}
	}
	clone := *job
	if job.FinishedAt != nil {
		finishedAt := *job.FinishedAt
		clone.FinishedAt = &finishedAt
	}
	return clone
}

func (s *Service) StartSelectedCheckJob(ctx context.Context, ids []string) (CheckJob, bool, error) {
	if len(ids) > 50 {
		return CheckJob{}, false, naverr.Validation("ids", "maximum", map[string]any{"max": 50})
	}
	ids = normalize(ids, max(len(ids), 1))
	if len(ids) == 0 {
		return CheckJob{}, false, naverr.Validation("ids", "required", nil)
	}
	if active, ok := s.checkJobs.active(); ok {
		return active, true, nil
	}
	links, err := s.store.LinksByIDs(ctx, ids)
	if err != nil {
		return CheckJob{}, false, err
	}
	if len(links) != len(ids) {
		return CheckJob{}, false, naverr.NotFound("one_or_more_links")
	}
	return s.startCheckJob(ctx, "selected", checkableLinks(links))
}

func (s *Service) StartFilteredCheckJob(ctx context.Context, filter dao.LinkFilter) (CheckJob, bool, error) {
	if active, ok := s.checkJobs.active(); ok {
		return active, true, nil
	}
	filter.Page = 0
	filter.Size = 0
	filter.Sort = "health"
	filter.ExcludeHealthCheckExempt = true
	links, err := s.store.Links(ctx, filter)
	if err != nil {
		return CheckJob{}, false, err
	}
	return s.startCheckJob(ctx, "filtered", links)
}

func (s *Service) startCheckJob(ctx context.Context, scope string, links []*model.Link) (CheckJob, bool, error) {
	job, reused := s.checkJobs.start(scope, len(links))
	if reused {
		return job, true, nil
	}
	jobContext := context.WithoutCancel(ctx)
	s.checkJobRunner(func() {
		s.executeCheckJob(jobContext, job.ID, links)
	})
	return job, false, nil
}

func (s *Service) executeCheckJob(ctx context.Context, id string, links []*model.Link) {
	var runErr error
	defer func() {
		if recovered := recover(); recovered != nil {
			runErr = fmt.Errorf("panic while running navigation link checks: %v", recovered)
		}
		if runErr != nil {
			g.Log().Errorf(ctx, "navigation link check job %s failed: %v", id, runErr)
		}
		s.checkJobs.finish(id, runErr)
	}()
	_, runErr = s.runLinkChecksWithProgress(ctx, links, func() {
		s.checkJobs.progress(id)
	})
}

func (s *Service) CheckJob(id string) (CheckJob, error) {
	job, ok := s.checkJobs.get(id)
	if !ok {
		return CheckJob{}, naverr.NotFound(strings.TrimSpace(id))
	}
	return job, nil
}
