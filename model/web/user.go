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
	Angkatan *int    `json:"angkatan"`
}

type UserDetailResponse struct {
	Id        int     `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Image     *string `json:"image"`
	Phone     *string `json:"phone"`
	Address   *string `json:"address"`
	Provinsi  *string `json:"provinsi"`
	Kabupaten *string `json:"kabupaten"`
	Job       *string `json:"job"`
	Gender    *string `json:"gender"`
}

type UserFullResponse struct {
	Id        int     `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Provinsi  *string `json:"provinsi"`
	Kabupaten *string `json:"kabupaten"`
	Gender    *string `json:"gender"`
	Job       *string `json:"job"`
	Image     *string `json:"image"`
	Phone     *string `json:"phone"`
	Address   *string `json:"address"`
}
type UserResponsePassword struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserUpdatePasswordRequest struct {
	Id              int    `json:"id"`
	Password        string `validate:"required,min=1,max=100" json:"password"`
	PasswordConfirm string `validate:"required,min=1,max=100" json:"password_confirm"`
}

type UserUpdateRequest struct {
	Id        int    `validate:"required"`
	Name      string `validate:"required,max=60,min=1" json:"name"`
	Email     string `validate:"required,email,max=100,min=1" json:"email"`
	Phone     string `validate:"required,max=20,min=1" json:"phone"`
	Address   string `validate:"required,max=100,min=1" json:"address"`
	Job       string `validate:"required,max=100,min=1" json:"job"`
	Gender    string `validate:"required,max=50,min=1" json:"gender"`
	Provinsi  string `validate:"required,max=100,min=1" json:"provinsi"`
	Kabupaten string `validate:"required,max=100,min=1" json:"kabupaten"`
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
