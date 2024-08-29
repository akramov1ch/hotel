package adjust

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"user-service/models"
	db "user-service/pkg/database/sql"
	"user-service/pkg/proto/notification"

	"golang.org/x/crypto/bcrypt"
)

type Database struct {
	Db *sql.DB
	N  notification.NotificationClient
}

func (u *Database) LogIn(ctx context.Context, req *models.LogInRequest) (*models.LogInResponse, error) {
	query, args, err := db.LogIn(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var password string
	if err := u.Db.QueryRow(query, args...).Scan(&password); err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Println(password)
	fmt.Println(req.Password)
	check := u.ComparePassword(password, req.Password)
	if check {
		_, err := u.N.Email(ctx, &notification.EmailSend{Email: req.Email, Message: "Siz qayta kirdingiz oldin ham kirgansiz 😉"})
		if err != nil {
			log.Println("xatyolik bor 😭", err)
		}
		return &models.LogInResponse{Status: true}, nil
	}
	return nil, errors.New("parol mos emas 🤨")
}

func (u *Database) CreateUser(ctx context.Context, req *models.RegisterUserRequest) (*models.GeneralResponse, error) {
	req.Password = u.Hashing(req.Password)
	query, args, err := db.Create(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var id int

	if err := u.Db.QueryRow(query, args...).Scan(&id); err != nil {
		log.Println(err)
		return nil, err
	}
	_, err = u.N.AddUser(ctx, &notification.AddnewUser{UserId: strconv.Itoa(id)})
	if err != nil {
		log.Println(err)
	}
	_, err = u.N.Notification(ctx, &notification.ProduceMessage{UserId: int32(id), Message: fmt.Sprintf("your account is created successfully with this is %v you can get notifications with it", id)})
	if err != nil {
		log.Println(err)
	}
	return &models.GeneralResponse{Message: fmt.Sprintf("Muvaffaqiyatli yaratildi ID %v   👍", id)}, nil
}

func (u *Database) GetUser(ctx context.Context, req *models.GetUserRequest) (*models.GetUserResponse, error) {
	query, args, err := db.Get(req)
	if err != nil {
		log.Println(err)
		return nil, errors.New("foydalanuvchi topilmadi 🤷‍♂️")
	}
	var res models.GetUserResponse
	if err := u.Db.QueryRow(query, args...).Scan(&res.ID, &res.Username, &res.Age, &res.Email, &res.LogOut); err != nil {
		log.Println(err)
		return nil, err
	}
	return &res, nil
}

func (u *Database) LastInserted(ctx context.Context, req *models.LastInsertedUser) (*models.GetUserResponse, error) {
	query, args, err := db.LastInserted()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var res models.GetUserResponse
	if err := u.Db.QueryRow(query, args...).Scan(&res.ID, &res.Username, &res.Age, &res.Email, &res.LogOut); err != nil {
		log.Println(err)
		return nil, err
	}
	return &res, nil
}

func (u *Database) UpdateUser(ctx context.Context, req *models.UpdateUserRequest) (*models.GeneralResponse, error) {
	req.Password = u.Hashing(req.Password)
	query, args, err := db.Update(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var id int

	if err := u.Db.QueryRow(query, args...).Scan(&id); err != nil {
		log.Println(err)
		return nil, err
	}
	_, err = u.N.Notification(ctx, &notification.ProduceMessage{UserId: int32(id), Message: fmt.Sprintf("you account is updated successfully with this is %v", id)})
	if err != nil {
		log.Println(err)
	}
	return &models.GeneralResponse{Message: fmt.Sprintf("Hisobingiz %v id si bilan yangilandi 🤓", id)}, nil
}

func (u *Database) LogOut(ctx context.Context, req *models.GetUserRequest) (*models.GeneralResponse, error) {
	query, args, err := db.LogOut(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	_, err = u.Db.Exec(query, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	_, err = u.N.Notification(ctx, &notification.ProduceMessage{UserId: int32(req.ID), Message: "hisobingizdan muvaffaqiyatli chiqdingiz, qayta kirishingiz mumkin! 👍"})
	if err != nil {
		log.Println(err)
	}
	return &models.GeneralResponse{Message: "Chiqish... 👉"}, nil
}

func (u *Database) DeletUser(ctx context.Context, req *models.GetUserRequest) (*models.GeneralResponse, error) {
	query, args, err := db.LogOut(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	_, err = u.Db.Exec(query, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	_, err = u.N.Notification(ctx, &notification.ProduceMessage{UserId: int32(req.ID), Message: "you have successfully deleted your account"})
	if err != nil {
		log.Println(err)
	}
	return &models.GeneralResponse{Message: "O'chirilmoqda... ☠️"}, nil
}

func (u *Database) ComparePassword(hashed, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (u *Database) Hashing(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return ""
	}
	return string(hashed)
}
