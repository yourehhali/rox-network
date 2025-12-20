import "./globals.css";
import { AuthProvider } from '../context/AuthContext';

export const metadata = {
  title: "Social UI Update",
  description: "A social media user interface built with Next.js",
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>
        <AuthProvider>
          {children}
        </AuthProvider>
      </body>
    </html>
  );
}
