# Frontend (Next.js)

- Next.js 16 app router
- Talks to API Gateway at NEXT_PUBLIC_API_URL (default http://localhost:8080)
- Uses cookie-based auth; fetches with credentials: include

## Setup

- Install and run
  - npm install
  - npm run dev
- Environment
  - NEXT_PUBLIC_API_URL=http://localhost:8080
  - See [api.js](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/frontend/lib/api.js) for usage

## Key Features

- Chat
  - Real-time via WS /chat/ws
  - Auto-scroll on new messages
  - Pagination: 10-by-10 chunks on history
  - Mobile responsive: sidebar/message view toggle
  - Implementation: [page.js](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/frontend/app/(main)/chat/page.js)
- Groups
  - List groups and membership status
  - Create posts with image upload
- Posts
  - Create text/image posts (multipart)
  - Reactions: like/dislike
  - Comments with optional image
- Notifications
  - WebSocket updates for invites and events

## API Usage Helper

- apiFetch(endpoint, options)
  - Adds base URL, JSON headers, credentials
  - See [api.js](file:///Users/macadmin/Desktop/my-area/repos-to-rust/rox-network/frontend/lib/api.js)
- toMediaUrl(path)
  - Builds gateway URL for media paths

## Folder Structure

- app/(main)
  - chat, groups, profile, settings
- app/components
  - UI, posts, notifications, groups, sidebar, emoji
- context
  - AuthContext, ChatContext, PostsContext, ToastContext
- lib
  - api.js base fetch and media URL

## Dev Notes

- Auth
  - Session cookie set by auth-service via gateway
  - apiFetch dispatches auth:unauthorized on 401
- Media
  - Protected via media-service; use toMediaUrl for images
- Ports
  - Frontend 3000; gateway 8080; services 8081–8087

