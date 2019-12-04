package items

import "errors"

// WARNING! completly unsafe in multi-goroutine use, need add mutexes

type PostsRepo struct {
	lastID uint32
	data   []*Post
}

func NewPostRepo() *PostsRepo {
	return &PostsRepo{
		data: make([]*Post, 0, 10),
	}
}

func (repo *PostsRepo) Add(p *Post) {
	p.Id = repo.lastID
	repo.data = append(repo.data, p)
	repo.lastID++
}

func (repo *PostsRepo) GetAll() ([]*Post, error) {
	return repo.data, nil
}

func (repo *PostsRepo) GetByID(id uint32) (*Post, error) {
	for _, post := range repo.data {
		if post.Id == id {
			return post, nil
		}
	}
	return nil, nil
}

func (repo *PostsRepo) GetUserPosts(login string) ([]*Post, error) {
	ps, err := repo.GetAll()
	if err != nil {
		return nil, errors.New("db error")
	}
	ups := make([]*Post, 0, 10)
	for _, p := range ps {
		if p.Author.Username == login {
			ups = append(ups, p)
		}
	}
	return ups, nil
}

func (repo *PostsRepo) Delete(id uint32) error {
	for i, v := range repo.data {
		if v.Id == id {
			lh := len(repo.data)
			repo.data[lh-1], repo.data[i] = repo.data[i], repo.data[lh-1]
			repo.data = repo.data[:lh-1]
			return nil
		}
	}
	return errors.New("not found post")
}
