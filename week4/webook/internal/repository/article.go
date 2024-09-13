package repository

import (
	"context"
	"gorm.io/gorm"
	"time"
	"webook/internal/domain"
	"webook/internal/repository/dao"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	Sync(ctx context.Context, art domain.Article) (int64, error)
	//Sync(ctx context.Context, art domain.Article) (int64, error)
	//SyncStatus(ctx context.Context, uid int64, id int64, status domain.ArticleStatus) error
	//GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	//GetById(ctx context.Context, id int64) (domain.Article, error)
	//GetPubById(ctx context.Context, id int64) (domain.Article, error)
}

type CachedArticleRepository struct {
	dao dao.ArticleDAO
	//cache cache.ArticleCache
	//// 因为如果你直接访问 UserDAO，你就绕开了 repository，
	//// repository 一般都有一些缓存机制
	//userRepo UserRepository
	//
	readerDAO dao.ArticleReaderDAO
	authorDAO dao.ArticleAuthorDAO
	//
	db *gorm.DB
}

func NewCachedArticleRepository(dao dao.ArticleDAO) ArticleRepository {
	return &CachedArticleRepository{
		dao: dao,
	}
}

func NewCachedArticleRepositoryV2(
	readerDAO dao.ArticleReaderDAO,
	authorDAO dao.ArticleAuthorDAO) *CachedArticleRepository {
	return &CachedArticleRepository{
		readerDAO: readerDAO,
		authorDAO: authorDAO,
	}
}

//func (c *CachedArticleRepository) Update(ctx context.Context, art domain.Article) error {
//	err := c.dao.UpdateById(ctx, c.toEntity(art))
//	if err == nil {
//		er := c.cache.DelFirstPage(ctx, art.Author.Id)
//		if er != nil {
//			// 也要记录日志
//		}
//	}
//	return err
//}

func (c *CachedArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	id, err := c.dao.Insert(ctx, c.toEntity(art))
	//if err == nil {
	//	er := c.cache.DelFirstPage(ctx, art.Author.Id)
	//	if er != nil {
	//		// 也要记录日志
	//	}
	//}
	return id, err
}

func (c *CachedArticleRepository) Update(ctx context.Context, art domain.Article) error {
	return c.dao.UpdateById(ctx, c.toEntity(art))
}

func (c *CachedArticleRepository) SyncV1(ctx context.Context, art domain.Article) (int64, error) {
	artn := c.toEntity(art)
	var (
		id  = art.Id
		err error
	)
	if id > 0 {
		err = c.authorDAO.Update(ctx, artn)
	} else {
		id, err = c.authorDAO.Create(ctx, artn)
	}
	if err != nil {
		return 0, err
	}
	artn.Id = id
	err = c.readerDAO.Upsert(ctx, artn)
	return id, err
}

func (c *CachedArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	tx := c.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	// 防止后面业务panic
	defer tx.Rollback()

	authorDAO := dao.NewArticleGORMAuthorDAO(tx)
	readerDAO := dao.NewArticleGORMReaderDAO(tx)

	artn := c.toEntity(art)
	var (
		id  = art.Id
		err error
	)
	if id > 0 {
		err = authorDAO.Update(ctx, artn)
	} else {
		id, err = authorDAO.Create(ctx, artn)
	}
	if err != nil {
		return 0, err
	}
	artn.Id = id
	err = readerDAO.UpsertV2(ctx, dao.PublishedArticle(artn))
	if err != nil {
		return 0, err
	}
	tx.Commit()
	return id, nil
}

func (c *CachedArticleRepository) toEntity(art domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		//Status:   uint8(art.Status),
		Status: art.Status.ToUint8(),
	}
}

func (c *CachedArticleRepository) toDomain(art dao.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			// 这里有一个错误
			Id: art.AuthorId,
		},
		Ctime:  time.UnixMilli(art.Ctime),
		Utime:  time.UnixMilli(art.Utime),
		Status: domain.ArticleStatus(art.Status),
	}
}
