package repositories

import (
	"go-ticketing/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EventRepository interface {
	Create(event *models.Event) error
	FindByID(id uint) (*models.Event, error)
	//ambil + kunci baris dalam transaksi untuk mencegah race condition
	FindByIDForUpdate(tx *gorm.DB, id uint) (*models.Event, error)
	Update(tx *gorm.DB, event *models.Event) error
	ListAll() ([]models.Event, error)
}

// simpan koneksi database, tiap kali eventRepositoryImpl dipanggil, dia sudah punya akses ke DB lewat field ini.
type eventRepository struct {
	db *gorm.DB
}

// constructor yg ngasih koneksi DB ke fungsi ini. returnnya return &eventRepository{db}) struct,
// tapi Go menganggapnya sebagai interface karena struct itu cocok dengan kontrak interface.
func NewEventRepository(db *gorm.DB) EventRepository {
	//Return EventRepository (interface) biar caller tidak tergantung detail struct, hanya tahu kontraknya saja.
	//Kita bikin instance dari struct eventRepository. Lalu kita simpan koneksi database db di dalamnya (jadi tiap method bisa akses DB).
	//Karena pakai & → yang dikembalikan adalah pointer ke struct itu.
	//repo aslinya adalah *eventRepository (struct). Tapi karena *eventRepository punya semua method yang diminta EventRepository, Go langsung bisa menganggapnya sebagai EventRepository (interface).
	return &eventRepository{db}
}

func (r *eventRepository) Create(event *models.Event) error {
	return r.db.Create(event).Error
}

func (r *eventRepository) FindByID(id uint) (*models.Event, error) {
	var e models.Event
	if err := r.db.First(&e, id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

// FindByIDForUpdate locks the selected row using SELECT ... FOR UPDATE
func (r *eventRepository) FindByIDForUpdate(tx *gorm.DB, id uint) (*models.Event, error) {
	var e models.Event
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&e, id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *eventRepository) Update(tx *gorm.DB, event *models.Event) error {
	// if tx provided, use it; otherwise use base db
	if tx != nil {
		return tx.Save(event).Error
	}
	return r.db.Save(event).Error
}

func (r *eventRepository) ListAll() ([]models.Event, error) {
	var list []models.Event
	if err := r.db.Order("start_at asc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
