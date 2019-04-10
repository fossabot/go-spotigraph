package services

import (
	"context"

	"go.zenithar.org/spotigraph/pkg/protocol/v1/spotigraph"
)

//go:generate mockgen -destination test/mock/user.gen.go -package mock go.zenithar.org/spotigraph/internal/services User

// User defines user service contract
type User interface {
	Create(ctx context.Context, req *spotigraph.UserCreateReq) (*spotigraph.SingleUserRes, error)
	Get(ctx context.Context, req *spotigraph.UserGetReq) (*spotigraph.SingleUserRes, error)
	Update(ctx context.Context, req *spotigraph.UserUpdateReq) (*spotigraph.SingleUserRes, error)
	Delete(ctx context.Context, req *spotigraph.UserGetReq) (*spotigraph.EmptyRes, error)
	Search(ctx context.Context, req *spotigraph.UserSearchReq) (*spotigraph.PaginatedUserRes, error)
}

//go:generate mockgen -destination test/mock/squad.gen.go -package mock go.zenithar.org/spotigraph/internal/services Squad

// Squad defines squad service contract
type Squad interface {
	Create(ctx context.Context, req *spotigraph.SquadCreateReq) (*spotigraph.SingleSquadRes, error)
	Get(ctx context.Context, req *spotigraph.SquadGetReq) (*spotigraph.SingleSquadRes, error)
	Update(ctx context.Context, req *spotigraph.SquadUpdateReq) (*spotigraph.SingleSquadRes, error)
	Delete(ctx context.Context, req *spotigraph.SquadGetReq) (*spotigraph.EmptyRes, error)
	Search(ctx context.Context, req *spotigraph.SquadSearchReq) (*spotigraph.PaginatedSquadRes, error)
}

//go:generate mockgen -destination test/mock/chapter.gen.go -package mock go.zenithar.org/spotigraph/internal/services Chapter

// Chapter defines chapter service contract
type Chapter interface {
	Create(ctx context.Context, req *spotigraph.ChapterCreateReq) (*spotigraph.SingleChapterRes, error)
	Get(ctx context.Context, req *spotigraph.ChapterGetReq) (*spotigraph.SingleChapterRes, error)
	Update(ctx context.Context, req *spotigraph.ChapterUpdateReq) (*spotigraph.SingleChapterRes, error)
	Delete(ctx context.Context, req *spotigraph.ChapterGetReq) (*spotigraph.EmptyRes, error)
	Search(ctx context.Context, req *spotigraph.ChapterSearchReq) (*spotigraph.PaginatedChapterRes, error)
}

//go:generate mockgen -destination test/mock/guild.gen.go -package mock go.zenithar.org/spotigraph/internal/services Guild

// Guild defines guild service contract
type Guild interface {
	Create(ctx context.Context, req *spotigraph.GuildCreateReq) (*spotigraph.SingleGuildRes, error)
	Get(ctx context.Context, req *spotigraph.GuildGetReq) (*spotigraph.SingleGuildRes, error)
	Update(ctx context.Context, req *spotigraph.GuildUpdateReq) (*spotigraph.SingleGuildRes, error)
	Delete(ctx context.Context, req *spotigraph.GuildGetReq) (*spotigraph.EmptyRes, error)
	Search(ctx context.Context, req *spotigraph.GuildSearchReq) (*spotigraph.PaginatedGuildRes, error)
}

//go:generate mockgen -destination test/mock/tribe.gen.go -package mock go.zenithar.org/spotigraph/internal/services Tribe

// Tribe defines tribe service contract
type Tribe interface {
	Create(ctx context.Context, req *spotigraph.TribeCreateReq) (*spotigraph.SingleTribeRes, error)
	Get(ctx context.Context, req *spotigraph.TribeGetReq) (*spotigraph.SingleTribeRes, error)
	Update(ctx context.Context, req *spotigraph.TribeUpdateReq) (*spotigraph.SingleTribeRes, error)
	Delete(ctx context.Context, req *spotigraph.TribeGetReq) (*spotigraph.EmptyRes, error)
	Search(ctx context.Context, req *spotigraph.TribeSearchReq) (*spotigraph.PaginatedTribeRes, error)
}

//go:generate mockgen -destination test/mock/graph.gen.go -package mock go.zenithar.org/spotigraph/internal/services Graph

// Graph defines graph service contract
type Graph interface {
	Expand(cx context.Context, req *spotigraph.NodeInfoReq) (*spotigraph.GraphRes, error)
}
