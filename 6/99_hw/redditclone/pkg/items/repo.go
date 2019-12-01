package items

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

/*
func (repo *PostsRepo) Add(post *Post) (uint32, error) {
	repo.lastID++
	post.Id = repo.lastID
	repo.data = append(repo.data, post)
	return repo.lastID, nil
}

func (repo *PostsRepo) Update(newItem *Post) (bool, error) {
	for _, post := range repo.data {
		if post.Id != newItem.Id {
			continue
		}
		post.Title = newItem.Title
		post.Text = newItem.Text
		return true, nil
	}
	return false, nil
}

/*
func (repo *PostsRepo) Delete(id uint32) (bool, error) {
	i := -1
	for idx, post := range repo.data {
		if post.Id != id {
			continue
		}
		i = idx
	}
	if i < 0 {
		return false, nil
	}

	if i < len(repo.data)-1 {
		copy(repo.data[i:], repo.data[i+1:])
	}
	repo.data[len(repo.data)-1] = nil // or the zero value of T
	repo.data = repo.data[:len(repo.data)-1]

	return true, nil
}
*/
