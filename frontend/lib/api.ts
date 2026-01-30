import { Note } from '../types';

const API_BASE = 'http://localhost:8080/api/notes';

export async function fetchNotes(): Promise<Note[]> {
    const res = await fetch(API_BASE);
    if (!res.ok) throw new Error('Failed to fetch notes');
    return res.json();
}

export async function fetchNote(id: string): Promise<Note> {
    const res = await fetch(`${API_BASE}/${id}`);
    if (!res.ok) throw new Error('Failed to fetch note');
    return res.json();
}

export async function searchNotes(query: string): Promise<Note[]> {
  const res = await fetch(`${API_BASE.replace('/api/notes', '/api/search')}?q=${encodeURIComponent(query)}`);
  if (!res.ok) throw new Error('Failed to search notes');
  return res.json();
}

export async function createNote(content: string): Promise<{ id: string }> {
    const res = await fetch(API_BASE, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ content }),
    });
    if (!res.ok) throw new Error('Failed to create note');
    return res.json();
}

export async function updateNote(id: string, content: string): Promise<void> {
    const res = await fetch(`${API_BASE}/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ content }),
    });
    if (!res.ok) throw new Error('Failed to update note');
}

export async function deleteNote(id: string): Promise<void> {
    const res = await fetch(`${API_BASE}/${id}`, {
        method: 'DELETE',
    });
    if (!res.ok) throw new Error('Failed to delete note');
}
