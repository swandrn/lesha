import React from "react";

interface Server {
  id: number;
  name: string;
  icon: string;
}

const sampleServers: Server[] = [
  { id: 1, name: "General", icon: "🔥" },
  { id: 2, name: "Gaming", icon: "🎮" },
  { id: 3, name: "Programming", icon: "💻" },
  { id: 4, name: "Music", icon: "🎵" },
  { id: 5, name: "Movies", icon: "🎬" },
];

interface ServerListProps {
  onServerSelect: (serverId: number) => void;
}

export function ServerList(
  { onServerSelect }: ServerListProps,
): React.JSX.Element {
  return (
    <div className="w-20 h-screen bg-blue-900 text-white flex flex-col items-center p-2 shadow-lg">
      {sampleServers.map((server) => (
        <div
          key={server.id}
          className="w-14 h-14 flex items-center justify-center text-2xl font-bold bg-blue-500 hover:bg-blue-400 transition rounded-full cursor-pointer my-3"
          title={server.name}
          onClick={() => onServerSelect(server.id)}
        >
          {server.icon}
        </div>
      ))}
    </div>
  );
}
