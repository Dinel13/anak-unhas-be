package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dinel13/anak-unhas-be/model/domain"
	"github.com/dinel13/anak-unhas-be/model/web"
)

type userRepositoryImpl struct {
}

func NewUserRepository() domain.UserRepository {
	return &userRepositoryImpl{}
}

func (m *userRepositoryImpl) IsExits(ctx context.Context, tx *sql.DB, email string) (bool, error) {
	stmt := `SELECT id FROM users WHERE email = $1`

	var id int
	err := tx.QueryRowContext(ctx, stmt, email).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	if id > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (m *userRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, user web.UserCreateRequest) (*web.UserResponse, error) {
	var newUser web.UserResponse

	stmt := `INSERT INTO users (name, password, email) VALUES ($1, $2, $3) returning id, name`

	err := tx.QueryRowContext(ctx, stmt,
		user.Name,
		user.Password,
		user.Email,
	).Scan(
		&newUser.Id,
		&newUser.Name,
	)

	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

// Detail get all field user
func (m *userRepositoryImpl) Detail(ctx context.Context, tx *sql.DB, id int) (*web.UserFullResponse, error) {
	var user web.UserFullResponse

	stmt := `SELECT id, name, email, image, wa, jurusan, fakultas, address, bio, gender, angkatan, ig, tertarik 
	  FROM users WHERE id = $1`

	err := tx.QueryRowContext(ctx, stmt, id).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Image,
		&user.Wa,
		&user.Jurusan,
		&user.Fakultas,
		&user.Address,
		&user.Bio,
		&user.Gender,
		&user.Angkatan,
		&user.Ig,
		&user.Tertarik,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *userRepositoryImpl) GetByEmail(ctx context.Context, tx *sql.DB, email string) (*web.UserResponsePassword, error) {
	var user web.UserResponsePassword

	stmt := `SELECT id, name, password FROM users WHERE email = $1`

	err := tx.QueryRowContext(ctx, stmt, email).Scan(&user.Id, &user.Name, &user.Password)

	if err == sql.ErrNoRows {
		return nil, errors.New("email tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *userRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, user web.UserUpdateRequest) (*web.UserDetailResponse, error) {
	stmt := `UPDATE users
				SET name = $1, email = $2, phone = $3, job =$4, address = $5, gender = $6, provinsi =$7, kabupaten = $8
				WHERE id = $9
				RETURNING id, name, email, image, phone, job, address, gender, provinsi, kabupaten`

	// change string born to date
	// layoutFormat := "2006-01-02"
	// dateBorn, err := time.Parse(layoutFormat, user.Born)
	// if err != nil {
	// 	return nil, err
	// }

	row := tx.QueryRowContext(ctx, stmt,
		user.Name,
		user.Email,
		user.Phone,
		user.Job,
		user.Address,
		user.Gender,
		user.Provinsi,
		user.Kabupaten,
		user.Id,
	)

	var userUpdated web.UserDetailResponse
	err := row.Scan(
		&userUpdated.Id,
		&userUpdated.Name,
		&userUpdated.Email,
		&userUpdated.Image,
		&userUpdated.Phone,
		&userUpdated.Job,
		&userUpdated.Address,
		&userUpdated.Gender,
		&userUpdated.Provinsi,
		&userUpdated.Kabupaten,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &userUpdated, nil
}

func (m *userRepositoryImpl) UpdatePasword(ctx context.Context, tx *sql.Tx, user web.UserUpdatePasswordRequest) (*web.UserResponse, error) {
	stmt := `UPDATE users SET password = $1 WHERE id = $2 
				RETURNING id, name`

	row := tx.QueryRowContext(ctx, stmt, user.Password, user.Id)

	var userUpdated web.UserResponse
	err := row.Scan(
		&userUpdated.Id,
		&userUpdated.Name,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return &userUpdated, nil
}

// UpdateUserImage updates user image
func (m *userRepositoryImpl) UpdateImage(ctx context.Context, tx *sql.Tx, user web.UserUpdateImageRequest) error {
	stmt := `UPDATE users SET image = $1 WHERE id = $2`

	_, err := tx.ExecContext(ctx, stmt, user.Image, user.Id)

	if err == sql.ErrNoRows {
		return errors.New("user not found")
	}
	if err != nil {
		return err
	}

	return nil
}

func (m *userRepositoryImpl) GetImage(ctx context.Context, tx *sql.DB, id int) (*string, error) {
	stmt := `SELECT image FROM users WHERE id = $1`

	var image *string
	err := tx.QueryRowContext(ctx, stmt, id).Scan(&image)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	if image == nil {
		return nil, nil
	}

	return image, nil
}

func (m *userRepositoryImpl) GetPhone(ctx context.Context, tx *sql.DB, id int) (*string, error) {
	stmt := `SELECT phone FROM users WHERE id = $1`

	var phone *string
	err := tx.QueryRowContext(ctx, stmt, id).Scan(&phone)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return phone, nil
}

func (m *userRepositoryImpl) GetAddress(ctx context.Context, tx *sql.DB, id int) (*web.AddressResponse, error) {
	stmt := `SELECT provinsi, kabupaten kabupaten FROM users WHERE id = $1`

	var address web.AddressResponse
	err := tx.QueryRowContext(ctx, stmt, id).Scan(&address.Provinsi, &address.Kabupaten)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &address, nil
}

// SERACH BY NAME
func (m *userRepositoryImpl) Search(ctx context.Context, DB *sql.DB, query web.SearchRequest) ([]*web.UserSortResponse, error) {
	stmt := `SELECT id, name, image, jurusan, angkatan FROM users WHERE name LIKE $1 ORDER BY id ASC LIMIT 20 OFFSET ($2 - 1) * 20`

	rows, err := DB.QueryContext(ctx, stmt, "%"+query.Query+"%", query.Page)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*web.UserSortResponse

	for rows.Next() {
		var user web.UserSortResponse

		err = rows.Scan(
			&user.Id,
			&user.Name,
			&user.Image,
			&user.Jurusan,
			&user.Angkatan,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

// filter by name, jurusan, angkatan
func (m *userRepositoryImpl) Filter(ctx context.Context, DB *sql.DB, query web.FilterRequest) ([]*web.UserSortResponse, error) {
	stmt := `SELECT id, name, image, jurusan, angkatan FROM users WHERE name LIKE $1 AND jurusan LIKE $2 AND angkatan LIKE $3 ORDER BY id ASC LIMIT 20 OFFSET ($4 - 1) * 20`

	rows, err := DB.QueryContext(ctx, stmt, "%"+query.Name+"%", "%"+query.Jurusan+"%", "%"+query.Angkatan+"%", query.Page)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*web.UserSortResponse

	for rows.Next() {
		var user web.UserSortResponse

		err = rows.Scan(
			&user.Id,
			&user.Name,
			&user.Image,
			&user.Jurusan,
			&user.Angkatan,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

// count total result search
func (m *userRepositoryImpl) TotalResultSearch(ctx context.Context, DB *sql.DB, keyword string) (int, error) {
	stmt := `SELECT COUNT(*) FROM users WHERE name LIKE $1`

	var count int
	err := DB.QueryRowContext(ctx, stmt, "%"+keyword+"%").Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

// count total result filter
func (m *userRepositoryImpl) TotalResultFilter(ctx context.Context, DB *sql.DB, filter web.FilterRequest) (int, error) {
	stmt := `SELECT COUNT(*) FROM users WHERE name LIKE $1 AND jurusan LIKE $2 AND angkatan LIKE $3`

	var count int
	err := DB.QueryRowContext(ctx, stmt, "%"+filter.Name+"%", "%"+filter.Jurusan+"%", "%"+filter.Angkatan+"%").Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}
