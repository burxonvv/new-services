package postgres

import (
	"log"
	"time"

	"github.com/burxondv/new-services/comment-service/storage/repo"
)

func (r *CommentRepo) WriteComment(comment repo.Comment) (repo.Comment, error) {
	var res repo.Comment
	err := r.db.QueryRow(`
		insert into 
			comments(id, post_id, user_id, text)
		values
			($1, $2, $3, $4) 
		returning 
			id, post_id, user_id, text, created_at`, comment.Id, comment.PostId, comment.UserId, comment.Text).Scan(&res.Id, &res.PostId, &res.UserId, &res.Text, &res.CreatedAt)

	if err != nil {
		log.Println("failed to create comment in sql: ", err)
		return repo.Comment{}, err
	}

	return res, nil
}

func (r *CommentRepo) GetComments(id string) ([]repo.Comment, error) {
	var res []repo.Comment
	rows, err := r.db.Query(`
		select 
			id, post_id, user_id, text, created_at 
		from 
			comments 
		where 
			post_id = $1 and deleted_at is null`, id)

	if err != nil {
		log.Println("failed to get comment in sql: ", err)
		return []repo.Comment{}, nil
	}

	for rows.Next() {
		comment := repo.Comment{}

		err = rows.Scan(
			&comment.Id,
			&comment.PostId,
			&comment.UserId,
			&comment.Text,
			&comment.CreatedAt,
		)

		if err != nil {
			log.Println("failed to scanning comment in sql: ", err)
			return []repo.Comment{}, err
		}

		res = append(res, comment)
	}

	return res, nil
}

func (r *CommentRepo) DeleteComment(id string) (repo.Comment, error) {
	var res repo.Comment
	err := r.db.QueryRow(`
		update 
			comments 
		set 
			deleted_at = $1 
		where 
			id = $2 and deleted_at is null
		returning 
			id, post_id, user_id, text, created_at`, time.Now(), id).Scan(&res.Id, &res.PostId, &res.UserId, &res.Text, &res.CreatedAt)

	if err != nil {
		log.Println("failed to delete comment in sql", err)
	}

	return res, nil
}
