package db

import (
	"context"
	"dreampicai/types"

	"github.com/google/uuid"
)

func CreateImage(image *types.Image) error {
	_, err := Bun.NewInsert().
		Model(image).
		Exec(context.Background())

	return err
}

func UpdateImage(image *types.Image) error {
	_, err := Bun.NewUpdate().
		Model(image).
		WherePK().
		Exec(context.Background())

	return err
}

func GetImagesByUserID(userID uuid.UUID) ([]types.Image, error) {
	var images []types.Image
	err := Bun.NewSelect().
		Model(&images).
		Where("deleted = ?", false).
		Where("user_id = ?", userID).
		Order("created_at desc").
		Scan(context.Background())
	return images, err
}
func GetImageByID(id uuid.UUID) (types.Image, error) {
	var images types.Image
	err := Bun.NewSelect().
		Model(&images).
		Where("id = ?", id).
		Scan(context.Background())
	return images, err
}

func CreateAccount(account *types.Account) error {
	_, err := Bun.NewInsert().
		Model(account).
		Exec(context.Background())
	return err
}

func UpdateAccount(account *types.Account) error {
	_, err := Bun.NewUpdate().
		Model(account).
		WherePK().
		Exec(context.Background())

	return err
}

func GetImagesByBatchID(batchID uuid.UUID) ([]types.Image, error) {
	var images []types.Image
	err := Bun.NewSelect().
		Model(&images).
		Where("batch_id = ?", batchID).
		Scan(context.Background())
	return images, err
}

func GetAccountByUserID(userID uuid.UUID) (types.Account, error) {
	account := types.Account{}
	err := Bun.NewSelect().
		Model(&account).
		Where("user_id = ?", userID).
		Scan(context.Background())

	return account, err
}
