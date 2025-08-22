-- Weave Database Initialization Script

-- Enable TimescaleDB extension
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- Enable UUID extension for generating UUIDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create schemas
CREATE SCHEMA IF NOT EXISTS weave;
CREATE SCHEMA IF NOT EXISTS analytics;

-- Set default schema
SET search_path TO weave, public;

-- Create basic tables (will be managed by GORM migrations)
-- This is just for initial setup

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    profile_image VARCHAR(500),
    bio TEXT,
    is_verified BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Channels table
CREATE TABLE IF NOT EXISTS channels (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    cover_image VARCHAR(500),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Weaves table (main content)
CREATE TABLE IF NOT EXISTS weaves (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    channel_id UUID NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    cover_image VARCHAR(500),
    content JSONB NOT NULL,
    version INTEGER DEFAULT 1,
    parent_weave_id UUID REFERENCES weaves(id) ON DELETE SET NULL,
    is_published BOOLEAN DEFAULT FALSE,
    is_featured BOOLEAN DEFAULT FALSE,
    view_count INTEGER DEFAULT 0,
    like_count INTEGER DEFAULT 0,
    fork_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create hypertable for time-series data (analytics)
CREATE TABLE IF NOT EXISTS analytics.weave_events (
    time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    weave_id UUID NOT NULL,
    user_id UUID,
    event_type VARCHAR(50) NOT NULL,
    event_data JSONB,
    ip_address INET,
    user_agent TEXT
);

-- Convert to hypertable for time-series optimization
SELECT create_hypertable('analytics.weave_events', 'time', if_not_exists => TRUE);

-- Insert default channels
INSERT INTO channels (name, slug, description) VALUES
('Recipes', 'recipes', '나만의 비밀 레시피 공유 및 발전'),
('Travel Plans', 'travel-plans', '모두가 함께 만드는 완벽한 여행 계획'),
('Workout Routines', 'workout-routines', '헬스, 홈트, 요가 등 목표별 운동 루틴'),
('Playlists', 'playlists', '상황과 기분에 맞는 최고의 플레이리스트'),
('Interior Tips', 'interior-tips', '셀프 인테리어 팁과 아이디어 모음'),
('Reading Lists', 'reading-lists', '특정 주제나 작가별 필독서 리스트'),
('Parenting Hacks', 'parenting-hacks', '연령별 육아 꿀팁과 정보 공유'),
('Creative Writing', 'creative-writing', '공동 세계관, 캐릭터 설정, 짧은 소설 창작'),
('Event Planning', 'event-planning', '결혼식, 파티, 워크샵을 위한 체크리스트와 아이디어'),
('Study Notes', 'study-notes', '시험, 자격증, 교양 과목 핵심 요약 노트'),
('Movie Reviews', 'movie-reviews', '깊이 있는 영화 해석 및 추천 리스트'),
('Personal Branding', 'personal-branding', '나만의 포트폴리오, 이력서 템플릿 만들기')
ON CONFLICT (slug) DO NOTHING;

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_weaves_user_id ON weaves(user_id);
CREATE INDEX IF NOT EXISTS idx_weaves_channel_id ON weaves(channel_id);
CREATE INDEX IF NOT EXISTS idx_weaves_created_at ON weaves(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_weaves_published ON weaves(is_published) WHERE is_published = TRUE;
CREATE INDEX IF NOT EXISTS idx_weaves_featured ON weaves(is_featured) WHERE is_featured = TRUE;
CREATE INDEX IF NOT EXISTS idx_weave_events_weave_id ON analytics.weave_events(weave_id);
CREATE INDEX IF NOT EXISTS idx_weave_events_user_id ON analytics.weave_events(user_id);
CREATE INDEX IF NOT EXISTS idx_weave_events_type ON analytics.weave_events(event_type);