package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/humbruno/glottr/internal/storage"
)

var usernames = []string{
	"alice", "bob", "charlie", "dave", "eve", "frank", "grace", "heidi",
	"ivan", "judy", "karl", "laura", "mallory", "nina", "oscar", "peggy",
	"quinn", "rachel", "steve", "trent", "ursula", "victor", "wendy", "xander",
	"yvonne", "zack", "amber", "brian", "carol", "doug", "eric", "fiona",
	"george", "hannah", "ian", "jessica", "kevin", "lisa", "mike", "natalie",
	"oliver", "peter", "queen", "ron", "susan", "tim", "uma", "vicky",
	"walter", "xenia", "yasmin", "zoe",
}

func Seed(storage storage.Storage, db *sql.DB) {
	ctx := context.Background()

	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)

	for _, user := range users {
		if err := storage.Users.CreateUser(ctx, tx, user.Email, user.Username); err != nil {
			_ = tx.Rollback()
			slog.Error("Error creating user", "err", err)
			return
		}
	}

	tx.Commit()
}

func generateUsers(num int) []*storage.User {
	users := make([]*storage.User, num)

	for i := 0; i < num; i++ {
		users[i] = &storage.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
		}
	}

	return users
}
