package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

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
func (m *userRepositoryImpl) Detail(ctx context.Context, tx *sql.DB, id int) (*web.UserDetailResponse, error) {
	var user web.UserDetailResponse

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

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// GetNameAndImage get name and image for list friend in chat
func (m *userRepositoryImpl) GetNameAndImage(ctx context.Context, tx *sql.DB, id int) (*web.UserNameImage, error) {
	var user web.UserNameImage

	stmt := `SELECT name, image FROM users WHERE id = $1`

	err := tx.QueryRowContext(ctx, stmt, id).Scan(&user.Name, &user.Image)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *userRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, user web.UserUpdateRequest) (*web.UserDetailResponse, error) {
	stmt := `UPDATE users
				SET name =$2, wa = $3, angkatan = $4, address = $5, gender = $6, fakultas=$7, jurusan=$8, tertarik=$9, bio=$10
				WHERE id = $1
				RETURNING id, name, email, image, wa, jurusan, fakultas, address, bio, gender, angkatan, ig, tertarik `

	row := tx.QueryRowContext(ctx, stmt,
		user.Id,
		user.Name,
		user.Wa,
		user.Angkatan,
		user.Address,
		user.Gender,
		user.Fakultas,
		user.Jurusan,
		user.Tertarik,
		user.Bio,
	)

	var userUpdated web.UserDetailResponse
	err := row.Scan(
		&userUpdated.Id,
		&userUpdated.Name,
		&userUpdated.Email,
		&userUpdated.Image,
		&userUpdated.Wa,
		&userUpdated.Jurusan,
		&userUpdated.Fakultas,
		&userUpdated.Address,
		&userUpdated.Bio,
		&userUpdated.Gender,
		&userUpdated.Angkatan,
		&userUpdated.Ig,
		&userUpdated.Tertarik,
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
	stmt := `SELECT id, name, image, jurusan, angkatan FROM users WHERE LOWER(name) LIKE LOWER($1) ORDER BY id ASC LIMIT 80 OFFSET ($2 - 1) * 80`

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

	// v := reflect.ValueOf(query)
	// values := make([]interface{}, v.NumField())

	// for i := 0; i < v.NumField(); i++ {
	// 	values[i] = v.Field(i).Interface()
	// }

	stmt := `SELECT id, name, image, jurusan, angkatan FROM users 
	WHERE LOWER(name) LIKE LOWER($1) AND LOWER(jurusan) LIKE LOWER($2) 
	and LOWER(fakultas) LIKE LOWER($3) AND LOWER(angkatan) LIKE LOWER($4) 
	ORDER BY id ASC LIMIT 100 OFFSET ($5 - 1) * 100`

	// stmt := `SELECT id, name, image, jurusan, angkatan FROM users
	// WHERE name LIKE $1 AND jurusan LIKE $2 AND fakultas LIKE $3
	// AND angkatan LIKE $4 ORDER BY id ASC LIMIT 20 OFFSET ($5 - 1) * 20`

	rows, err := DB.QueryContext(ctx, stmt, "%"+query.Name+"%", "%"+query.Jurusan+"%", "%"+query.Fakultas+"%", "%"+query.Angkatan+"%", query.Page)

	if err != nil {
		log.Println("ggg", err)
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
	stmt := `SELECT COUNT(*) FROM users WHERE name LIKE $1 AND jurusan LIKE $2 AND jurusan LIKE $3 AND angkatan LIKE $4`

	var count int
	err := DB.QueryRowContext(ctx, stmt, "%"+filter.Name+"%", "%"+filter.Jurusan+"%", "%"+filter.Jurusan+"%", "%"+filter.Angkatan+"%").Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}
