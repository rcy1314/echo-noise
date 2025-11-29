package services

import (
	"errors"
	"fmt"

	"github.com/lin-snow/ech0/internal/dto"
	"github.com/lin-snow/ech0/internal/models"
	"github.com/lin-snow/ech0/internal/repository"
	"github.com/lin-snow/ech0/pkg"
	"golang.org/x/crypto/bcrypt"
)

func Register(userdto dto.RegisterDto) error {
	if userdto.Username == "" || userdto.Password == "" {
		return errors.New(models.UsernameOrPasswordCannotBeEmptyMessage)
	}

	// 使用 bcrypt 存储新用户密码
	hashed := models.HashPassword(userdto.Password)
	if hashed == "" {
		return errors.New("密码加密失败")
	}

	newuser := models.User{
		Username: userdto.Username,
		Password: hashed,
		IsAdmin:  false,
		Token:    models.GenerateToken(32),
	}

	user, err := repository.GetUserByUsername(userdto.Username)
	if err == nil && user != nil && user.ID != 0 {
		return errors.New(models.UsernameAlreadyExistsMessage)
	}

	users, err := repository.GetAllUsers()
	if err != nil {
		return errors.New(models.GetAllUsersFailMessage)
	}
	if len(users) == 0 {
		newuser.IsAdmin = true
	}

	if err := repository.CreateUser(&newuser); err != nil {
		return errors.New(models.CreateUserFailMessage)
	}

	return nil
}

func Login(userdto dto.LoginDto) (*models.User, error) {
	if userdto.Username == "" || userdto.Password == "" {
		return nil, errors.New(models.UsernameOrPasswordCannotBeEmptyMessage)
	}

	plain := userdto.Password
	md5pwd := pkg.MD5Encrypt(plain)

	user, err := repository.GetUserByUsername(userdto.Username)
	if err != nil {
		return nil, errors.New(models.UserNotFoundMessage)
	}

	// 兼容两种密码存储：MD5 与 bcrypt
	// 1) 直接匹配 MD5
	if user.Password == md5pwd {
		// ok
	} else {
		// 2) 尝试按 bcrypt 校验
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plain)); err != nil {
			return nil, errors.New(models.PasswordIncorrectMessage)
		}
	}

	// 只在 token 为空时生成新的 token
	if user.Token == "" {
		user.Token = models.GenerateToken(32)
		if err := repository.UpdateUser(user); err != nil {
			return nil, fmt.Errorf("生成用户 token 失败: %v", err)
		}
	}

	return user, nil
}

func GetStatus() (models.Status, error) {
	sysuser, err := repository.GetSysAdmin()
	if err != nil {
		return models.Status{}, errors.New(models.UserNotFoundMessage)
	}

	var users []models.UserStatus
	allusers, err := repository.GetAllUsers()
	if err != nil {
		return models.Status{}, errors.New(models.GetAllUsersFailMessage)
	}
	for _, user := range allusers {
		users = append(users, models.UserStatus{
			ID:       user.ID,
			Username: user.Username,
			IsAdmin:  user.IsAdmin,
		})
	}

	status := models.Status{}

	messages, err := repository.GetAllMessages(true)
	if err != nil {
		return status, errors.New(models.GetAllMessagesFailMessage)
	}

	status.SysAdminID = sysuser.ID
	status.Username = sysuser.Username
	status.Users = users
	status.TotalMessages = len(messages)

	return status, nil
}

func GetUserByID(userID uint) (*models.User, error) {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return nil, errors.New(models.UserNotFoundMessage)
	}
	return user, nil
}

func IsUserAdmin(userID uint) (bool, error) {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return false, errors.New(models.UserNotFoundMessage)
	}
	return user.IsAdmin, nil
}

func UpdateUser(user *models.User, userdto dto.UserInfoDto) error {
	if user.Username == userdto.Username {
		return nil
	}

	if userdto.Username == "" {
		return errors.New(models.UsernameCannotBeEmptyMessage)
	}

	user.Username = userdto.Username
	if err := repository.UpdateUser(user); err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func ChangePassword(user *models.User, userdto dto.UserInfoDto) error {
	if user == nil {
		return errors.New("用户信息不能为空")
	}

	if userdto.Password == "" {
		return errors.New(models.PasswordCannotBeEmptyMessage)
	}

	// 如果新密码与旧密码一致，则拒绝
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userdto.Password)); err == nil {
		return errors.New(models.PasswordCannotBeSameAsBeforeMessage)
	}

	// 使用 bcrypt 更新密码
	hashed := models.HashPassword(userdto.Password)
	if hashed == "" {
		return fmt.Errorf("密码加密失败")
	}
	user.Password = hashed

	if err := repository.UpdateUser(user); err != nil {
		return fmt.Errorf("更新密码失败: %v", err)
	}

	return nil
}

func UpdateUserAdmin(userID uint, currentUserID uint) error {
    user, err := repository.GetUserByID(userID)
    if err != nil {
        return err
    }
    // 不允许取消当前登录用户的管理员身份
    if userID == currentUserID && user.IsAdmin {
        return fmt.Errorf("不允许取消当前登录用户的管理员身份")
    }
    // 至少保留一位管理员
    if user.IsAdmin {
        count, err := repository.CountAdmins()
        if err != nil { return err }
        if count <= 1 {
            return fmt.Errorf("系统至少保留一位管理员")
        }
    }
    user.IsAdmin = !user.IsAdmin
    if err := repository.UpdateUser(user); err != nil {
        return err
    }
    return nil
}

func GetUserByUsername(username string) (*models.User, error) {
    user, err := repository.GetUserByUsername(username)
    if err != nil {
        return nil, err
    }
    return user, nil
}
