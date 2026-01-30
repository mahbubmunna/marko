'use client';

import { Search } from 'lucide-react';

interface SearchBarProps {
    onSearch: (query: string) => void;
    isSearching: boolean;
}

export default function SearchBar({ onSearch, isSearching }: SearchBarProps) {
    return (
        <div className="relative">
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-stone-400">
                <Search size={14} />
            </div>
            <input
                type="text"
                placeholder="Search notes..."
                onChange={(e) => onSearch(e.target.value)}
                className="w-full pl-9 pr-3 py-1.5 bg-white border border-stone-200 rounded-md text-sm text-stone-700 placeholder-stone-400 focus:outline-none focus:border-stone-400 focus:ring-1 focus:ring-stone-400 transition-colors"
            />
            {isSearching && (
                <div className="absolute inset-y-0 right-0 pr-3 flex items-center">
                    <div className="w-3 h-3 border-2 border-stone-300 border-t-stone-500 rounded-full animate-spin"></div>
                </div>
            )}
        </div>
    );
}
