import { notFound } from 'next/navigation';
import { fetchNote } from '@/lib/api';
import Editor from '@/components/Editor';

interface PageProps {
    params: Promise<{ id: string }>;
}

export default async function NotePage({ params }: PageProps) {
    const { id } = await params;
    // Decode the ID since it might be URL encoded
    const decodedId = decodeURIComponent(id);

    try {
        const note = await fetchNote(decodedId);
        return <Editor note={note} />;
    } catch (e) {
        notFound();
    }
}
