import React, { useEffect } from "react";
import axios from "axios";

interface Server {
  id: number;
  name: string;
  icon: string;
}

const sampleServers: Server[] = [
  { id: 0, name: "Edit Account", icon: "🛠️" },
  { id: 1, name: "Friends", icon: "👫" },
  { id: 3, name: "Programming", icon: "💻" },
  { id: 4, name: "Music", icon: "🎵" },
  { id: 5, name: "Movies", icon: "🎬" },
];

interface ServerListProps {
  onServerSelect: (serverId: number) => void;
  onCreateNewServer: () => void; // Ajoutez une fonction pour gérer la création d'un serveur
}

export function ServerList(
  { onServerSelect, onCreateNewServer }: ServerListProps,
): React.JSX.Element {
  return (
    <div className="w-20 h-screen bg-blue-900 text-white flex flex-col items-center p-2 shadow-lg">
      {sampleServers.map((server, index) => (
        <React.Fragment key={server.id}>
          <div
            className="w-14 h-14 flex items-center justify-center text-2xl font-bold bg-blue-500 hover:bg-blue-400 transition rounded-full cursor-pointer my-3"
            title={server.name}
            onClick={() => onServerSelect(server.id)}
          >
            {server.icon}
          </div>

          {/* Ajouter une barre sous le 2ème serveur (index 1) */}
          {index === 1 && <div className="w-12 h-1 bg-gray-500 my-3" />}
        </React.Fragment>
      ))}

      {/* Ajouter une option pour créer un nouveau serveur */}
      <div
        className="w-14 h-14 flex items-center justify-center text-2xl font-bold bg-green-500 hover:bg-green-400 transition rounded-full cursor-pointer my-3"
        title="Create New Server"
        onClick={onCreateNewServer}
      >
        ➕
      </div>
    </div>
  );
}
