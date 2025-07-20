package datasource

import (
	"database/sql"
	"fmt"
	"marketplace/internal/app"
)


type MarketRepo struct{
	db *sql.DB
}

func NewMarketRepo(db *sql.DB) *MarketRepo {
	return &MarketRepo{db: db}
}

func (s *MarketRepo) SaveAd(ad app.Ad) (app.Ad, error){
	stmt, err := s.db.Prepare(`INSERT INTO ads (uuid, title, description, price, img, user_uuid, created_at) 
	VALUES (?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return app.Ad{}, fmt.Errorf("prepare error DB:%w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(ad.UUID.String(), ad.Title, ad.Description, ad.Price, ad.ImageURL, ad.UserID, ad.CreatedAt)
	if err != nil {
		return app.Ad{}, fmt.Errorf("exec error DB:%w", err)
	}

	return ad, nil
}

func (s *MarketRepo) GetAdsList(params app.AdsListParams, user_id string) ([]app.Ad, error) {
	var ads []app.Ad

	query := `
		SELECT 
			a.id,
			a.uuid,
			a.title,
			a.description,
			a.img AS image_url,
			a.user_uuid,
			a.price,
			a.created_at
		FROM ads a
		JOIN users u ON a.user_uuid = u.uuid
		WHERE a.price >= ? AND a.price <= ?
	`

	sortBy := "a.created_at"
	if params.SortBy == "price" {
		sortBy = "a.price"
	}
	order := "ASC"
	if params.Order == "desc" {
		order = "DESC"
	}
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, order)

	offset := (params.Page - 1) * params.Limit
	query += " LIMIT ? OFFSET ?"

	rows, err := s.db.Query(query, params.MinPrice, params.MaxPrice, params.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query error DB: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ad app.Ad
		err := rows.Scan(
			&ad.ID,
			&ad.UUID,
			&ad.Title,
			&ad.Description,
			&ad.ImageURL,
			&ad.UserID,
			&ad.Price,
			&ad.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error DB: %w", err)
		}
		if ad.UserID.String() == user_id {
			ad.Owner = true
			
		}

		ads = append(ads, ad)
	}

	return ads, nil
}