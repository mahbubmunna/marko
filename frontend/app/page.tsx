export default function Home() {
  return (
    <div className="flex flex-col items-center justify-center h-full text-stone-400">
      <div className="text-center space-y-4">
        <h2 className="text-2xl font-medium text-stone-600">No note selected</h2>
        <p>Select a note from the sidebar or create a new one.</p>
      </div>
    </div>
  );
}
