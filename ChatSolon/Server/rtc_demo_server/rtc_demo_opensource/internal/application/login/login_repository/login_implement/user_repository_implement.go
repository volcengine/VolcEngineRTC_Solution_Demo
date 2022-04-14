package login_implement

import (
	"context"
	"time"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/util"

	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/application/login/login_entity"
	"github.com/volcengine/VolcEngineRTC_Solution_Demo/internal/pkg/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const UserProfileTable = "user_profile"

type UserRepositoryImpl struct {
}

func (impl *UserRepositoryImpl) Save(ctx context.Context, user *login_entity.UserProfile) error {
	defer util.CheckPanic()

	user.UpdatedAt = time.Now()
	err := db.Client.WithContext(ctx).Debug().Table(UserProfileTable).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			UpdateAll: true,
		}).Create(&user).Error
	return err
}

func (impl *UserRepositoryImpl) GetUser(ctx context.Context, userID string) (*login_entity.UserProfile, error) {
	defer util.CheckPanic()

	var rs *login_entity.UserProfile
	err := db.Client.WithContext(ctx).Debug().Table(UserProfileTable).Where("user_id = ?", userID).First(&rs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rs, nil
}
