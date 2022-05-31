package users

import (
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/githubzjm/tuo/api/v1/users/def"
	"github.com/githubzjm/tuo/internal/pkg/jwt"
	"github.com/githubzjm/tuo/internal/pkg/mysql"
	"github.com/githubzjm/tuo/internal/pkg/mysql/models"
)

func CreateUser(value *models.User) (isConflict bool, err error) {
	// TODO may create a init cluster automatically
	return mysql.Create(value)
}

// if not exists, return nil
// if exists, return the user info
func QueryUser(query *models.User) (*models.User, error) {
	var user models.User
	isExist, err := mysql.Query(query, &user)
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, nil
	}
	return &user, nil
}
func UpdateUser(query, values *models.User) error {
	return mysql.Update(query, values)
}
func DeleteUser(query *models.User) error {
	return mysql.Delete(query, &models.User{})
}

// delete association records
func DeleteUserWithAssociation(query *models.User) error {
	return mysql.DeleteWithAssociatons(query)
}

func GenerateToken(c *gin.Context, user *models.User) (string, error) {
	j := jwt.NewJWT()

	// init claim of payload
	claims := jwt.CustomClaims{
		UserID: user.ID,
		StandardClaims: jwtgo.StandardClaims{
			NotBefore: int64(time.Now().Unix()),                                                  // token active time
			ExpiresAt: int64(time.Now().Unix() + int64(def.TOKEN_VALID_DURATION_NS/time.Second)), // token expire time
			Issuer:    "tuo",
		},
	}
	// create token
	token, err := j.CreateToken(claims)
	if err != nil {
		return "", err
	}

	return token, nil
}

// TODO just add a admin user, can do anything
func CheckPrivilege(visitor uint, target uint) (bool, error) {
	if visitor != target {
		// for now, user can only operate its own info
		// may adapt privilege system later
		return false, nil
	} else {
		return true, nil
	}
}
