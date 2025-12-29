import "./globals.css";
import { ToastProvider } from '../context/ToastContext';
import { AuthProvider } from '../context/AuthContext';
import { PostsProvider } from '../context/PostsContext';
import { ChatProvider } from '../context/ChatContext';

export const metadata = {
  title: "Social Network",
  description: "A social media web app built by zone01 talents for zone01 talents",
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>
        <ToastProvider>
          <AuthProvider>
            <PostsProvider>
              <ChatProvider>
                {children}
              </ChatProvider>
            </PostsProvider>
          </AuthProvider>
        </ToastProvider>
      </body>
    </html>
  );
}
