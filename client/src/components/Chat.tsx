import React, { useEffect, useRef, useState } from "react";
import { useUser } from "../hooks/useUser";

interface Message {
  ID?: number;
  Content: string;
  UserID: number;
  User: {
    id: number;
    name: string;
    displayName: string;
  };
  medias: {
    id: number;
    type: string;
    extension: string;
    url: string;
  }[];
  reactions?: {
    id: number;
    emoji: string;
    userId: number;
  }[];
  channelId: number;
  createdAt: string;
}

interface ChatProps {
  channelId: number;
  onToggleChannels: () => void;
}

export function Chat({ channelId, onToggleChannels }: ChatProps): React.JSX.Element {
  const { user } = useUser();
  const [showReactions, setShowReactions] = useState(0);
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState("");
  const [file, setFile] = useState<File | null>(null);
  const messagesEndRef = useRef<HTMLDivElement | null>(null);
  const socketRef = useRef<WebSocket | null>(null);
  const availableReactions = ["👍", "👎", "👏", "🤣", "👀","👌"];

  // Load message history
  useEffect(() => {
    const controller = new AbortController();
    const signal = controller.signal;

    const fetchMessages = async () => {
      try {
        const res = await fetch(`http://localhost:8080/channels/${channelId}/messages`, {
          credentials: "include",
          signal,
        });

        if (!res.ok) throw new Error("Failed to fetch messages");

        const data = await res.json();
        console.log("received data", data);

        setMessages(
          data.map((msg: any) => ({
            ID: msg.id,
            Content: msg.content,
            UserID: msg.user.id,
            User: msg.user,
            medias: msg.medias,
            reactions: msg.reactions,
            channelId: msg.channelId,
            createdAt: msg.createdAt,
          }))
        );
      } catch (err: any) {
        if (err.name === "AbortError") {
          console.log("Message fetch aborted");
          return;
        }
        console.error("Error fetching messages:", err);
      }
    };

    fetchMessages();

    return () => {
      controller.abort();
    };
  }, [channelId]);

  // Scroll to bottom
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
    console.log(messages);
  }, [messages]);

  // Connect to WebSocket
  useEffect(() => {
    const ws = new WebSocket("ws://localhost:8080/ws");
    socketRef.current = ws;

    let isActive = true;

    ws.onopen = () => {
      if (!isActive) return;
      console.log("WebSocket connected");

      // Join this channel explicitly
      ws.send(
        JSON.stringify({
          type: "JOIN_CHANNEL",
          channel_id: channelId,
        })
      );
    };

    ws.onmessage = (event) => {
      if (!isActive) return;
      try {
        const data = JSON.parse(event.data);
        if (data.type === "MESSAGE" && data.channel_id === channelId) {
          const newMsg: Message = {
            ID: data.id,
            User: data.user,
            medias: data.medias,
            reactions: data.reactions,
            Content: data.content,
            channelId: data.channel_id,
            UserID: data.sender,
            createdAt: data.timestamp,
          };
          setMessages((prev) => [...prev, newMsg]);
        }
      } catch (err) {
        console.error("Invalid WebSocket message:", err);
      }
    };

    ws.onerror = (err) => {
      if (!isActive) return;
      console.error("WebSocket error:", err);
    };

    ws.onclose = () => {
      if (!isActive) return;
      console.log("WebSocket closed");
    };

    return () => {
      isActive = false;
      ws.close();
    };
  }, [channelId]);

  // Send message via WebSocket only
  const sendMessage = (e?: React.FormEvent) => {
    e?.preventDefault();
    if (!input.trim() && !file) return;
    if (file && !input.trim()) {
      setInput(file.name);
    }

    if (socketRef.current?.readyState === WebSocket.OPEN) {
      socketRef.current.send(
        JSON.stringify({
          type: "MESSAGE",
          channel_id: channelId,
          content: input,
          file: file,
        })
      );
      setInput("");
      setFile(null);
    } else {
      console.warn("WebSocket not ready");
    }
  };

  const handleReactionClick = (reaction: string, messageId: number) => {
    
    if (socketRef.current?.readyState === WebSocket.OPEN) {
      socketRef.current.send(
        JSON.stringify({
          type: "REACTION",
          message_id: messageId,
          reaction: reaction,
        })
      );
    } else {
      console.warn("WebSocket not ready");
    }
  };

  return (
    <div className="flex flex-col h-screen w-full bg-gray-100">
      {/* Header */}
      <div className="p-4 bg-blue-600 text-white text-lg font-semibold shadow-md flex justify-between items-center">
        <button onClick={onToggleChannels} className="bg-blue-400 px-3 py-1 rounded-md hover:bg-blue-500 transition">
          ☰ Channels
        </button>
        <span>Channel {channelId}</span>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto p-4 space-y-2">
        {messages.map((msg, i) => (
          <div
            key={i}
            className={`px-4 py-2 rounded-lg w-fit max-w-[75%] break-words ${msg.UserID === user?.user.id
              ? "bg-blue-500 text-white self-end ml-auto"
              : "bg-gray-300 text-black self-start mr-auto"
              }`}
            onMouseEnter={() => setShowReactions(msg.ID ?? 0)}
            onMouseLeave={() => setShowReactions(0)}
          >
            {msg.medias?.length > 0 ? (
              <div className="flex items-center gap-2">
                {msg.medias.map((media, index) => (
                  <img
                    key={index}
                    src={"http://localhost:8080/" + media.url}
                    alt={media.type}
                    className="w-10 h-10 object-cover rounded-md"
                  />
                ))}
              </div>
            ) : (
              <p>{msg.Content}</p>
            )}
           
            <div className="text-xs text-right mt-1 text-gray-500">
              {new Date(msg.createdAt).toLocaleTimeString()}
            </div>
            {msg.reactions && msg.reactions.length > 0 && (
              <div className="text-md text-right mt-1 text-gray-500">
                {msg.reactions.map((reaction) => (
                  <span key={reaction.emoji}>{reaction.emoji}</span>
                ))}
              </div>
            )}

            {showReactions === msg.ID && (
              <div className="flex items-center gap-2">
              {availableReactions.map((reaction) => (
                <button key={reaction} onClick={() => handleReactionClick(reaction, msg.ID ?? 0)} className="text-lg text-gray-500 hover:text-gray-700 cursor-pointer hover:scale-110 transition-all duration-300">
                  {reaction}
                </button>
                ))}
              </div>
            )}
          </div>
        ))}
        <div ref={messagesEndRef} />
      </div>

      {/* Input */}
      <form
        onSubmit={sendMessage}
        className="flex p-4 bg-white border-t border-gray-300"
      >
        <input
          type="file"
          id="fileInput"
          style={{ display: "none" }}
          onChange={(e) => setFile(e.target.files?.[0] || null)}
          accept="image/*,video/*,audio/*,.pdf,.doc,.docx,.txt"
        />
        <button
          type="button"
          onClick={() => document.getElementById("fileInput")?.click()}
          className="ml-2 px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 transition"
        >
          Attach
        </button>

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