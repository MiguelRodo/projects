package git

import "context"

// MockClient is a mock implementation of Client for testing.
type MockClient struct {
	CloneFunc         func(ctx context.Context, url, targetDir, branch string) error
	PullFunc          func(ctx context.Context, repoDir string) (string, error)
	FetchFunc         func(ctx context.Context, repoDir string) error
	StatusFunc        func(ctx context.Context, repoDir string) (string, error)
	CurrentBranchFunc func(ctx context.Context, repoDir string) (string, error)
	IsCleanFunc       func(ctx context.Context, repoDir string) (bool, error)
	IsRepoFunc        func(path string) bool
}

func (m *MockClient) Clone(ctx context.Context, url, targetDir, branch string) error {
	if m.CloneFunc != nil {
		return m.CloneFunc(ctx, url, targetDir, branch)
	}
	return nil
}

func (m *MockClient) Pull(ctx context.Context, repoDir string) (string, error) {
	if m.PullFunc != nil {
		return m.PullFunc(ctx, repoDir)
	}
	return "Already up to date.", nil
}

func (m *MockClient) Fetch(ctx context.Context, repoDir string) error {
	if m.FetchFunc != nil {
		return m.FetchFunc(ctx, repoDir)
	}
	return nil
}

func (m *MockClient) Status(ctx context.Context, repoDir string) (string, error) {
	if m.StatusFunc != nil {
		return m.StatusFunc(ctx, repoDir)
	}
	return "", nil
}

func (m *MockClient) CurrentBranch(ctx context.Context, repoDir string) (string, error) {
	if m.CurrentBranchFunc != nil {
		return m.CurrentBranchFunc(ctx, repoDir)
	}
	return "main", nil
}

func (m *MockClient) IsClean(ctx context.Context, repoDir string) (bool, error) {
	if m.IsCleanFunc != nil {
		return m.IsCleanFunc(ctx, repoDir)
	}
	return true, nil
}

func (m *MockClient) IsRepo(path string) bool {
	if m.IsRepoFunc != nil {
		return m.IsRepoFunc(path)
	}
	return true
}
