import React, { useEffect, useRef, useState } from "react";

interface Message {
  id: number;
  text: string;
  sender: "user" | "bot";
}

interface ChatProps {
  channelId: number;
  onToggleChannels: () => void;
}

export function Chat(
  { channelId, onToggleChannels }: ChatProps,
): React.JSX.Element {
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState("");
  const messagesEndRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    setMessages([]); // Clear chat when switching channels
  }, [channelId]);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  const sendMessage = (event?: React.FormEvent) => {
    event?.preventDefault();
    if (!input.trim()) return;

    const newMessage: Message = {
      id: messages.length + 1,
      text: input,
      sender: "user",
    };

    setMessages([...messages, newMessage]);
    setInput("");

    setTimeout(() => {
      const botMessage: Message = {
        id: messages.length + 2,
        text: "Bot reply in channel " + channelId,
        sender: "bot",
      };
      setMessages((prevMessages) => [...prevMessages, botMessage]);
    }, 1000);
  };

  return (
    <div className="flex flex-col h-screen w-full bg-gray-100">
      {/* Chat Header */}
      <div className="p-4 bg-blue-600 text-white text-lg font-semibold shadow-md flex justify-between items-center">
        <button
          onClick={onToggleChannels}
          className="bg-blue-400 px-3 py-1 rounded-md hover:bg-blue-500 transition"
        >
          ☰ Channels
        </button>
        <span>Channel {channelId}</span>
      </div>

      {/* Chat Messages */}
      <div className="flex-1 overflow-y-auto p-4 space-y-2">
        {messages.map((msg) => (
          <div
            key={msg.id}
            className={`px-4 py-2 rounded-lg w-fit max-w-[75%] break-words ${
              msg.sender === "user"
                ? "bg-blue-500 text-white self-end ml-auto"
                : "bg-gray-300 text-black self-start mr-auto"
            }`}
          >
            {msg.text}
          </div>
        ))}
        <div ref={messagesEndRef} />
      </div>

      {/* Chat Input */}
      <form
        onSubmit={sendMessage}
        className="flex p-4 bg-white border-t border-gray-300"
      >
        <input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="Type a message..."
        />
        <button
          type="submit"
          className="ml-2 px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 transition"
        >
          Send
        </button>
      </form>
    </div>
  );
}
