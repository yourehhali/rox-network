import "./globals.css";
import { AuthProvider } from '../context/AuthContext';
import { ToastProvider } from '../context/ToastContext';
import { ChatProvider } from '../context/ChatContext';

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
            <ChatProvider>
              {children}
            </ChatProvider>
          </AuthProvider>
        </ToastProvider>
      </body>
    </html>
  );
}
