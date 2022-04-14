package login_service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/login/login_entity"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/login/login_repository/login_facade"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/models/custom_error"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/config"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/logs"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/redis_cli/general"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/redis_cli/lock"
)

const (
	chars             = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	userNameLen       = 6
	maxRetryTimes     = 3
	retryBackOff      = 8 * time.Millisecond
	localUserIDLength = 8
	LocalUserIDPrefix = "8181"
	TokenExpiration   = 24 * 7 * time.Hour
)

const (
	retryDelay    = 8 * time.Millisecond
	maxRetryDelay = 128 * time.Millisecond
	maxRetryNum   = 10
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type UserService struct {
	userRepo       login_facade.UserRepositoryInterface
	loginTokenRepo login_facade.LoginTokenRepositoryInterface
}

var userService *UserService

func GetUserService() *UserService {
	if userService == nil {
		userService = &UserService{
			userRepo:       login_facade.GetUserRepository(),
			loginTokenRepo: login_facade.GetLoginTokenRepository(),
		}
	}
	return userService
}

func (s *UserService) Login(ctx context.Context, userID, token string) error {
	err := s.loginTokenRepo.Save(ctx, &login_entity.LoginToken{
		Token:      token,
		UserID:     userID,
		CreateTime: time.Now(),
	})
	if err != nil {
		return custom_error.InternalError(err)
	}
	return nil
}

func (s *UserService) GetUserID(ctx context.Context, token string) string {
	return s.loginTokenRepo.GetUserID(ctx, token)
}

func (s *UserService) GetUserName(ctx context.Context, userID string) (string, error) {
	user, err := s.userRepo.GetUser(ctx, userID)
	if user == nil {
		logs.CtxWarn(ctx, "user not exist")
		return "", nil
	}

	if err != nil {
		logs.CtxError(ctx, "failed to get user, err: %v", err)
		return "", custom_error.InternalError(err)
	}

	logs.CtxInfo(ctx, "get user name: %v", user.UserName)

	return user.UserName, nil
}

func (s *UserService) SetUserName(ctx context.Context, userID, userName string) error {
	user := &login_entity.UserProfile{
		UserID:   userID,
		UserName: userName,
	}

	err := s.userRepo.Save(ctx, user)
	if err != nil {
		logs.CtxError(ctx, "failed to set user name, err: %v", err)
		return custom_error.InternalError(err)
	}

	logs.CtxInfo(ctx, "set user name: %v", userName)
	return nil
}

func (s *UserService) GetLoginUserName(ctx context.Context, userID string) string {
	userName, err := s.GetUserName(ctx, userID)
	if err != nil {
		return GetRandomUserName()
	}

	if userName == "" {
		userName = GetRandomUserName()

		_ = s.SetUserName(ctx, userID, userName)
	}

	return userName
}

func (s *UserService) CheckLoginToken(ctx context.Context, token string) error {
	if !config.Configs().CheckLoginToken {
		logs.CtxWarn(ctx, "no need to check login token")
		return nil
	}

	if token == "" {
		logs.CtxWarn(ctx, "empty token")
		return custom_error.ErrorTokenEmpty
	}

	var exist bool
	var err error

	for i := 0; i < maxRetryTimes; i++ {
		exist, err = s.loginTokenRepo.ExistToken(ctx, token)
		if err == nil {
			break
		}

		if i < maxRetryTimes {
			retryDelay := time.Duration(math.Pow(2, float64(i))) * retryBackOff
			time.Sleep(retryDelay)
		}
	}

	if err != nil {
		logs.CtxError(ctx, "get token failed, err: %v", err)
		return custom_error.InternalError(err)
	}

	if !exist {
		logs.CtxWarn(ctx, "login token expiry")
		return custom_error.ErrorTokenExpiry
	}

	return nil
}

// RefreshLoginToken avoid token expiry when user in room.
func (s *UserService) RefreshLoginToken(ctx context.Context, token string) error {
	if !config.Configs().CheckLoginToken {
		logs.CtxWarn(ctx, "no need to check login token")
		return nil
	}

	if token == "" {
		logs.CtxWarn(ctx, "empty token")
		return custom_error.ErrorTokenEmpty
	}

	var createdAt int64
	var err error

	for i := 0; i < maxRetryTimes; i++ {
		createdAt, err = s.loginTokenRepo.GetTokenCreatedAt(ctx, token)
		if err == nil {
			break
		}

		if i < maxRetryTimes {
			retryDelay := time.Duration(math.Pow(2, float64(i))) * retryBackOff
			time.Sleep(retryDelay)
		}
	}

	if err != nil {
		logs.CtxError(ctx, "get token failed, err: %v", err)
		return custom_error.InternalError(err)
	}

	now := time.Now().UnixNano()

	if now-createdAt > int64(TokenExpiration) {
		logs.CtxWarn(ctx, "login token expiry")
		return custom_error.ErrorTokenExpiry
	}

	s.loginTokenRepo.SetTokenExpiration(ctx, token, TokenExpiration)

	return nil
}

func (s *UserService) IsAuditorLogin(_ context.Context, phoneNo, code string) bool {
	if config.Configs().AuditorPhoneCode == nil {
		return false
	}

	val, ok := config.Configs().AuditorPhoneCode[phoneNo]
	if ok && val == code {
		return true
	}
	return false
}

func (s *UserService) GenerateLocalUserIDWithRetry(ctx context.Context) (string, error) {
	userID, err := s.generateLocalUserID(ctx)
	for i := 0; userID == 0 && i <= maxRetryNum; i++ {
		backOff := time.Duration(int64(math.Pow(2, float64(i)))) * retryDelay
		if backOff > maxRetryDelay {
			backOff = maxRetryDelay
		}
		time.Sleep(backOff)
		userID, err = s.generateLocalUserID(ctx)
	}
	if userID == 0 {
		logs.CtxError(ctx, "failed to generate userID, err: %s", err)
		return "", custom_error.InternalError(errors.New("make user err"))
	}
	return strconv.FormatInt(userID, 10), nil
}

func (s *UserService) generateLocalUserID(ctx context.Context) (int64, error) {
	ok, lt := lock.LockLocalUserIDAssign(ctx)
	if !ok {
		return 0, custom_error.ErrLockRedis
	}

	defer lock.UnLockLocalUserIDAssign(ctx, lt)

	userID, err := general.GetGeneratedUserID(ctx)
	if err != nil {
		return 0, custom_error.InternalError(err)
	}

	baseline := int64(math.Pow10(localUserIDLength))
	minUserID := int64(math.Pow10(localUserIDLength - 1))

	if userID == 0 {
		userID = time.Now().Unix() % baseline
	} else {
		userID = (userID + 1) % baseline
	}

	if userID < minUserID {
		userID += minUserID
	}

	general.SetGeneratedUserID(ctx, userID)

	return userID, nil
}

func (s *UserService) GenerateLocalLoginToken(_ context.Context, userID string, createdTime int64) string {
	strCreateTime := strconv.FormatInt(createdTime, 10)
	text := userID + strCreateTime

	hasher := md5.New()
	hasher.Write([]byte(text))

	return hex.EncodeToString(hasher.Sum(nil))
}

func GetRandomUserName() string {
	return RandString(userNameLen)
}

func RandString(l uint) string {
	s := make([]byte, l)

	for i := 0; i < int(l); i++ {
		s[i] = chars[rand.Intn(len(chars))]
	}

	return string(s)
}
