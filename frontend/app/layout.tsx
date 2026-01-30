import { fetchNotes } from '@/lib/api';
import { Note } from '@/types';
import Sidebar from '@/components/Sidebar';
import './globals.css';
import { Inter } from 'next/font/google';

const inter = Inter({ subsets: ['latin'] });

export const metadata = {
  title: 'Dev Notes Vault',
  description: 'Local-first markdown notes',
};

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  // Fetch notes on the server
  // In a real app we might handle error state better
  let notes: Note[] = [];
  try {
    notes = await fetchNotes();
  } catch (e) {
    console.error("Failed to fetch notes:", e);
  }

  return (
    <html lang="en" suppressHydrationWarning>
      <body className={`${inter.className} text-stone-800 bg-white h-screen flex overflow-hidden`} suppressHydrationWarning>
        <aside className="h-full flex-shrink-0">
          <Sidebar notes={notes} />
        </aside>
        <main className="flex-1 h-full overflow-hidden flex flex-col relative">
          {children}
        </main>
      </body>
    </html>
  );
}
