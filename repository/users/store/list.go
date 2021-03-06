package store

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	domain "github.com/kanok-p/go-clean-architecture/domain/users"
)

func (s *Store) List(ctx context.Context, offset, limit int64, filter bson.M) (int64, []*domain.Users, error) {
	total, err := s.collection().CountDocuments(ctx, filter)
	if err != nil {
		return total, nil, err
	}

	cursor, err := s.collection().Find(ctx, filter, options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.M{"createdAt": -1}))
	if err != nil {
		return total, nil, err
	}
	defer func() { _ = cursor.Close(ctx) }()

	list := make([]*domain.Users, 0)
	for cursor.Next(ctx) {
		user := &domain.Users{}
		if err := cursor.Decode(user); err != nil {
			return total, nil, err
		}

		list = append(list, user)
	}

	return total, list, nil
}
