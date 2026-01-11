# Ant Tower - 시스템 모니터링 대시보드

시스템 프로세스를 실시간으로 모니터링하고 기록하는 프로그램입니다.

## 📋 프로젝트 개요

- **백엔드**: Go (Gin 웹프레임워크)
- **프론트엔드**: React + TypeScript
- **실시간 통신**: WebSocket
- **데이터베이스**: SQLite (GORM ORM)
- **시스템 정보**: gopsutil

---

## 📁 프로젝트 구조

```
ant-tower/
├── main.go                 # 애플리케이션 진입점
├── go.mod / go.sum        # Go 의존성
├── server/                # HTTP API 및 WebSocket 서버
│   ├── router.go          # 라우터 설정 및 API 엔드포인트
│   ├── handler.go         # 핸들러 함수
│   ├── ws_hub.go          # WebSocket 연결 관리
│   └── response.go        # 응답 포맷
├── sysagent/              # 시스템 모니터링 로직
│   ├── monitor.go         # 프로세스 모니터링 서비스
│   ├── console.go         # 콘솔 제어 (Windows)
│   └── windowapi.go       # 윈도우 제어
├── data/                  # 데이터베이스 관련
│   └── db.go              # GORM 설정 및 데이터 모델
├── config/                # 설정 관리
│   └── config.go          # 환경설정 로드
├── logger/                # 로깅 설정
│   └── logger.go          # zerolog 설정
└── dashboard/             # React 프론트엔드
    ├── src/               # 소스코드
    ├── public/            # 정적 파일
    ├── package.json       # Node.js 의존성
    └── vite.config.ts     # Vite 설정
```

---

## 🔧 백엔드 (Go)

### 메인 흐름 (`main.go`)

1. 설정 로드
2. 로거 초기화
3. SQLite 데이터베이스 초기화
4. MonitorService 시작 (백그라운드 프로세스 모니터링)
5. WebSocket 허브 시작
6. HTTP 서버 시작
7. Graceful shutdown 처리

### API 엔드포인트 (`server/router.go`)

| 메서드 | 경로                | 설명                           |
| ------ | ------------------- | ------------------------------ |
| GET    | `/api/uptime`       | 시스템 업타임 조회             |
| GET    | `/api/login`        | 현재 로그인 사용자 정보        |
| GET    | `/api/process`      | 실행 중인 프로세스 목록        |
| POST   | `/api/process/kill` | 프로세스 강제 종료             |
| POST   | `/api/power`        | 시스템 전원 제어 (종료/재부팅) |
| GET    | `/ws`               | WebSocket 연결                 |

### WebSocket 허브 (`server/ws_hub.go`)

- 클라이언트 연결 관리
- 실시간 메시지 브로드캐스트
- 자동 재연결 처리
- 연결 해제 시 정리

### 시스템 모니터링 (`sysagent/monitor.go`)

**MonitorService**:

- 1초 주기로 실행 중인 프로세스 스캔
- 새로운 프로세스 시작 감지 → 브로드캐스트
- 프로세스 종료 감지 → 데이터베이스 저장 → 브로드캐스트

**브로드캐스트 메시지 형식**:

```json
// 프로세스 시작
{
  "type": "START",
  "data": {
    "pid": 1234,
    "name": "process.exe",
    "start_time": "2026-01-10T12:30:45Z",
    "uptime":
  }
}

// 프로세스 종료
{
  "type": "END",
  "data": {
    "pid": 1234,
    "name": "process.exe",
    "start_time": "2026-01-10T12:30:45Z",
    "end_time": "2026-01-10T12:35:50Z",
    "duration": "5m5s",
    "uptime":
  }
}
```

### 데이터베이스 (`data/db.go`)

**ProcessLog 테이블**:
| 컬럼 | 타입 | 설명 |
|-----|------|------|
| PID | int32 | 프로세스 ID |
| Name | string | 프로세스 이름 |
| StartTime | time.Time | 시작 시간 |
| EndTime | \*time.Time | 종료 시간 (NULL이면 실행 중) |
| Duration | string | 실행 지속 시간 |

---

## ⚛️ 프론트엔드 (React)

### 기술 스택

- **React 19** + TypeScript
- **Vite** - 빠른 빌드 도구
- **Ant Design** - UI 컴포넌트 라이브러리
- **ag-Grid** - 데이터 테이블 (프로세스 목록)
- **Axios** - HTTP 클라이언트
- **date-fns** - 날짜 포맷팅

### 주요 의존성

```json
{
  "@ant-design/icons": "^6.1.0",
  "ag-grid-react": "^35.0.0",
  "antd": "^6.1.4",
  "axios": "^1.13.2",
  "date-fns": "^4.1.0",
  "react": "^19.2.0",
  "react-dom": "^19.2.0"
}
```

### 빌드 및 실행

```bash
cd dashboard

# 개발 서버 시작
npm run dev

# 프로덕션 빌드
npm run build

# 프리뷰
npm run preview

# 린트
npm run lint
```

---

## 🔄 데이터 흐름

```
┌─────────────────────────────────────────┐
│   MonitorService (1초 주기)              │
│   - 시스템 프로세스 스캔                  │
│   - 시작/종료 이벤트 감지                │
└────────────┬────────────────────────────┘
             │
             │ broadcast 채널
             ▼
┌─────────────────────────────────────────┐
│   WebSocket Hub                         │
│   - 모든 연결된 클라이언트에 전송        │
└────────────┬────────────────────────────┘
             │
             │ WebSocket
             ▼
┌─────────────────────────────────────────┐
│   React 대시보드                        │
│   - 실시간 프로세스 업데이트             │
│   - 프로세스 목록 표시                   │
│   - API를 통한 제어                      │
└─────────────────────────────────────────┘
             │
             │ API 요청
             ▼
┌─────────────────────────────────────────┐
│   HTTP API Server (Gin)                 │
│   - 프로세스 제어                        │
│   - 시스템 정보 조회                     │
└─────────────────────────────────────────┘
```

---

## 📦 주요 의존성

### Go

| 패키지                          | 용도                  |
| ------------------------------- | --------------------- |
| `github.com/gin-gonic/gin`      | HTTP 웹프레임워크     |
| `github.com/gorilla/websocket`  | WebSocket 통신        |
| `github.com/shirou/gopsutil/v4` | OS/프로세스 정보 조회 |
| `gorm.io/gorm`                  | ORM (데이터베이스)    |
| `gorm.io/driver/sqlite`         | SQLite 드라이버       |
| `github.com/rs/zerolog`         | 구조화된 로깅         |
| `golang.org/x/sys`              | 시스템 호출           |

---

## 🚀 실행 방법

### 사전 요구사항

- Go 1.25.5 이상 (mingw-w64-gcc)
- Node.js 및 npm (프론트엔드 빌드)

### 1. 프론트엔드 빌드

```bash
cd dashboard
npm install
npm run build
```

### 2. 백엔드 실행

```bash
go run main.go
```

**옵션**:

- `--port=5001` - 서버 포트 (기본값: 5001)
- `--console` - 콘솔 창 표시 (기본값: false)
- `TWN_PORT` - 환경변수로 포트 설정 가능

### 3. 웹 브라우저 접속

```
http://localhost:5001
```

---

## 🔐 보안 주의사항

- WebSocket CORS: `CheckOrigin` 함수에서 모든 출처 허용 (필요시 제한 권장)
- API 엔드포인트: 현재 인증 없음 (프로덕션 배포 시 추가 필요)
- 프로세스 제어: 관리자 권한 필요

---

## 📝 라이선스

[라이선스 정보 추가]

---

## 🤝 기여

[기여 가이드 추가]
