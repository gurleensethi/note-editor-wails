import "./App.css";
import { main } from "../wailsjs/go/models";
import {
  CreateNote,
  GetAllNotes,
  UpdateNoteTitle,
  DeleteArticle,
  GetNoteByID,
  UpdateNoteText,
} from "../wailsjs/go/main/App";
import { useEffect, useState } from "react";

export const App = () => {
  const [notes, setNotes] = useState<main.Note[]>([]);
  const [selectedNote, setSelectedNote] = useState<main.Note | null>(null);

  const fetchNotes = () => {
    GetAllNotes().then((res) => {
      setNotes(res);
    });
  };

  useEffect(() => {
    fetchNotes();
  }, []);

  const handleOnNoteCreate = () => {
    CreateNote({ title: "New Note" }).then(() => {
      fetchNotes();
    });
  };

  const handleSelectNote = (id: number) => {
    GetNoteByID(id).then((note) => {
      console.log(note);
      setSelectedNote(note);
    });
  };

  const handleTitleUpdate = (title: string) => {
    if (!selectedNote) {
      return;
    }

    setNotes((notes) => {
      return notes.map((item) => {
        if (item.id === selectedNote.id) {
          item.title = title;
        }
        return item;
      });
    });

    setSelectedNote((v) => {
      if (!v) {
        return v;
      }

      const note = new main.Note();
      note.id = v.id;
      note.note = v.note;
      note.createdAt = v.createdAt;
      note.title = title;

      return note;
    });

    UpdateNoteTitle(selectedNote?.id, title).then(() => {});
  };

  const handleNoteUpdate = (note: string) => {
    if (!selectedNote) {
      return;
    }

    UpdateNoteText(selectedNote.id, note).then(() => {});
    setSelectedNote((v) => {
      if (!v) {
        return v;
      }

      return main.Note.createFrom({ note });
    });
  };

  const handleArticleDelete = () => {
    if (!selectedNote) {
      return;
    }

    setSelectedNote(null);
    DeleteArticle(selectedNote.id).then((res) => {
      fetchNotes();
    });
  };

  return (
    <div className="app">
      <div className="sidebar">
        <button className="create-note-btn" onClick={handleOnNoteCreate}>
          + Create Note
        </button>
        <NoteList
          notes={notes}
          onNoteSelected={(note) => handleSelectNote(note.id)}
        />
      </div>
      {selectedNote && (
        <div className="note-editor">
          <input
            className="title"
            placeholder="Enter title here..."
            value={selectedNote.title}
            onChange={(e) => handleTitleUpdate(e.target.value)}
          />
          <textarea
            className="editor"
            placeholder="Start writing here..."
            onChange={(e) => handleNoteUpdate(e.target.value)}
            value={selectedNote.note}
          />
          <div className="bottom">
            <div className="note-date">{selectedNote.createdAt}</div>
            <button className="delete-note-btn" onClick={handleArticleDelete}>
              Delete
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export const NoteList = ({
  notes,
  onNoteSelected,
}: {
  notes: main.Note[];
  onNoteSelected: (note: main.Note) => void;
}) => {
  return (
    <div>
      {notes.map((item) => (
        <div
          key={item.createdAt}
          className="note-list-item"
          onClick={() => onNoteSelected(item)}
        >
          {item.title || "(Empty title)"}
        </div>
      ))}
    </div>
  );
};

export default App;
