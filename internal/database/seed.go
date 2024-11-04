package database

import (
	"context"
	"database/sql"
	"fmt"

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

	users := generateUsers(1)

	for _, user := range users {
		storage.Users.Create(ctx, user)
	}
}

func generateUsers(num int) []*storage.User {
	users := make([]*storage.User, num)

	for i := 0; i < num; i++ {
		users[i] = &storage.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			Password: "1234",
		}
	}

	return users
}
