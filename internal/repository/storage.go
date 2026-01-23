package repository

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/Ishkhan88/go-study/internal/model"
)

var (
	user         []model.User
	concert      []model.Concert
	booking      []model.Booking
	notification []model.Notification

	muUser         sync.Mutex
	muConcert      sync.Mutex
	muBooking      sync.Mutex
	muNotification sync.Mutex

	// папка и файлы хранения
	dataDir           = "data"
	usersFile         = filepath.Join(dataDir, "users.jsonl")
	concertsFile      = filepath.Join(dataDir, "concerts.jsonl")
	bookingsFile      = filepath.Join(dataDir, "bookings.jsonl")
	notificationsFile = filepath.Join(dataDir, "notifications.jsonl")
)

// ---------- PUBLIC API ----------

// LoadFromFiles нужно вызвать при старте программы.
// Она наполнит все слайсы данными из файлов.
func LoadFromFile() error {
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return fmt.Errorf("mkdir data dir: %w", err)
	}

	loadedUsers, err := readJSONLines[model.User](usersFile)
	if err != nil {
		return fmt.Errorf("read users: %w", err)
	}
	muUser.Lock()
	user = loadedUsers
	muUser.Unlock()

	loadedConcerts, err := readJSONLines[model.Concert](concertsFile)
	if err != nil {
		return fmt.Errorf("read concerts: %w", err)
	}
	muConcert.Lock()
	concert = loadedConcerts
	muConcert.Unlock()

	loadedBookings, err := readJSONLines[model.Booking](bookingsFile)
	if err != nil {
		return fmt.Errorf("read bookings: %w", err)
	}
	muBooking.Lock()
	booking = loadedBookings
	muBooking.Unlock()

	loadedNotifications, err := readJSONLines[model.Notification](notificationsFile)
	if err != nil {
		return fmt.Errorf("read notifications: %w", err)
	}
	muNotification.Lock()
	notification = loadedNotifications
	muNotification.Unlock()

	log.Println("Loaded data from files:",
		"users=", len(loadedUsers),
		"concerts=", len(loadedConcerts),
		"bookings=", len(loadedBookings),
		"notifications=", len(loadedNotifications),
	)

	return nil
}

// SaveEntity распределяет сущности по типам.
// Теперь это “сохранить в память + сохранить в файл”.
func SaveEntity(e model.Entity) error {
	switch v := e.(type) {
	case model.User:
		return addUser(v)
	case model.Concert:
		return addConcert(v)
	case model.Booking:
		return addBooking(v)
	case model.Notification:
		return addNotification(v)
	default:
		return fmt.Errorf("unknown entity type: %T", v)
	}
}

// безопасные копии
func GetUsersSafeCopy() []model.User {
	muUser.Lock()
	defer muUser.Unlock()
	copySlice := make([]model.User, len(user))
	copy(copySlice, user)
	return copySlice
}

func GetConcertsSafeCopy() []model.Concert {
	muConcert.Lock()
	defer muConcert.Unlock()
	copySlice := make([]model.Concert, len(concert))
	copy(copySlice, concert)
	return copySlice
}

func GetBookingsSafeCopy() []model.Booking {
	muBooking.Lock()
	defer muBooking.Unlock()
	copySlice := make([]model.Booking, len(booking))
	copy(copySlice, booking)
	return copySlice
}

func GetNotificationsSafeCopy() []model.Notification {
	muNotification.Lock()
	defer muNotification.Unlock()
	copySlice := make([]model.Notification, len(notification))
	copy(copySlice, notification)
	return copySlice
}

// ---------- INTERNAL: add + persist ----------

func addUser(u model.User) error {
	muUser.Lock()
	user = append(user, u)
	muUser.Unlock()
	return appendJSONLine(usersFile, u)
}

func addConcert(c model.Concert) error {
	muConcert.Lock()
	concert = append(concert, c)
	muConcert.Unlock()
	return appendJSONLine(concertsFile, c)
}

func addBooking(b model.Booking) error {
	muBooking.Lock()
	booking = append(booking, b)
	muBooking.Unlock()
	return appendJSONLine(bookingsFile, b)
}

func addNotification(n model.Notification) error {
	muNotification.Lock()
	notification = append(notification, n)
	muNotification.Unlock()
	return appendJSONLine(notificationsFile, n)
}

// ---------- helpers: jsonl read/write ----------

func appendJSONLine[T any](filePath string, v T) error {
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	// каждая сущность = отдельная строка JSON
	if _, err := f.Write(append(b, '\n')); err != nil {
		return err
	}
	return nil
}

func readJSONLines[T any](filePath string) ([]T, error) {
	// если файла нет — это не ошибка, просто вернём пустой слайс
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []T{}, nil
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var out []T
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) == 0 {
			continue
		}
		var item T
		if err := json.Unmarshal(line, &item); err != nil {
			return nil, fmt.Errorf("bad json line in %s: %w", filePath, err)
		}
		out = append(out, item)
	}

	if err := sc.Err(); err != nil {
		return nil, err
	}

	return out, nil
}
