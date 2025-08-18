'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { createNote } from '@/lib/api';

export default function NewNotePage() {
    const router = useRouter();

    useEffect(() => {
        // Create a generic new note and redirect
        const init = async () => {
            try {
                const { id } = await createNote('# New Note\n\nStarting writing...');
                router.refresh(); // Refresh sidebar
                router.push(`/note/${id}`);
            } catch (e) {
                console.error('Failed to create note', e);
            }
        };
        init();
    }, [router]);

    return (
        <div className="flex items-center justify-center h-full text-stone-400">
            Creating note...
        </div>
    );
}
