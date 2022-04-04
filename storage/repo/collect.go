package repo

type CollectStorage interface {
	CollectPostsStart() error
	CollectPostsFinish() error
	CheckFinished() (bool, error)
	CheckStarted() (bool, error)
}
