package postgres

import (
	"fmt"
	"log"
	"time"

	"github.com/new-york-services/post_service/storage/repo"
)

func (r *PostRepo) CreatePost(post repo.Post) (repo.Post, error) {
	var res repo.Post
	err := r.db.QueryRow(`
		insert into 
			posts(id, title, description, user_id) 
		values
			($1, $2, $3, $4) 
		returning 
			id, title, description, likes, user_id, created_at, updated_at`, post.Id, post.Title, post.Description, post.UserId).Scan(&res.Id, &res.Title, &res.Description, &res.Likes, &res.UserId, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		log.Println("failed to create post in sql: ", err)
		return repo.Post{}, err
	}

	return res, nil
}

func (r *PostRepo) GetPostById(id string) (repo.Post, error) {
	res := repo.Post{}
	err := r.db.QueryRow(`
		select 
			id, title, description, likes, user_id, created_at, updated_at 
		from 
			posts 
		where 
			id = $1 and deleted_at is null`, id).Scan(&res.Id, &res.Title, &res.Description, &res.Likes, &res.UserId, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		log.Println("failed to get post in sql: ", err)
		return repo.Post{}, err
	}

	return res, nil
}

func (r *PostRepo) GetPostByUserId(id string) ([]repo.Post, error) {
	res := []repo.Post{}
	rows, err := r.db.Query(`
		select 
			id, title, description, likes, user_id, created_at, updated_at 
		from 
			posts 
		where 
			user_id = $1 and deleted_at is null`, id)

	if err != nil {
		log.Println("failed to get post by user id in sql: ", err)
		return []repo.Post{}, err
	}

	for rows.Next() {
		post := repo.Post{}

		err = rows.Scan(
			&post.Id,
			&post.Title,
			&post.Description,
			&post.Likes,
			&post.UserId,
			&post.CreatedAt,
			&post.UpdatedAt,
		)

		if err != nil {
			log.Println("failed to scan post by user id in sql: ", err)
			return []repo.Post{}, err
		}

		res = append(res, post)
	}

	return res, nil
}

func (r *PostRepo) GetPostForUser(id string) ([]repo.Post, error) {
	res := []repo.Post{}
	rows, err := r.db.Query(`
		select 
			id, title, description, likes, user_id, created_at, updated_at 
		from 
			posts 
		where
			user_id = $1 and deleted_at is null`, id)

	if err != nil {
		log.Println("failed to get post for user in sql: ", err)
		return []repo.Post{}, err
	}

	for rows.Next() {
		post := repo.Post{}

		err = rows.Scan(
			&post.Id,
			&post.Title,
			&post.Description,
			&post.Likes,
			&post.UserId,
			&post.CreatedAt,
			&post.UpdatedAt,
		)

		if err != nil {
			log.Println("failed to scan post: in sql: ", err)
			return []repo.Post{}, nil
		}

		res = append(res, post)
	}

	return res, nil
}

func (r *PostRepo) GetPostForComment(id string) (repo.Post, error) {
	res := repo.Post{}
	err := r.db.QueryRow(`
		select 
			id, title, description, likes, user_id, created_at, updated_at 
		from 
			posts 
		where 
			id = $1 and deleted_at is null`, id).Scan(&res.Id, &res.Title, &res.Description, &res.Likes, &res.UserId, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		log.Println("failed to get post in sql: ", err)
		return repo.Post{}, err
	}

	return res, nil
}

func (r *PostRepo) SearchPosts(title string) ([]repo.Post, error) {
	res := []repo.Post{}
	query := fmt.Sprint("select id, title, description, likes, user_id, created_at, updated_at from posts where title ilike '%" + title + "%' and deleted_at is null")

	rows, err := r.db.Query(query)
	if err != nil {
		log.Println("failed to search post in sql: ", err)
		return []repo.Post{}, nil
	}

	for rows.Next() {
		post := repo.Post{}

		err = rows.Scan(
			&post.Id,
			&post.Title,
			&post.Description,
			&post.Likes,
			&post.UserId,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			log.Println("failed to scanning post in sql: ", err)
			return []repo.Post{}, nil
		}

		res = append(res, post)
	}

	return res, nil
}

func (r *PostRepo) LikePost(postId string, isLiked bool) (repo.Post, error) {
	res := repo.Post{}
	if isLiked {
		err := r.db.QueryRow(`
			update 
				posts 
			set 
				likes = likes + 1 
			where 
				id = $1 and deleted_at is null
			returning 
				id, title, description, likes, user_id, created_at, updated_at`, postId).Scan(&res.Id, &res.Title, &res.Description, &res.Likes, &res.UserId, &res.CreatedAt, &res.UpdatedAt)
		if err != nil {
			log.Println("failed to like post in sql: ", err)
			return repo.Post{}, err
		}
	} else {
		err := r.db.QueryRow(`
			select 
				id, title, description, likes + 1, user_id, created_at, updated_at 
			from 
				posts 
			where 
				id = $1 and deleted_at is null`, postId).Scan(&res.Id, &res.Title, &res.Description, &res.Likes, &res.UserId, &res.CreatedAt, &res.UpdatedAt)

		if err != nil {
			log.Println("failed to like post in sql: ", err)
			return repo.Post{}, err
		}
	}

	return res, nil
}

func (r *PostRepo) UpdatePost(post repo.Post) (repo.Post, error) {
	res := repo.Post{}
	err := r.db.QueryRow(`
		update
			posts 
		set 
			title = $1, description = $2, updated_at = $3 
		where 
			id = $4 and deleted_at is null
		returning id, title, description, likes, user_id, created_at, updated_at`, post.Title, post.Description, time.Now(), post.Id).Scan(&res.Id, &res.Title, &res.Description, &res.Likes, &res.UserId, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		log.Println("failed to update post in sql in sql: ", err)
		return repo.Post{}, err
	}

	return res, nil
}

func (r *PostRepo) DeletePost(id string) (repo.Post, error) {
	post := repo.Post{}
	err := r.db.QueryRow(`
		update 
			posts 
		set 
			deleted_at = $1 
		where 
			id = $2 and deleted_at is null
		returning 
			id, title, description, likes, user_id, created_at, updated_at`, time.Now(), id).Scan(&post.Id, &post.Title, &post.Description, &post.Likes, &post.UserId, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		log.Println("failed to delete post in sql: ", err)
		return repo.Post{}, err
	}

	return post, nil
}
