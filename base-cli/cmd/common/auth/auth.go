package auth

import (
	"os"
	"path"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/magaluCloud/mgccli/cmd/common/structs"
	"github.com/magaluCloud/mgccli/cmd/common/workspace"
	"gopkg.in/yaml.v3"
)

// access_key_id: ""
// access_token: ""
// current_environment: ""
// refresh_token: ""
// secret_access_key: ""

type AuthFile struct {
	AccessKeyID     string `yaml:"access_key_id"`
	AccessToken     string `yaml:"access_token"`
	RefreshToken    string `yaml:"refresh_token"`
	SecretAccessKey string `yaml:"secret_access_key"`
}

type Auth interface {
	GetAccessKeyID() string
	GetAccessToken() string
	GetRefreshToken() string
	GetSecretAccessKey() string

	SetAccessToken(token string) error
	SetRefreshToken(token string) error
	SetSecretAccessKey(key string) error
	SetAccessKeyID(key string) error

	ValidateToken() error
	RefreshToken() error
}

type authValue struct {
	authValue AuthFile
	workspace workspace.Workspace
}

func NewAuth(workspace workspace.Workspace) Auth {
	authFile := path.Join(workspace.Dir(), "auth.yaml")
	authContent, err := structs.LoadFileToStruct[AuthFile](authFile)
	if err != nil {
		//TODO: Handle error
		panic(err)
	}
	return &authValue{workspace: workspace, authValue: authContent}
}

func (a *authValue) GetAccessKeyID() string {
	return a.authValue.AccessKeyID
}

func (a *authValue) GetAccessToken() string {
	return a.authValue.AccessToken
}

func (a *authValue) GetRefreshToken() string {
	return a.authValue.RefreshToken
}

func (a *authValue) GetSecretAccessKey() string {
	return a.authValue.SecretAccessKey
}

func (a *authValue) SetAccessToken(token string) error {
	a.authValue.AccessToken = token
	return a.Write()
}

func (a *authValue) SetRefreshToken(token string) error {
	a.authValue.RefreshToken = token
	return a.Write()
}

func (a *authValue) SetSecretAccessKey(key string) error {
	a.authValue.SecretAccessKey = key
	return a.Write()
}

func (a *authValue) SetAccessKeyID(key string) error {
	a.authValue.AccessKeyID = key
	return a.Write()
}

func (a *authValue) Logout(name string) error {
	a.SetAccessToken("")
	a.SetRefreshToken("")
	a.SetSecretAccessKey("")
	a.SetAccessKeyID("")
	return a.Write()
}

func (a *authValue) Write() error {
	data, err := yaml.Marshal(a.authValue)
	if err != nil {
		return err
	}
	err = os.WriteFile(path.Join(a.workspace.Dir(), "auth.yaml"), data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (a *authValue) ValidateToken() error {
	//extract iat from token, if expires in less than 30 sec, run refresh operation
	token, err := jwt.Parse(a.authValue.AccessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.authValue.SecretAccessKey), nil
	})
	if err != nil {
		return err
	}
	iat := token.Claims.(jwt.MapClaims)["exp"].(float64)
	if iat < float64(time.Now().Unix()-60) {
		return a.RefreshToken()
	}
	return nil
}

func (a *authValue) RefreshToken() error {
	// not implemented yet
	return nil
}
