package mongo

import (
	"context"
	"fmt"
	"github.com/zhashkevych/go-clean-architecture/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Bookmark struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID primitive.ObjectID `bson:"userId"`
	URL    string             `bson:"url"`
	Title  string             `bson:"title"`
	Tags   []string           `bson:"tags"`
}

type BookmarkRepository struct {
	db *mongo.Collection
}

func NewBookmarkRepository(db *mongo.Database, collection string) *BookmarkRepository {
	return &BookmarkRepository{
		db: db.Collection(collection),
	}
}

func (r BookmarkRepository) CreateBookmark(ctx context.Context, user *models.User, bm *models.Bookmark) error {
	bm.UserID = user.ID

	model := toModel(bm)

	res, err := r.db.InsertOne(ctx, model)
	if err != nil {
		return err
	}

	bm.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r BookmarkRepository) GetBookmarks(ctx context.Context, user *models.User) ([]*models.Bookmark, error) {
	uid, _ := primitive.ObjectIDFromHex(user.ID)
	cur, err := r.db.Find(ctx, bson.M{
		"userId": uid,
	})
	defer cur.Close(ctx)

	if err != nil {
		return nil, err
	}

	out := make([]*Bookmark, 0)

	for cur.Next(ctx) {
		user := new(Bookmark)
		err := cur.Decode(user)
		if err != nil {
			return nil, err
		}

		out = append(out, user)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return toBookmarks(out), nil
}

func (r BookmarkRepository) GetBookmarksByTags(ctx context.Context, user *models.User, tags []string) ([]*models.Bookmark, error) {
	uid, _ := primitive.ObjectIDFromHex(user.ID)

	filter := bson.M{
		"userId": uid,
		"tags":   bson.M{"$all": tags},
	}

	cur, err := r.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	out := make([]*Bookmark, 0)

	for cur.Next(ctx) {
		bookmark := new(Bookmark)
		err := cur.Decode(bookmark)
		if err != nil {
			return nil, err
		}

		out = append(out, bookmark)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return toBookmarks(out), nil
}

func (r BookmarkRepository) DeleteBookmark(ctx context.Context, user *models.User, id string) error {
	objID, _ := primitive.ObjectIDFromHex(id)
	uID, _ := primitive.ObjectIDFromHex(user.ID)

	_, err := r.db.DeleteOne(ctx, bson.M{"_id": objID, "userId": uID})
	return err
}

func (r BookmarkRepository) UpdateBookmarkTags(ctx context.Context, user *models.User, id string, tags []string) error {
	objID, _ := primitive.ObjectIDFromHex(id)
	uID, _ := primitive.ObjectIDFromHex(user.ID)

	filter := bson.M{
		"_id":    objID,
		"userId": uID,
	}

	update := bson.M{
		"$set": bson.M{
			"tags": tags,
		},
	}

	_, err := r.db.UpdateOne(ctx, filter, update)
	return err
}

func (r BookmarkRepository) MergeTags(ctx context.Context, user *models.User, fromTag string, toTag string) error {
	uid, _ := primitive.ObjectIDFromHex(user.ID)

	filter := bson.M{
		"userId": uid,
		"tags":   fromTag,
	}

	update := bson.M{
		"$pull": bson.M{
			"tags": fromTag,
		},
		"$addToSet": bson.M{
			"tags": toTag,
		},
	}

	_, err := r.db.UpdateMany(ctx, filter, update)
	return err
}

func (r BookmarkRepository) BatchAddTags(ctx context.Context, user *models.User, bookmarkIDs []string, tags []string) error {
	uid, _ := primitive.ObjectIDFromHex(user.ID)

	var objIDs []primitive.ObjectID
	for _, id := range bookmarkIDs {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return fmt.Errorf("invalid bookmark ID: %s", id)
		}
		objIDs = append(objIDs, objID)
	}

	filter := bson.M{
		"_id":    bson.M{"$in": objIDs},
		"userId": uid,
	}

	update := bson.M{
		"$addToSet": bson.M{
			"tags": bson.M{"$each": tags},
		},
	}

	_, err := r.db.UpdateMany(ctx, filter, update)
	return err
}

func (r BookmarkRepository) BatchRemoveTags(ctx context.Context, user *models.User, bookmarkIDs []string, tags []string) error {
	uid, _ := primitive.ObjectIDFromHex(user.ID)

	var objIDs []primitive.ObjectID
	for _, id := range bookmarkIDs {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return fmt.Errorf("invalid bookmark ID: %s", id)
		}
		objIDs = append(objIDs, objID)
	}

	filter := bson.M{
		"_id":    bson.M{"$in": objIDs},
		"userId": uid,
	}

	update := bson.M{
		"$pull": bson.M{
			"tags": bson.M{"$in": tags},
		},
	}

	_, err := r.db.UpdateMany(ctx, filter, update)
	return err
}

func (r BookmarkRepository) GetAllTags(ctx context.Context, user *models.User) ([]string, error) {
	uid, _ := primitive.ObjectIDFromHex(user.ID)

	matchStage := bson.M{
		"$match": bson.M{
			"userId": uid,
		},
	}

	unwindStage := bson.M{
		"$unwind": "$tags",
	}

	groupStage := bson.M{
		"$group": bson.M{
			"_id": nil,
			"tags": bson.M{
				"$addToSet": "$tags",
			},
		},
	}

	cur, err := r.db.Aggregate(ctx, mongo.Pipeline{matchStage, unwindStage, groupStage})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var result bson.M
	if cur.Next(ctx) {
		if err := cur.Decode(&result); err != nil {
			return nil, err
		}

		tags, ok := result["tags"].(primitive.A)
		if !ok {
			return []string{}, nil
		}

		out := make([]string, 0, len(tags))
		for _, tag := range tags {
			if t, ok := tag.(string); ok {
				out = append(out, t)
			}
		}

		return out, nil
	}

	return []string{}, nil
}

func toModel(b *models.Bookmark) *Bookmark {
	uid, _ := primitive.ObjectIDFromHex(b.UserID)

	return &Bookmark{
		UserID: uid,
		URL:    b.URL,
		Title:  b.Title,
		Tags:   b.Tags,
	}
}

func toBookmark(b *Bookmark) *models.Bookmark {
	return &models.Bookmark{
		ID:     b.ID.Hex(),
		UserID: b.UserID.Hex(),
		URL:    b.URL,
		Title:  b.Title,
		Tags:   b.Tags,
	}
}

func toBookmarks(bs []*Bookmark) []*models.Bookmark {
	out := make([]*models.Bookmark, len(bs))

	for i, b := range bs {
		out[i] = toBookmark(b)
	}

	return out
}
