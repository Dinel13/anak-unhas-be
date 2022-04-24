package web

type UserCreateRequest struct {
	Name     string `validate:"required,min=1,max=100" json:"name"`
	Password string `validate:"required,min=1,max=100" json:"password"`
	Email    string `validate:"required,min=1,max=100" json:"email"`
}

type UserLoginRequest struct {
	Password string `validate:"required,min=1,max=100" json:"password"`
	Email    string `validate:"required,min=1,max=100" json:"email"`
}

// for outh google
type UserAuthGoogle struct {
	TokenId  string `validate:"required" json:"tokenId"`
	Email    string `validate:"required,min=1,max=100" json:"email"`
	GoogleId string `validate:"required,min=1,max=100" json:"googleId"`
}

type UserResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

type UserSortResponse struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Image    *string `json:"image"`
	Jurusan  *string `json:"jurusan"`
	Angkatan *string `json:"angkatan"`
}

type UserDetailResponse struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Gender   *string `json:"gender"`
	Wa       *string `json:"wa"`
	Image    *string `json:"image"`
	Jurusan  *string `json:"jurusan"`
	Fakultas *string `json:"fakultas"`
	Address  *string `json:"address"`
	Bio      *string `json:"bio"`
	Angkatan *string `json:"angkatan"`
	Ig       *string `json:"ig"`
	Tertarik *string `json:"tertarik"`
}
type UserResponsePassword struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserNameImage struct {
	Name  string  `json:"name"`
	Image *string `json:"image"`
}

type UserUpdatePasswordRequest struct {
	Id              int    `json:"id"`
	Password        string `validate:"required,min=1,max=100" json:"password"`
	PasswordConfirm string `validate:"required,min=1,max=100" json:"password_confirm"`
}

type UserUpdateRequest struct {
	Id       int    `validate:"required"`
	Name     string `validate:"required,max=60,min=1" json:"name"`
	Gender   string `json:"gender"`
	Wa       string `json:"wa"`
	Jurusan  string `json:"jurusan"`
	Fakultas string `json:"fakultas"`
	Address  string `json:"address"`
	Bio      string `json:"bio"`
	Angkatan string `json:"angkatan"`
	Ig       string `json:"ig"`
	Tertarik string `json:"tertarik"`
}

type UserUpdateImageRequest struct {
	Id    int    `validate:"required"`
	Image string `validate:"required,max=150,min=1" json:"image"`
}

type UserForgetPasswordRequest struct {
	Email string `validate:"required,email,max=100,min=1" json:"email"`
}

type AddressResponse struct {
	Provinsi  *string `json:"provinsi"`
	Kabupaten *string `json:"kabupaten"`
}

type FilterRequest struct {
	Name     string `json:"name"`
	Fakultas string `json:"fakultas"`
	Jurusan  string `json:"jurusan"`
	Angkatan string `json:"angkatan"`
	Skill    string `json:"skill"`
	Page     int    `json:"page"`
}

type SearchRequest struct {
	Query string `json:"query"`
	Page  int    `json:"page"`
}

type FilterResponse struct {
	Total int                 `json:"total"`
	Users []*UserSortResponse `json:"users"`
}

type SearchResponse struct {
	Total int                 `json:"total"`
	Users []*UserSortResponse `json:"users"`
}

type UserCreateEs struct {
	Id       int
	Name     string
	Gender   string
	Jurusan  string
	Fakultas string
	Address  string
	Angkatan string
	Tertarik string
}
