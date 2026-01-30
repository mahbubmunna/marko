'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useState, useEffect } from 'react';
import { FileText, Plus } from 'lucide-react';
import { Note } from '../types';
import { searchNotes } from '../lib/api';
import SearchBar from './SearchBar';
import clsx from 'clsx';

interface SidebarProps {
    notes: Note[];
}

export default function Sidebar({ notes }: SidebarProps) {
    const pathname = usePathname();
    const [searchQuery, setSearchQuery] = useState('');
    const [searchResults, setSearchResults] = useState<Note[] | null>(null);
    const [isSearching, setIsSearching] = useState(false);

    // Guard against undefined notes
    const safeNotes = Array.isArray(notes) ? notes : [];

    // Debounced search
    useEffect(() => {
        if (!searchQuery.trim()) {
            setSearchResults(null);
            return;
        }

        const timer = setTimeout(async () => {
            setIsSearching(true);
            try {
                const results = await searchNotes(searchQuery);
                setSearchResults(results);
            } catch (e) {
                console.error("Search failed", e);
            } finally {
                setIsSearching(false);
            }
        }, 300);

        return () => clearTimeout(timer);
    }, [searchQuery]);

    // Use search results if active, otherwise sorted full list
    const displayNotes = searchResults || safeNotes;

    // Sort: if searching, preserve rank. If not, sort by UpdatedAt
    const sortedNotes = searchResults ? displayNotes : [...displayNotes].sort((a, b) =>
        new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime()
    );

    return (
        <div className="w-64 h-full border-r border-stone-200 bg-stone-50 flex flex-col">
            <div className="p-4 border-b border-stone-200 space-y-3">
                <div className="flex items-center justify-between">
                    <h1 className="font-semibold text-stone-700">Dev Notes</h1>
                    <Link
                        href="/new"
                        className="p-1.5 hover:bg-stone-200 rounded-md text-stone-600 transition-colors"
                        title="New Note"
                    >
                        <Plus size={18} />
                    </Link>
                </div>
                <SearchBar onSearch={setSearchQuery} isSearching={isSearching} />
            </div>

            <div className="flex-1 overflow-y-auto p-2">
                <nav className="space-y-0.5">
                    {sortedNotes.map((note) => {
                        const isActive = pathname === `/note/${note.id}`;
                        const title = note.title || 'Untitled';

                        return (
                            <Link
                                key={note.id}
                                href={`/note/${note.id}`}
                                className={clsx(
                                    'block px-3 py-2 rounded-md text-sm transition-colors flex flex-col gap-1',
                                    isActive
                                        ? 'bg-white text-stone-900 shadow-sm border border-stone-200'
                                        : 'text-stone-600 hover:bg-stone-200/50 hover:text-stone-900'
                                )}
                            >
                                <div className="flex items-center gap-2">
                                    <FileText size={14} className={clsx(isActive ? "text-stone-800" : "text-stone-400")} />
                                    <span className="truncate font-medium">{title}</span>
                                </div>
                                {searchResults && note.content && (
                                    <div
                                        className="text-xs text-stone-400 pl-6 line-clamp-2"
                                        dangerouslySetInnerHTML={{ __html: note.content }}
                                    />
                                )}
                            </Link>
                        );
                    })}

                    {sortedNotes.length === 0 && (
                        <div className="text-center py-8 text-stone-400 text-sm">
                            {searchQuery ? 'No results found.' : 'No notes found.'}
                        </div>
                    )}
                </nav>
            </div>
        </div>
    );
}
