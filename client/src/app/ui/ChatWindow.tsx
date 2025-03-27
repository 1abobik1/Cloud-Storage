import React from "react";

type ChatWindowProps = {
  selectedChat: number | null;
};

const messages = [
  { id: 1, text: "Привет!", sender: "me" },
  { id: 2, text: "Как дела?", sender: "other" },
  { id: 3, text: "Все хорошо, спасибо!", sender: "me" },
];

export default function ChatWindow({ selectedChat }: ChatWindowProps) {
  if (!selectedChat) {
    return <div className="flex-1 flex items-center justify-center">Выберите чат</div>;
  }

  return (
    <div className="flex-1 flex flex-col">
      <div className="p-4 bg-gray-100 border-b font-bold">Чат {selectedChat}</div>
      <div className="flex-1 p-4 overflow-auto">
        {messages.map((msg) => (
          <div
            key={msg.id}
            className={`mb-2 p-2 rounded-lg max-w-xs ${
              msg.sender === "me" ? "bg-blue-500 text-white ml-auto" : "bg-gray-300 text-black"
            }`}
          >
            {msg.text}
          </div>
        ))}
      </div>
    </div>
  );
}
