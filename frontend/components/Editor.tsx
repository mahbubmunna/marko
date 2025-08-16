'use client';

import { useState, useEffect, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import ReactMarkdown from 'react-markdown';
import { Eye, Edit2, Save, Trash2 } from 'lucide-react';
import clsx from 'clsx';
import { updateNote, deleteNote } from '@/lib/api';
import { Note } from '@/types';

interface EditorProps {
    note: Note;
}

export default function Editor({ note }: EditorProps) {
    const router = useRouter();
    const [content, setContent] = useState(note.content || '');
    const [isPreview, setIsPreview] = useState(true);
    const [isSaving, setIsSaving] = useState(false);
    const [lastSaved, setLastSaved] = useState(note.updatedAt);

    // Debounced save
    useEffect(() => {
        if (content === (note.content || '')) return;

        const handler = setTimeout(async () => {
            setIsSaving(true);
            try {
                await updateNote(note.id, content);
                setLastSaved(new Date().toISOString());
                router.refresh(); // Refresh sidebar for title updates
            } catch (e) {
                console.error('Autosave failed', e);
            } finally {
                setIsSaving(false);
            }
        }, 1000);

        return () => clearTimeout(handler);
    }, [content, note.id, note.content, router]);

    const handleDelete = async () => {
        if (confirm('Are you sure you want to delete this note?')) {
            await deleteNote(note.id);
            router.refresh();
            router.push('/');
        }
    };

    return (
        <div className="flex flex-col h-full bg-white">
            {/* Toolbar */}
            <div className="flex items-center justify-between px-6 py-3 border-b border-stone-100 bg-white/50 backdrop-blur-sm z-10">
                <div className="flex items-center gap-4 text-xs text-stone-400 font-mono">
                    <span>{isSaving ? 'Saving...' : 'Saved'}</span>
                    <span>{new Date(lastSaved).toLocaleTimeString()}</span>
                </div>

                <div className="flex items-center gap-2">
                    <button
                        onClick={() => setIsPreview(!isPreview)}
                        className="p-2 text-stone-500 hover:bg-stone-100 rounded-md transition-colors"
                        title={isPreview ? "Edit" : "Preview"}
                    >
                        {isPreview ? <Edit2 size={18} /> : <Eye size={18} />}
                    </button>

                    <button
                        onClick={handleDelete}
                        className="p-2 text-red-400 hover:bg-red-50 rounded-md transition-colors"
                        title="Delete Note"
                    >
                        <Trash2 size={18} />
                    </button>
                </div>
            </div>

            {/* Content */}
            <div className="flex-1 overflow-y-auto">
                <div className="max-w-3xl mx-auto px-8 py-12 min-h-full">
                    {isPreview ? (
                        <div className="prose prose-stone prose-lg max-w-none">
                            <ReactMarkdown>{content}</ReactMarkdown>
                        </div>
                    ) : (
                        <textarea
                            value={content}
                            onChange={(e) => setContent(e.target.value)}
                            className="w-full h-full min-h-[70vh] resize-none outline-none font-mono text-stone-800 bg-transparent leading-relaxed"
                            placeholder="Start writing..."
                            spellCheck={false}
                            autoFocus
                        />
                    )}
                </div>
            </div>
        </div>
    );
}
