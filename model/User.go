package model

import (
	"context"
	"ginblog/utils/errmsg"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
)

type User struct {
	gorm.Model
	// Username string `gorm:"type:varchar(20);not null " json:"username" validate:"required,min=4,max=12" label:"用户名"`
	Username string `gorm:"type:varchar(255);not null " json:"username" validate:"required,email" label:"用户名"`
	Password string `gorm:"type:varchar(500);not null" json:"password" validate:"required,min=6,max=120" label:"密码"`
	// 角色码：缺省为2 普通用户 要求角色码必须大于等于2 特权用户的角色码为1
	Role int `gorm:"type:int;DEFAULT:2" json:"role" validate:"required,gte=2" label:"角色码"`
	// 新添加的表字段
	Avatar   string `gorm:"type:varchar(200)" json:"avatar" label:"用户头像"`
	Nickname string `gorm:"type:varchar(20);not null " json:"nickname" label:"用户名"`
	Intro    string `gorm:"type:varchar(200)" json:"intro" label:"自我介绍"`
	Website  string `gorm:"type:varchar(200)" json:"webSite" label:"个人网站"`
}

// CheckUser 查询用户是否存在
func CheckUser(name string) (code int) {
	var user User
	db.Select("id").Where("username = ?", name).First(&user)
	if user.ID > 0 {
		return errmsg.ERROR_USERNAME_USED // 1001
	}
	return errmsg.SUCCSE
}

// CheckUpUser 更新查询
func CheckUpUser(id int, name string) (code int) {
	var user User
	db.Select("id, username").Where("username = ?", name).First(&user)
	// 如果编辑的用户名和自己的用户名相同则能够进行更新
	if user.ID == uint(id) {
		return errmsg.SUCCSE
	}
	// 如果查出来的用户ID大于0，说明用户名已经存在，返回错误
	if user.ID > 0 {
		return errmsg.ERROR_USERNAME_USED // 1001
	}
	return errmsg.SUCCSE
}

// CreateUser 新增用户
func CreateUser(data *User) int {
	//data.Password = ScryptPw(data.Password)
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR // 500
	}
	return errmsg.SUCCSE
}

// GetUser 查询用户
func GetUser(id int) (User, int) {
	var user User
	err := db.Limit(1).Where("ID = ?", id).Find(&user).Error
	if err != nil {
		return user, errmsg.ERROR
	}
	return user, errmsg.SUCCSE
}

// GetUsers 查询用户列表
func GetUsers(username string, pageSize int, pageNum int) ([]User, int64) {
	var users []User
	var total int64

	if username != "" {
		db.Select("id,username,role").Where(
			"username LIKE ?", username+"%",
		).Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users)
		db.Model(&users).Count(&total)
		return users, total
	}
	db.Select("id,username,role").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users)
	db.Model(&users).Count(&total)

	return users, total
}

// EditUser 编辑用户信息
func EditUser(id int, data *User) int {
	var user User
	var maps = make(map[string]interface{})
	maps["username"] = data.Username
	maps["role"] = data.Role
	err = db.Model(&user).Where("id = ? ", id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// ChangePassword 修改密码
func ChangePassword(id int, data *User) int {
	//var user User
	//var maps = make(map[string]interface{})
	//maps["password"] = data.Password

	err = db.Select("password").Where("id = ?", id).Updates(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// ResetPasswordFront 前端修改用户密码
func ResetPasswordFront(username string, password string) int {
	var (
		user User
		err  error
	)
	if err = db.Model(&User{}).Where("username = ?", username).First(&user).Error; err != nil {
		return errmsg.ERROR
	}
	if user.ID > 0 {
		user.Password = password
		if err = db.Select("password").Where("id = ?", user.ID).Updates(&user).Error; err != nil {
			return errmsg.ERROR
		}
		return errmsg.SUCCSE
	}
	return errmsg.ERROR_USERNAME_NOT_EXIST
}

// DeleteUser 删除用户
func DeleteUser(id int) int {
	var user User
	err = db.Where("id = ? ", id).Delete(&user).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// BeforeCreate 密码加密&权限控制
func (u *User) BeforeCreate(_ *gorm.DB) (err error) {
	u.Password = ScryptPw(u.Password)
	u.Role = 2
	return nil
}

func (u *User) BeforeUpdate(_ *gorm.DB) (err error) {
	u.Password = ScryptPw(u.Password)
	return nil
}

// ScryptPw 生成密码
func ScryptPw(password string) string {
	const cost = 10

	HashPw, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		log.Fatal(err)
	}

	return string(HashPw)
}

// CheckLogin 后台登录验证
func CheckLogin(username string, password string) (User, int) {
	var user User
	var PasswordErr error

	db.Where("username = ?", username).First(&user)

	// 如果用户不存在则直接返回
	if user.ID == 0 {
		return user, errmsg.ERROR_USER_NOT_EXIST
	}
	// 如果用户存在但角色权限不是1，为普通用户
	if user.Role != 1 {
		return user, errmsg.ERROR_USER_NO_RIGHT
	}
	// 验证哈希值
	PasswordErr = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if PasswordErr != nil {
		return user, errmsg.ERROR_PASSWORD_WRONG
	}

	return user, errmsg.SUCCSE
}

// CheckLoginFront 前台登录
func CheckLoginFront(username string, password string) (User, int) {
	var user User
	var PasswordErr error

	db.Select("id, username, password, intro, website, avatar, nickname, role").Where("username = ?", username).First(&user)

	PasswordErr = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if user.ID == 0 {
		return user, errmsg.ERROR_USER_NOT_EXIST
	}
	if PasswordErr != nil {
		return user, errmsg.ERROR_PASSWORD_WRONG
	}
	return user, errmsg.SUCCSE
}

func UpdateUserInfoFront(id int, m map[string]interface{}) int {
	if err := db.Model(&User{}).Where("id=?", id).Updates(m).Error; err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

type UserWithLikeSet struct {
	User
	ArticleLikeSet []int `json:"articleLikeSet"`
	CommentLikeSet []int `json:"commentLikeSet"`
}

// GetUserLikeSet 获取当前用户的点赞集合，包括文章点赞集合和评论点赞集合
func GetUserLikeSet(userInfo *User) *UserWithLikeSet {
	var (
		uwls         UserWithLikeSet
		artLikeSet   []int
		comLikeSet   []int
		artId, comId int
		rdh          *RedisHelper
		ctx          context.Context
		keys, res    []string
		key          string
		isMember     bool
		cursor       uint64
		err          error
	)
	uwls.User = *userInfo
	// 获取redis
	rdh = GetRedisHelper()
	ctx = context.Background()
	// 获取当前用户的文章点赞集合
	artLikeSet = make([]int, 0)
	cursor = 0
	for {
		if keys, cursor, err = rdh.Scan(ctx, cursor, "article:like:*", 10).Result(); err != nil {
			return nil
		}
		// 迭代keys中的每个键，判断用户是否在对应的点赞集合中
		for _, key = range keys {
			if isMember, err = rdh.SIsMember(ctx, key, userInfo.ID).Result(); err != nil {
				return nil
			}
			// 用户点赞了当前文章
			if isMember {
				res = strings.Split(key, ":")
				artId, _ = strconv.Atoi(res[len(res)-1])
				artLikeSet = append(artLikeSet, artId)
			}
		}
		if cursor == 0 {
			break
		}
	}
	// 获取当前用户的评论点赞集合
	comLikeSet = make([]int, 0)
	cursor = 0
	for {
		if keys, cursor, err = rdh.Scan(ctx, cursor, "comment:like:*", 10).Result(); err != nil {
			return nil
		}
		// 迭代keys中的每个键，判断用户是否在对应的点赞集合中
		for _, key = range keys {
			if isMember, err = rdh.SIsMember(ctx, key, userInfo.ID).Result(); err != nil {
				return nil
			}
			// 用户点赞了当前文章
			if isMember {
				res = strings.Split(key, ":")
				comId, _ = strconv.Atoi(res[len(res)-1])
				comLikeSet = append(comLikeSet, comId)
			}
		}
		if cursor == 0 {
			break
		}
	}
	uwls.ArticleLikeSet = artLikeSet
	uwls.CommentLikeSet = comLikeSet
	return &uwls
}
