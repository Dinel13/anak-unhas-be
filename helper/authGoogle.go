package helper

import "github.com/dinel13/anak-unhas-be/model/domain"

func NewGoogleClient(id, secret string) *domain.GoogleCred {
	return &domain.GoogleCred{
		ClientID:     id,
		ClientSecret: secret,
	}
}

// func VerifyGoogleToken(idToken string) error) {
// 	var httpClient = &http.Client{}
// 	oauth2Service, err := oauth2.New(httpClient)
// 	tokenInfoCall := oauth2Service.Tokeninfo()
// 	tokenInfoCall.IdToken(idToken)
// 	tokenInfo, err := tokenInfoCall.Do()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return tokenInfo, nil
// }
