# Weave - The Living Document Network

Weave는 아이디어가 진화하는 과정 자체가 콘텐츠가 되는 차세대 소셜 네트워크입니다.

## 서비스 구성

- **weave-be**: 메인 백엔드 API 서버 (Go + Gin + GORM)
- **weave-fe**: 프론트엔드 웹 애플리케이션 (React)
- **weave-module**: MSA 공통 모듈 (Go)
- **weave-scheduler**: 스케줄러 서버 (알림, 주기적 작업)
- **weave-worker**: 워커 서버 (백그라운드 작업 처리)

## 기술 스택

- **Backend**: Go, Gin, GORM
- **Frontend**: React
- **Database**: PostgreSQL, TimescaleDB
- **Cache**: Redis
- **Queue**: RabbitMQ
- **Infrastructure**: Docker, Docker Compose

## 개발 환경 설정

### 로컬 개발
```bash
# 인프라 서비스 시작 (PostgreSQL, Redis 등)
make dev-infra

# 각 서버를 개별적으로 실행
make dev-be      # 백엔드 서버
make dev-fe      # 프론트엔드 서버
make dev-scheduler  # 스케줄러 서버
make dev-worker  # 워커 서버
```

### 배포 환경
```bash
# 전체 서비스 배포
make deploy
```

## 핵심 기능

- **Weave**: 버전 관리가 내장된 살아있는 포스트
- **Timeline**: 아이디어 진화 과정의 타임랩스 뷰
- **Channels**: 주제별 커뮤니티 (w/recipes, w/travel-plans 등)
- **Lab**: 협업을 위한 토론 공간
- **Profile**: 과정이 포트폴리오가 되는 개인 페이지