import './globals.css';
import type { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'FlagForge Dashboard',
  description: 'Manage feature flags and remote config'
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body className="min-h-screen bg-slate-100 text-slate-900">{children}</body>
    </html>
  );
}
