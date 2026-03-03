package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Ishkhan88/go-study/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
	mongoDB     *mongo.Database

	colUsers         *mongo.Collection
	colConcerts      *mongo.Collection
	colBookings      *mongo.Collection
	colNotifications *mongo.Collection
)

const mongoOpTimeout = 5 * time.Second

func mustMongo() {
	if mongoClient == nil ||
		mongoDB == nil ||
		colUsers == nil ||
		colConcerts == nil ||
		colBookings == nil ||
		colNotifications == nil {
		panic("mongo is not initialized: call repository.InitMongo() on startup")
	}
}

// InitMongo нужно вызвать при старте сервиса.
func InitMongo(ctx context.Context, uri, dbName string) error {
	if uri == "" {
		return fmt.Errorf("mongo uri is empty")
	}
	if dbName == "" {
		return fmt.Errorf("mongo db name is empty")
	}

	cctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(cctx, options.Client().ApplyURI(uri))
	if err != nil {
		return fmt.Errorf("mongo connect: %w", err)
	}

	// ping
	pctx, pcancel := context.WithTimeout(ctx, 5*time.Second)
	defer pcancel()
	if err := client.Ping(pctx, nil); err != nil {
		_ = client.Disconnect(context.Background())
		return fmt.Errorf("mongo ping: %w", err)
	}

	mongoClient = client
	mongoDB = client.Database(dbName)

	colUsers = mongoDB.Collection("users")
	colConcerts = mongoDB.Collection("concerts")
	colBookings = mongoDB.Collection("bookings")
	colNotifications = mongoDB.Collection("notifications")

	// уникальный индекс по id в каждой коллекции
	if err := ensureUniqueIDIndex(ctx, colUsers); err != nil {
		return err
	}
	if err := ensureUniqueIDIndex(ctx, colConcerts); err != nil {
		return err
	}
	if err := ensureUniqueIDIndex(ctx, colBookings); err != nil {
		return err
	}
	if err := ensureUniqueIDIndex(ctx, colNotifications); err != nil {
		return err
	}

	return nil
}

func CloseMongo(ctx context.Context) error {
	if mongoClient == nil {
		return nil
	}
	cctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return mongoClient.Disconnect(cctx)
}

func ensureUniqueIDIndex(ctx context.Context, col *mongo.Collection) error {
	ictx, cancel := context.WithTimeout(ctx, mongoOpTimeout)
	defer cancel()

	_, err := col.Indexes().CreateOne(ictx, mongo.IndexModel{
		Keys:    bson.D{{Key: "id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("create unique index on %s.id: %w", col.Name(), err)
	}
	return nil
}

// ===================== DOCS (bson tags) =====================

type userDoc struct {
	ID        int       `bson:"id"`
	FirstName string    `bson:"first_name"`
	LastName  string    `bson:"last_name"`
	Email     string    `bson:"email"`
	Phone     string    `bson:"phone"`
	AvatarURL string    `bson:"avatar_url"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func toUserDoc(u model.User) userDoc {
	return userDoc{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Phone:     u.Phone,
		AvatarURL: u.AvatarURL,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
func fromUserDoc(d userDoc) model.User {
	return model.User{
		ID:        d.ID,
		FirstName: d.FirstName,
		LastName:  d.LastName,
		Email:     d.Email,
		Phone:     d.Phone,
		AvatarURL: d.AvatarURL,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

type concertDoc struct {
	ID             int       `bson:"id"`
	Title          string    `bson:"title"`
	Date           time.Time `bson:"date"`
	Location       string    `bson:"location"`
	TicketPrice    float64   `bson:"ticket_price"`
	TicketsTotal   int       `bson:"tickets_total"`
	TicketsLeft    int       `bson:"tickets_left"`
	OrganizerEmail string    `bson:"organizer_email"`
	CreatedAt      time.Time `bson:"created_at"`
	UpdatedAt      time.Time `bson:"updated_at"`
}

func toConcertDoc(c model.Concert) concertDoc {
	return concertDoc{
		ID:             c.ID,
		Title:          c.Title,
		Date:           c.Date,
		Location:       c.Location,
		TicketPrice:    c.TicketPrice,
		TicketsTotal:   c.TicketsTotal,
		TicketsLeft:    c.TicketsLeft,
		OrganizerEmail: c.OrganizerEmail,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}
}
func fromConcertDoc(d concertDoc) model.Concert {
	return model.Concert{
		ID:             d.ID,
		Title:          d.Title,
		Date:           d.Date,
		Location:       d.Location,
		TicketPrice:    d.TicketPrice,
		TicketsTotal:   d.TicketsTotal,
		TicketsLeft:    d.TicketsLeft,
		OrganizerEmail: d.OrganizerEmail,
		CreatedAt:      d.CreatedAt,
		UpdatedAt:      d.UpdatedAt,
	}
}

type bookingDoc struct {
	ID        int       `bson:"id"`
	UserID    int       `bson:"user_id"`
	ConcertID int       `bson:"concert_id"`
	Status    string    `bson:"status"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func toBookingDoc(b model.Booking) bookingDoc {
	return bookingDoc{
		ID:        b.ID,
		UserID:    b.UserID,
		ConcertID: b.ConcertID,
		Status:    b.Status,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}
func fromBookingDoc(d bookingDoc) model.Booking {
	return model.Booking{
		ID:        d.ID,
		UserID:    d.UserID,
		ConcertID: d.ConcertID,
		Status:    d.Status,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

type notificationDoc struct {
	ID        int       `bson:"id"`
	ConcertID int       `bson:"concert_id"`
	UserID    int       `bson:"user_id"`
	Status    string    `bson:"status"`
	SentAt    time.Time `bson:"sent_at"`
}

func toNotificationDoc(n model.Notification) notificationDoc {
	return notificationDoc{
		ID:        n.ID,
		ConcertID: n.ConcertID,
		UserID:    n.UserID,
		Status:    n.Status,
		SentAt:    n.SentAt,
	}
}
func fromNotificationDoc(d notificationDoc) model.Notification {
	return model.Notification{
		ID:        d.ID,
		ConcertID: d.ConcertID,
		UserID:    d.UserID,
		Status:    d.Status,
		SentAt:    d.SentAt,
	}
}

// ===================== helpers =====================

func nextID(ctx context.Context, col *mongo.Collection) (int, error) {
	octx, cancel := context.WithTimeout(ctx, mongoOpTimeout)
	defer cancel()

	var r struct {
		ID int `bson:"id"`
	}

	err := col.FindOne(
		octx,
		bson.D{},
		options.FindOne().
			SetSort(bson.D{{Key: "id", Value: -1}}).
			SetProjection(bson.D{{Key: "id", Value: 1}}),
	).Decode(&r)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return 1, nil
	}
	if err != nil {
		return 0, err
	}
	return r.ID + 1, nil
}

// ===================== USERS =====================

func GetNextUserID() int {
	mustMongo()
	id, _ := nextID(context.Background(), colUsers)
	return id
}

func AddUser(u model.User) error {
	mustMongo()

	ctx, cancel := context.WithTimeout(context.Background(), mongoOpTimeout)
	defer cancel()

	_, err := colUsers.InsertOne(ctx, toUserDoc(u))
	if mongo.IsDuplicateKeyError(err) {
		return fmt.Errorf("user id already exists")
	}
	if err == nil {
		AuditEvent("user", u.ID, "create", u)
	}
	return err
}

func GetUserByID(id int) (model.User, bool) {
	mustMongo()

	ctx, cancel := context.WithTimeout(context.Background(), mongoOpTimeout)
	defer cancel()

	var d userDoc
	err := colUsers.FindOne(ctx, bson.M{"id": id}).Decode(&d)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return model.User{}, false
	}
	if err != nil {
		return model.User{}, false
	}
	return fromUserDoc(d), true
}

func UpdateUser(id int, upd model.User) (model.User, bool, error) {
	mustMongo()

	ctx, cancel := context.WithTimeout(context.Background(), mongoOpTimeout)
	defer cancel()

	upd.ID = id
	res, err := colUsers.ReplaceOne(ctx, bson.M{"id": id}, toUserDoc(upd))
	if err != nil {
		return model.User{}, false, err
	}
	if res.MatchedCount == 0 {
		return model.User{}, false, nil
	}

	AuditEvent("user", id, "update", upd)
	return upd, true, nil
}

func DeleteUser(id int) (bool, error) {
	mustMongo()

	ctx, cancel := context.WithTimeout(context.Background(), mongoOpTimeout)
	defer cancel()

	res, err := colUsers.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return false, err
	}
	deleted := res.DeletedCount > 0
	if deleted {
		AuditEvent("user", id, "delete", map[string]any{"id": id})
	}
	return deleted, nil
}

func GetUserSafeCopy() []model.User {
	mustMongo()

	ctx, cancel := context.WithTimeout(context.Background(), mongoOpTimeout)
	defer cancel()

	cur, err := colUsers.Find(ctx, bson.D{})
	if err != nil {
		return []model.User{}
	}
	defer cur.Close(ctx)

	var out []model.User
	for cur.Next(ctx) {
		var d userDoc
		if err := cur.Decode(&d); err == nil {
			out = append(out, fromUserDoc(d))
		}
	}
	return out
}

// ===================== CONCERTS =====================

func GetNextConcertID() int {
	mustMongo()
	id, _ := nextID(context.Background(), colConcerts)
	return id
}

func AddConcert(c model.Concert) error {
	mustMongo()

	ctx, cancel := context.WithTimeout(context.Background(), mongoOpTimeout)
	defer cancel()

	_, err := colConcerts.InsertOne(ctx, toConcertDoc(c))
	if mongo.IsDuplicateKeyError(err) {
		return fmt.Errorf("concert id already exists")
	}
	if err == nil {
		AuditEvent("concert", c.ID, "create", c)
	}
	return err
}

func GetConcertByID(id int) (model.Concert, bool) {
	mustMongo()

	ctx, cancel := context.WithTimeout(context.Background(), mongoOpTimeout)
	defer cancel()

	var d concertDoc
	err := colConcerts.FindOne(ctx, bson.M{"id": id}).Decode(&d)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return model.Concert{}, false
	}
	if err != nil {
		return model.Concert{}, false
	}
	return fromConcertDoc(d), true
}

func UpdateConcert(id int, upd model.Concert) (model.Concert, bool, error) {
	mustMongo()

	ctx, cancel := context.WithTimeout(context.Background(), mongoOpTimeout)
	defer cancel()

	upd.ID = id
	res, err := colConcerts.ReplaceOne(ctx, bson.M{"id": id}, toConcertDoc(upd))
	if err != nil {
		return model.Concert{}, false, err
	}
	if res.MatchedCount == 0 {
		return model.Concert{}, false, nil
	}

	AuditEvent("concert", id, "update", upd)
	return upd, true, nil
}

func DeleteConcert(id int) (bool, error) {
	mustMongo()

	ctx, cancel := context.WithTimeout(context.Background(), mongoOpTimeout)
	defer cancel()

	res, err := colConcerts.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return false, err
	}
	deleted := res.DeletedCount > 0
	if deleted {
		AuditEvent("concert", id, "delete", map[string]any{"id": id})
	}
	return deleted, nil
}

func GetConcertSafeCopy() []model.Concert {
	mustMongo()

	ctx, cancel := context.WithTimeout(context.Background(), mongoOpTimeout)
	defer cancel()

	cur, err := colConcerts.Find(ctx, bson.D{})
	if err != nil {
		return []model.Concert{}
	}
	defer cur.Close(ctx)

	var out []model.Concert
	for cur.Next(ctx) {
		var d concertDoc
		if err := cur.Decode(&d); err == nil {
			out = append(out, fromConcertDoc(d))
		}
	}
	return out
}

// ===================== BOOKINGS =====================

func GetNextBookingID() int {
	mustMongo()
	id, _ := nextID(context.Background(), colBookings)
	return id
}

func AddBooking(b model.Booking) error {
	mustMongo()

	ctx, cancel := context.WithTimeout(context.Background(), mongoOpTimeout)
	defer cancel()

	_, err := colBookings.InsertOne(ctx, toBookingDoc(b))
	if mongo.IsDuplicateKeyError(err) {
		return fmt.Errorf("booking id already exists")
	}
	if err == nil {
		AuditEvent("booking", b.ID, "create", b)
	}
	return err
}

func GetBookingByID(id int) (model.Booking, bool) {
	mustMongo()

	ctx, cancel := context.WithTimeout(context.Background(), mongoOpTimeout)
	defer cancel()

	var d bookingDoc
	err := colBookings.FindOne(ctx, bson.M{"id": id}).Decode(&d)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return model.Booking{}, false
	}
	if err != nil {
		return model.Booking{}, false
	}
	return fromBookingDoc(d), true
}

func UpdateBooking(id int, upd model.Booking) (model.Booking, bool, error) {
	mustMongo()

	ctx, cancel := context.WithTimeout(context.Background(), mongoOpTimeout)
	defer cancel()

	upd.ID = id
	res, err := colBookings.ReplaceOne(ctx, bson.M{"id": id}, toBookingDoc(upd))
	if err != nil {
		return model.Booking{}, false, err
	}
	if res.MatchedCount == 0 {
		return model.Booking{}, false, nil
	}

	AuditEvent("booking", id, "update", upd)
	return upd, true, nil
}

func DeleteBooking(id int) (bool, error) {
	mustMongo()

	ctx, cancel := context.WithTimeout(context.Background(), mongoOpTimeout)
	defer cancel()

	res, err := colBookings.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return false, err
	}
	deleted := res.DeletedCount > 0
	if deleted {
		AuditEvent("booking", id, "delete", map[string]any{"id": id})
	}
	return deleted, nil
}

func GetBookingSafeCopy() []model.Booking {
	mustMongo()

	ctx, cancel := context.WithTimeout(context.Background(), mongoOpTimeout)
	defer cancel()

	cur, err := colBookings.Find(ctx, bson.D{})
	if err != nil {
		return []model.Booking{}
	}
	defer cur.Close(ctx)

	var out []model.Booking
	for cur.Next(ctx) {
		var d bookingDoc
		if err := cur.Decode(&d); err == nil {
			out = append(out, fromBookingDoc(d))
		}
	}
	return out
}

// ===================== NOTIFICATIONS =====================

func GetNextNotificationID() int {
	mustMongo()
	id, _ := nextID(context.Background(), colNotifications)
	return id
}

func AddNotification(n model.Notification) error {
	mustMongo()

	ctx, cancel := context.WithTimeout(context.Background(), mongoOpTimeout)
	defer cancel()

	_, err := colNotifications.InsertOne(ctx, toNotificationDoc(n))
	if mongo.IsDuplicateKeyError(err) {
		return fmt.Errorf("notification id already exists")
	}
	if err == nil {
		AuditEvent("notification", n.ID, "create", n)
	}
	return err
}

func GetNotificationByID(id int) (model.Notification, bool) {
	mustMongo()

	ctx, cancel := context.WithTimeout(context.Background(), mongoOpTimeout)
	defer cancel()

	var d notificationDoc
	err := colNotifications.FindOne(ctx, bson.M{"id": id}).Decode(&d)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return model.Notification{}, false
	}
	if err != nil {
		return model.Notification{}, false
	}
	return fromNotificationDoc(d), true
}

func UpdateNotification(id int, upd model.Notification) (model.Notification, bool, error) {
	mustMongo()

	ctx, cancel := context.WithTimeout(context.Background(), mongoOpTimeout)
	defer cancel()

	upd.ID = id
	res, err := colNotifications.ReplaceOne(ctx, bson.M{"id": id}, toNotificationDoc(upd))
	if err != nil {
		return model.Notification{}, false, err
	}
	if res.MatchedCount == 0 {
		return model.Notification{}, false, nil
	}

	AuditEvent("notification", id, "update", upd)
	return upd, true, nil
}

func DeleteNotification(id int) (bool, error) {
	mustMongo()

	ctx, cancel := context.WithTimeout(context.Background(), mongoOpTimeout)
	defer cancel()

	res, err := colNotifications.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return false, err
	}
	deleted := res.DeletedCount > 0
	if deleted {
		AuditEvent("notification", id, "delete", map[string]any{"id": id})
	}
	return deleted, nil
}

func GetNotificationSafeCopy() []model.Notification {
	mustMongo()

	ctx, cancel := context.WithTimeout(context.Background(), mongoOpTimeout)
	defer cancel()

	cur, err := colNotifications.Find(ctx, bson.D{})
	if err != nil {
		return []model.Notification{}
	}
	defer cur.Close(ctx)

	var out []model.Notification
	for cur.Next(ctx) {
		var d notificationDoc
		if err := cur.Decode(&d); err == nil {
			out = append(out, fromNotificationDoc(d))
		}
	}
	return out
}
