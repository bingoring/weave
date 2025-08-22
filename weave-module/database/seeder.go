package database

import (
	"fmt"
	"time"

	"weave-module/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// SeedData creates initial data for development
func SeedData() error {
	if DB == nil {
		return fmt.Errorf("database connection not established")
	}

	// Check if data already exists
	var userCount int64
	DB.Model(&models.User{}).Count(&userCount)
	if userCount > 0 {
		fmt.Println("Seed data already exists, skipping...")
		return nil
	}

	fmt.Println("Creating seed data...")

	// Create seed users
	users, err := createSeedUsers()
	if err != nil {
		return fmt.Errorf("failed to create seed users: %w", err)
	}

	// Create seed channels
	channels, err := createSeedChannels(users)
	if err != nil {
		return fmt.Errorf("failed to create seed channels: %w", err)
	}

	// Create seed weaves
	weaves, err := createSeedWeaves(users, channels)
	if err != nil {
		return fmt.Errorf("failed to create seed weaves: %w", err)
	}

	// Create seed collaborations
	err = createSeedCollaborations(users, weaves)
	if err != nil {
		return fmt.Errorf("failed to create seed collaborations: %w", err)
	}

	fmt.Println("Seed data created successfully!")
	return nil
}

func createSeedUsers() ([]models.User, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	
	users := []models.User{
		{
			ID:           uuid.New(),
			Username:     "alice_creator",
			Email:        "alice@weave.local",
			PasswordHash: string(hashedPassword),
			Bio:          stringPtr("혁신적인 아이디어를 통해 세상을 바꾸고 싶은 창작자입니다."),
			IsVerified:   true,
			IsActive:     true,
		},
		{
			ID:           uuid.New(),
			Username:     "bob_developer",
			Email:        "bob@weave.local",
			PasswordHash: string(hashedPassword),
			Bio:          stringPtr("기술과 디자인의 완벽한 조화를 추구하는 개발자입니다."),
			IsVerified:   true,
			IsActive:     true,
		},
		{
			ID:           uuid.New(),
			Username:     "charlie_designer",
			Email:        "charlie@weave.local",
			PasswordHash: string(hashedPassword),
			Bio:          stringPtr("사용자 경험을 최우선으로 생각하는 UX/UI 디자이너입니다."),
			IsVerified:   true,
			IsActive:     true,
		},
		{
			ID:           uuid.New(),
			Username:     "diana_researcher",
			Email:        "diana@weave.local",
			PasswordHash: string(hashedPassword),
			Bio:          stringPtr("데이터 기반의 인사이트로 미래를 예측하는 연구자입니다."),
			IsVerified:   true,
			IsActive:     true,
		},
	}

	for i := range users {
		if err := DB.Create(&users[i]).Error; err != nil {
			return nil, err
		}

		// Create user profile
		profile := models.UserProfile{
			UserID: users[i].ID,
		}
		if err := DB.Create(&profile).Error; err != nil {
			return nil, err
		}

		// Create user analytics
		analytics := models.UserAnalytics{
			UserID: users[i].ID,
		}
		if err := DB.Create(&analytics).Error; err != nil {
			return nil, err
		}
	}

	// Create some follow relationships
	follows := []models.UserFollow{
		{FollowerID: users[0].ID, FollowingID: users[1].ID},
		{FollowerID: users[0].ID, FollowingID: users[2].ID},
		{FollowerID: users[1].ID, FollowingID: users[0].ID},
		{FollowerID: users[1].ID, FollowingID: users[3].ID},
		{FollowerID: users[2].ID, FollowingID: users[0].ID},
		{FollowerID: users[3].ID, FollowingID: users[1].ID},
	}

	for _, follow := range follows {
		if err := DB.Create(&follow).Error; err != nil {
			return nil, err
		}
	}

	return users, nil
}

func createSeedChannels(users []models.User) ([]models.Channel, error) {
	channels := []models.Channel{
		{
			ID:          uuid.New(),
			Name:        "Technology",
			Description: stringPtr("최신 기술 동향과 혁신적인 아이디어를 공유하는 공간"),
			IsPublic:    true,
		},
		{
			ID:          uuid.New(),
			Name:        "Design",
			Description: stringPtr("창의적인 디자인과 사용자 경험에 대한 통찰을 나누는 곳"),
			IsPublic:    true,
		},
		{
			ID:          uuid.New(),
			Name:        "Research",
			Description: stringPtr("학술적 연구와 데이터 분석 결과를 공유하는 채널"),
			IsPublic:    true,
		},
		{
			ID:          uuid.New(),
			Name:        "Startup Ideas",
			Description: stringPtr("혁신적인 스타트업 아이디어와 비즈니스 모델을 논의하는 공간"),
			IsPublic:    true,
		},
	}

	for i := range channels {
		if err := DB.Create(&channels[i]).Error; err != nil {
			return nil, err
		}
	}

	return channels, nil
}

func createSeedWeaves(users []models.User, channels []models.Channel) ([]models.Weave, error) {
	weaves := []models.Weave{
		{
			ID:                  uuid.New(),
			UserID:              users[0].ID,
			ChannelID:           channels[0].ID,
			Title:               "AI 기반 협업 플랫폼의 미래",
			Description:         stringPtr("인공지능이 창작 과정에 미치는 영향과 협업의 새로운 패러다임"),
			Content:             `{"type":"doc","content":[{"type":"heading","attrs":{"level":1},"content":[{"type":"text","text":"AI 기반 협업 플랫폼의 미래"}]},{"type":"paragraph","content":[{"type":"text","text":"인공지능 기술의 발전은 우리의 협업 방식을 근본적으로 변화시키고 있습니다. 이 문서에서는 AI가 창작 과정에 미치는 영향과 앞으로의 가능성에 대해 탐구해보겠습니다."}]}]}`,
			Status:              models.WeaveStatusPublished,
			Type:                models.WeaveTypeOriginal,
			IsCollaborationOpen: true,
			PublishedAt:         timePtr(time.Now().Add(-24 * time.Hour)),
		},
		{
			ID:                  uuid.New(),
			UserID:              users[1].ID,
			ChannelID:           channels[1].ID,
			Title:               "사용자 중심 디자인 시스템 구축",
			Description:         stringPtr("확장 가능하고 일관성 있는 디자인 시스템을 만드는 방법론"),
			Content:             `{"type":"doc","content":[{"type":"heading","attrs":{"level":1},"content":[{"type":"text","text":"사용자 중심 디자인 시스템 구축"}]},{"type":"paragraph","content":[{"type":"text","text":"현대적인 웹 애플리케이션에서 일관성 있는 사용자 경험을 제공하기 위해서는 체계적인 디자인 시스템이 필요합니다."}]}]}`,
			Status:              models.WeaveStatusPublished,
			Type:                models.WeaveTypeOriginal,
			IsCollaborationOpen: true,
			PublishedAt:         timePtr(time.Now().Add(-18 * time.Hour)),
		},
		{
			ID:                  uuid.New(),
			UserID:              users[2].ID,
			ChannelID:           channels[2].ID,
			Title:               "빅데이터 분석을 통한 사용자 행동 패턴 연구",
			Description:         stringPtr("대규모 데이터셋을 활용한 사용자 행동 예측 모델 개발"),
			Content:             `{"type":"doc","content":[{"type":"heading","attrs":{"level":1},"content":[{"type":"text","text":"빅데이터 분석을 통한 사용자 행동 패턴 연구"}]},{"type":"paragraph","content":[{"type":"text","text":"사용자의 디지털 발자취를 분석하여 개인화된 서비스를 제공하는 것은 현대 IT 서비스의 핵심입니다."}]}]}`,
			Status:              models.WeaveStatusPublished,
			Type:                models.WeaveTypeOriginal,
			IsCollaborationOpen: true,
			PublishedAt:         timePtr(time.Now().Add(-12 * time.Hour)),
		},
	}

	for i := range weaves {
		if err := DB.Create(&weaves[i]).Error; err != nil {
			return nil, err
		}

		// Create weave analytics
		analytics := models.WeaveAnalytics{
			WeaveID:     weaves[i].ID,
			TotalViews:  10 + i*5,
			UniqueViews: 8 + i*3,
			TotalLikes:  2 + i,
		}
		if err := DB.Create(&analytics).Error; err != nil {
			return nil, err
		}

		// Create timeline entry
		timeline := models.WeaveTimeline{
			WeaveID:     weaves[i].ID,
			UserID:      weaves[i].UserID,
			EventType:   models.TimelineCreated,
			Title:       "Weave 생성됨",
			Description: stringPtr("새로운 Weave가 생성되었습니다."),
		}
		if err := DB.Create(&timeline).Error; err != nil {
			return nil, err
		}
	}

	return weaves, nil
}

func createSeedCollaborations(users []models.User, weaves []models.Weave) error {
	// Create some contributions
	contributions := []models.Contribution{
		{
			UserID:          users[1].ID,
			WeaveID:         weaves[0].ID,
			Type:            models.ContributionTypeSuggestion,
			Title:           "AI 윤리적 고려사항 추가",
			Description:     stringPtr("AI 기술 발전과 함께 고려해야 할 윤리적 측면들을 추가하면 좋을 것 같습니다."),
			ProposedContent: stringPtr(`{"type":"paragraph","content":[{"type":"text","text":"AI 기술의 발전과 함께 개인정보 보호, 알고리즘 편향성, 투명성 등의 윤리적 고려사항들이 중요해지고 있습니다."}]}`),
			Status:          models.ContributionStatusPending,
		},
		{
			UserID:          users[2].ID,
			WeaveID:         weaves[1].ID,
			Type:            models.ContributionTypeContentEdit,
			Title:           "접근성 가이드라인 추가",
			Description:     stringPtr("WCAG 접근성 가이드라인을 반영한 내용을 추가했습니다."),
			ProposedContent: stringPtr(`{"type":"paragraph","content":[{"type":"text","text":"웹 접근성 가이드라인(WCAG)을 준수하여 모든 사용자가 동등하게 서비스를 이용할 수 있도록 해야 합니다."}]}`),
			Status:          models.ContributionStatusAccepted,
		},
	}

	for i := range contributions {
		if err := DB.Create(&contributions[i]).Error; err != nil {
			return err
		}
	}

	// Create some lab comments
	comments := []models.LabComment{
		{
			UserID:  users[3].ID,
			WeaveID: weaves[0].ID,
			Type:    models.CommentTypeQuestion,
			Content: "AI 기반 협업에서 사람의 창의성은 어떻게 보장될 수 있을까요?",
		},
		{
			UserID:  users[0].ID,
			WeaveID: weaves[1].ID,
			Type:    models.CommentTypeSuggestion,
			Content: "컴포넌트 라이브러리 예시를 추가하면 더 실용적일 것 같습니다.",
		},
	}

	for i := range comments {
		if err := DB.Create(&comments[i]).Error; err != nil {
			return err
		}
	}

	// Create some likes
	likes := []models.WeaveLike{
		{UserID: users[1].ID, WeaveID: weaves[0].ID},
		{UserID: users[2].ID, WeaveID: weaves[0].ID},
		{UserID: users[0].ID, WeaveID: weaves[1].ID},
		{UserID: users[3].ID, WeaveID: weaves[1].ID},
		{UserID: users[0].ID, WeaveID: weaves[2].ID},
	}

	for i := range likes {
		if err := DB.Create(&likes[i]).Error; err != nil {
			return err
		}
	}

	return nil
}

func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}