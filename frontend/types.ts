export interface Note {
  id: string;
  title: string;
  content?: string;
  createdAt: string;
  updatedAt: string;
}

export interface NoteMetadata {
  title: string;
  tags?: string[];
  created?: string;
  updated?: string;
}
