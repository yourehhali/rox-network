# Social Network Frontend

This is the frontend application for the Social Network project, built with [Next.js 16](https://nextjs.org/) and [React 19](https://react.dev/).

## Features

- **Authentication**: Login and Registration pages.
- **Social Feed**: View posts, events, and comments.
- **Groups**: Create and join groups, view group posts.
- **Chat**: Real-time chat functionality.
- **Profile**: User profile management and settings.
- **Notifications**: Real-time notifications for events and interactions.

## Tech Stack

- **Framework**: Next.js 16 (App Router)
- **Library**: React 19
- **Styling**: CSS Modules
- **State Management**: React Context (`AuthContext`, `ChatContext`, `PostsContext`, `ToastContext`)

## Getting Started

### Prerequisites

- Node.js (v18 or later recommended)
- npm, yarn, pnpm, or bun

### Installation

1. Navigate to the frontend directory:
   ```bash
   cd ui-frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   # or
   yarn install
   ```

### Running the Development Server

```bash
npm run dev
# or
yarn dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

## Project Structure

```
ui-frontend/
├── app/                  # App Router pages and layouts
│   ├── (main)/           # Main application routes (authenticated)
│   ├── components/       # Reusable UI components
│   ├── context/          # React Context providers
│   ├── login/            # Login page
│   ├── register/         # Registration page
│   └── lib/              # Utility functions (API clients, etc.)
├── public/               # Static assets
└── ...
```

## Configuration

The frontend is configured to communicate with the backend API Gateway.
By default, it expects the backend to be running on `http://localhost:8080`.

To change the backend URL, check `lib/api.js` or environment configuration (if applicable).
