import "./globals.css";
import { AuthProvider } from '../context/AuthContext';
import { ToastProvider } from '../context/ToastContext';
import { ChatProvider } from '../context/ChatContext';
import { PostsProvider } from '../context/PostsContext';

export const metadata = {
  title: "Social UI Update",
  description: "A social media user interface built with Next.js",
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
