package repository

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/Ishkhan88/go-study/internal/model"
)

var (
	users         []model.User
	concerts      []model.Concert
	bookings      []model.Booking
	notifications []model.Notification

	muUsers         sync.Mutex
	muConcerts      sync.Mutex
	muBookings      sync.Mutex
	muNotifications sync.Mutex

	dataDirs = "data"
)

func ensureDataDir() error { return os.MkdirAll(dataDir, 0o755) }

// ===================== LOAD ALL ON START =====================

func LoadFromFiles() error {
	if err := ensureDataDir(); err != nil {
		return err
	}
	if err := loadUsers(); err != nil {
		return fmt.Errorf("load users: %w", err)
	}
	if err := loadConcerts(); err != nil {
		return fmt.Errorf("load concerts: %w", err)
	}
	if err := loadBookings(); err != nil {
		return fmt.Errorf("load bookings: %w", err)
	}
	if err := loadNotifications(); err != nil {
		return fmt.Errorf("load notifications: %w", err)
	}
	return nil
}

// ===================== USERS =====================

func usersPath() string { return filepath.Join(dataDir, "users.csv") }

func loadUsers() error {
	rows, err := readCSV(usersPath())
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}

	var loaded []model.User
	for i, r := range rows {
		if i == 0 && len(r) > 0 && r[0] == "id" { // header
			continue
		}
		if len(r) < 8 {
			continue
		}
		id, _ := strconv.Atoi(r[0])
		createdAt, _ := time.Parse(time.RFC3339, r[6])
		updatedAt, _ := time.Parse(time.RFC3339, r[7])

		loaded = append(loaded, model.User{
			ID:        id,
			FirstName: r[1],
			LastName:  r[2],
			Email:     r[3],
			Phone:     r[4],
			AvatarURL: r[5],
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}

	muUser.Lock()
	user = loaded
	muUser.Unlock()
	return nil
}

func saveUsers() error {
	muUser.Lock()
	defer muUser.Unlock()

	header := []string{"id", "first_name", "last_name", "email", "phone", "avatar_url", "created_at", "updated_at"}
	var rows [][]string
	rows = append(rows, header)

	for _, u := range user {
		rows = append(rows, []string{
			strconv.Itoa(u.ID),
			u.FirstName,
			u.LastName,
			u.Email,
			u.Phone,
			u.AvatarURL,
			u.CreatedAt.Format(time.RFC3339),
			u.UpdatedAt.Format(time.RFC3339),
		})
	}
	return writeCSV(usersPath(), rows)
}

func GetUserSafeCopy() []model.User {
	muUser.Lock()
	defer muUser.Unlock()
	cp := make([]model.User, len(user))
	copy(cp, user)
	return cp
}

func GetUserByID(id int) (model.User, bool) {
	muUser.Lock()
	defer muUser.Unlock()
	for _, u := range user {
		if u.ID == id {
			return u, true
		}
	}
	return model.User{}, false
}

func AddUser(u model.User) error {
	muUser.Lock()
	for _, ex := range user {
		if ex.ID == u.ID {
			muUser.Unlock()
			return fmt.Errorf("user id already exists")
		}
	}
	user = append(user, u)
	muUser.Unlock()
	return saveUsers()
}

func GetNextUserID() int {
	muUser.Lock()
	defer muUser.Unlock()

	maxID := 0
	for _, u := range user {
		if u.ID > maxID {
			maxID = u.ID
		}
	}
	return maxID + 1
}

func UpdateUser(id int, upd model.User) (model.User, bool, error) {
	muUser.Lock()
	for i := range user {
		if user[i].ID == id {
			upd.ID = id
			user[i] = upd
			muUser.Unlock()
			return upd, true, saveUsers()
		}
	}
	muUser.Unlock()
	return model.User{}, false, nil
}

func DeleteUser(id int) (bool, error) {
	muUser.Lock()
	for i := range user {
		if user[i].ID == id {
			user = append(user[:i], user[i+1:]...)
			muUser.Unlock()
			return true, saveUsers()
		}
	}
	muUser.Unlock()
	return false, nil
}

// ===================== CONCERTS =====================

func concertsPath() string { return filepath.Join(dataDir, "concerts.csv") }

func loadConcerts() error {
	rows, err := readCSV(concertsPath())
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}

	var loaded []model.Concert
	for i, r := range rows {
		if i == 0 && len(r) > 0 && r[0] == "id" {
			continue
		}
		if len(r) < 10 {
			continue
		}
		id, _ := strconv.Atoi(r[0])
		date, _ := time.Parse(time.RFC3339, r[2])
		price, _ := strconv.ParseFloat(r[4], 64)
		total, _ := strconv.Atoi(r[5])
		left, _ := strconv.Atoi(r[6])
		createdAt, _ := time.Parse(time.RFC3339, r[8])
		updatedAt, _ := time.Parse(time.RFC3339, r[9])

		loaded = append(loaded, model.Concert{
			ID:             id,
			Title:          r[1],
			Date:           date,
			Location:       r[3],
			TicketPrice:    price,
			TicketsTotal:   total,
			TicketsLeft:    left,
			OrganizerEmail: r[7],
			CreatedAt:      createdAt,
			UpdatedAt:      updatedAt,
		})
	}

	muConcert.Lock()
	concert = loaded
	muConcert.Unlock()
	return nil
}

func saveConcerts() error {
	muConcert.Lock()
	defer muConcert.Unlock()

	header := []string{"id", "title", "date", "location", "ticket_price", "tickets_total", "tickets_left", "organizer_email", "created_at", "updated_at"}
	var rows [][]string
	rows = append(rows, header)

	for _, c := range concert {
		rows = append(rows, []string{
			strconv.Itoa(c.ID),
			c.Title,
			c.Date.Format(time.RFC3339),
			c.Location,
			strconv.FormatFloat(c.TicketPrice, 'f', 2, 64),
			strconv.Itoa(c.TicketsTotal),
			strconv.Itoa(c.TicketsLeft),
			c.OrganizerEmail,
			c.CreatedAt.Format(time.RFC3339),
			c.UpdatedAt.Format(time.RFC3339),
		})
	}
	return writeCSV(concertsPath(), rows)
}

func GetConcertSafeCopy() []model.Concert {
	muConcert.Lock()
	defer muConcert.Unlock()
	cp := make([]model.Concert, len(concert))
	copy(cp, concert)
	return cp
}

func GetConcertByID(id int) (model.Concert, bool) {
	muConcert.Lock()
	defer muConcert.Unlock()
	for _, c := range concert {
		if c.ID == id {
			return c, true
		}
	}
	return model.Concert{}, false
}

func AddConcert(c model.Concert) error {
	muConcert.Lock()
	for _, ex := range concert {
		if ex.ID == c.ID {
			muConcert.Unlock()
			return fmt.Errorf("concert id already exists")
		}
	}
	concert = append(concert, c)
	muConcert.Unlock()
	return saveConcerts()
}

func GetNextConcertID() int {
	muConcert.Lock()
	defer muConcert.Unlock()

	maxID := 0
	for _, c := range concert {
		if c.ID > maxID {
			maxID = c.ID
		}
	}
	return maxID + 1
}

func UpdateConcert(id int, upd model.Concert) (model.Concert, bool, error) {
	muConcert.Lock()
	for i := range concert {
		if concert[i].ID == id {
			upd.ID = id
			concert[i] = upd
			muConcert.Unlock()
			return upd, true, saveConcerts()
		}
	}
	muConcert.Unlock()
	return model.Concert{}, false, nil
}

func DeleteConcert(id int) (bool, error) {
	muConcert.Lock()
	for i := range concert {
		if concert[i].ID == id {
			concert = append(concert[:i], concert[i+1:]...)
			muConcert.Unlock()
			return true, saveConcerts()
		}
	}
	muConcert.Unlock()
	return false, nil
}

// ===================== BOOKINGS =====================

func bookingsPath() string { return filepath.Join(dataDir, "bookings.csv") }

func loadBookings() error {
	rows, err := readCSV(bookingsPath())
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}

	var loaded []model.Booking
	for i, r := range rows {
		if i == 0 && len(r) > 0 && r[0] == "id" {
			continue
		}
		if len(r) < 6 {
			continue
		}
		id, _ := strconv.Atoi(r[0])
		uid, _ := strconv.Atoi(r[1])
		cid, _ := strconv.Atoi(r[2])
		createdAt, _ := time.Parse(time.RFC3339, r[4])
		updatedAt, _ := time.Parse(time.RFC3339, r[5])

		loaded = append(loaded, model.Booking{
			ID:        id,
			UserID:    uid,
			ConcertID: cid,
			Status:    r[3],
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}

	muBooking.Lock()
	booking = loaded
	muBooking.Unlock()
	return nil
}

func saveBookings() error {
	muBooking.Lock()
	defer muBooking.Unlock()

	header := []string{"id", "user_id", "concert_id", "status", "created_at", "updated_at"}
	var rows [][]string
	rows = append(rows, header)

	for _, b := range booking {
		rows = append(rows, []string{
			strconv.Itoa(b.ID),
			strconv.Itoa(b.UserID),
			strconv.Itoa(b.ConcertID),
			b.Status,
			b.CreatedAt.Format(time.RFC3339),
			b.UpdatedAt.Format(time.RFC3339),
		})
	}
	return writeCSV(bookingsPath(), rows)
}

func GetBookingSafeCopy() []model.Booking {
	muBooking.Lock()
	defer muBooking.Unlock()
	cp := make([]model.Booking, len(booking))
	copy(cp, booking)
	return cp
}

func GetBookingByID(id int) (model.Booking, bool) {
	muBooking.Lock()
	defer muBooking.Unlock()
	for _, b := range booking {
		if b.ID == id {
			return b, true
		}
	}
	return model.Booking{}, false
}

func AddBooking(b model.Booking) error {
	muBooking.Lock()
	for _, ex := range booking {
		if ex.ID == b.ID {
			muBooking.Unlock()
			return fmt.Errorf("booking id already exists")
		}
	}
	booking = append(booking, b)
	muBooking.Unlock()
	return saveBookings()
}

func GetNextBookingID() int {
	muBooking.Lock()
	defer muBooking.Unlock()

	maxID := 0
	for _, b := range booking {
		if b.ID > maxID {
			maxID = b.ID
		}
	}
	return maxID + 1
}

func UpdateBooking(id int, upd model.Booking) (model.Booking, bool, error) {
	muBooking.Lock()
	for i := range booking {
		if booking[i].ID == id {
			upd.ID = id
			booking[i] = upd
			muBooking.Unlock()
			return upd, true, saveBookings()
		}
	}
	muBooking.Unlock()
	return model.Booking{}, false, nil
}

func DeleteBooking(id int) (bool, error) {
	muBooking.Lock()
	for i := range booking {
		if booking[i].ID == id {
			booking = append(booking[:i], booking[i+1:]...)
			muBooking.Unlock()
			return true, saveBookings()
		}
	}
	muBooking.Unlock()
	return false, nil
}

// ===================== NOTIFICATIONS =====================

func notificationsPath() string { return filepath.Join(dataDir, "notifications.csv") }

func loadNotifications() error {
	rows, err := readCSV(notificationsPath())
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}

	var loaded []model.Notification
	for i, r := range rows {
		if i == 0 && len(r) > 0 && r[0] == "id" {
			continue
		}
		if len(r) < 6 {
			continue
		}
		id, _ := strconv.Atoi(r[0])
		uid, _ := strconv.Atoi(r[1])
		cid, _ := strconv.Atoi(r[2])
		sentAt, _ := time.Parse(time.RFC3339, r[5])

		loaded = append(loaded, model.Notification{
			ID:        id,
			UserID:    uid,
			ConcertID: cid,
			Status:    r[3],
			SentAt:    sentAt,
		})
	}

	muNotification.Lock()
	notification = loaded
	muNotification.Unlock()
	return nil
}

func saveNotifications() error {
	muNotification.Lock()
	defer muNotification.Unlock()

	header := []string{"id", "user_id", "concert_id", "status", "sent_at", "sent_at_rfc3339"}
	var rows [][]string
	rows = append(rows, header)

	for _, n := range notification {
		rows = append(rows, []string{
			strconv.Itoa(n.ID),
			strconv.Itoa(n.UserID),
			strconv.Itoa(n.ConcertID),
			n.Status,
			"", // резерв
			n.SentAt.Format(time.RFC3339),
		})
	}
	return writeCSV(notificationsPath(), rows)
}

func GetNotificationSafeCopy() []model.Notification {
	muNotification.Lock()
	defer muNotification.Unlock()
	cp := make([]model.Notification, len(notification))
	copy(cp, notification)
	return cp
}

func GetNotificationByID(id int) (model.Notification, bool) {
	muNotification.Lock()
	defer muNotification.Unlock()
	for _, n := range notification {
		if n.ID == id {
			return n, true
		}
	}
	return model.Notification{}, false
}

func AddNotification(n model.Notification) error {
	muNotification.Lock()
	for _, ex := range notification {
		if ex.ID == n.ID {
			muNotification.Unlock()
			return fmt.Errorf("notification id already exists")
		}
	}
	notification = append(notification, n)
	muNotification.Unlock()
	return saveNotifications()
}

func GetNextNotificationID() int {
	muNotification.Lock()
	defer muNotification.Unlock()

	maxID := 0
	for _, n := range notification {
		if n.ID > maxID {
			maxID = n.ID
		}
	}
	return maxID + 1
}

func UpdateNotification(id int, upd model.Notification) (model.Notification, bool, error) {
	muNotification.Lock()
	for i := range notification {
		if notification[i].ID == id {
			upd.ID = id
			notification[i] = upd
			muNotification.Unlock()
			return upd, true, saveNotifications()
		}
	}
	muNotification.Unlock()
	return model.Notification{}, false, nil
}

func DeleteNotification(id int) (bool, error) {
	muNotification.Lock()
	for i := range notification {
		if notification[i].ID == id {
			notification = append(notification[:i], notification[i+1:]...)
			muNotification.Unlock()
			return true, saveNotifications()
		}
	}
	muNotification.Unlock()
	return false, nil
}

// ===================== CSV HELPERS =====================

func readCSV(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	return r.ReadAll()
}

func writeCSV(path string, rows [][]string) error {
	if err := ensureDataDir(); err != nil {
		return err
	}
	f, err := os.Create(path) // перезапись целиком
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	if err := w.WriteAll(rows); err != nil {
		return err
	}
	w.Flush()
	return w.Error()
}
