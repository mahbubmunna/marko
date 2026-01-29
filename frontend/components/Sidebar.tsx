'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { FileText, Folder, Plus } from 'lucide-react';
import { Note } from '../types';
import clsx from 'clsx';

interface SidebarProps {
    notes: Note[];
}

export default function Sidebar({ notes }: SidebarProps) {
    const pathname = usePathname();

    // Guard against undefined notes
    const safeNotes = Array.isArray(notes) ? notes : [];

    // Simple grouping by folder logic could go here.
    // For now, flat list sorted by UpdatedAt
    const sortedNotes = [...safeNotes].sort((a, b) =>
        new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime()
    );

    return (
        <div className="w-64 h-full border-r border-stone-200 bg-stone-50 flex flex-col">
            <div className="p-4 border-b border-stone-200 flex items-center justify-between">
                <h1 className="font-semibold text-stone-700">Dev Notes</h1>
                <Link
                    href="/new"
                    className="p-1.5 hover:bg-stone-200 rounded-md text-stone-600 transition-colors"
                    title="New Note"
                >
                    <Plus size={18} />
                </Link>
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
                                    'block px-3 py-2 rounded-md text-sm transition-colors flex items-center gap-2',
                                    isActive
                                        ? 'bg-white text-stone-900 shadow-sm border border-stone-200'
                                        : 'text-stone-600 hover:bg-stone-200/50 hover:text-stone-900'
                                )}
                            >
                                <FileText size={14} className={clsx(isActive ? "text-stone-800" : "text-stone-400")} />
                                <span className="truncate">{title}</span>
                            </Link>
                        );
                    })}

                    {sortedNotes.length === 0 && (
                        <div className="text-center py-8 text-stone-400 text-sm">
                            No notes found.
                        </div>
                    )}
                </nav>
            </div>
        </div>
    );
}
