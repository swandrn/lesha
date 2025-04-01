import React from "react";
import { Server } from "./MainLayout";

interface SidebarProps {
  servers: Server[];
  onServerSelect: (serverId: number) => void;
  onCreateNewServer: () => void;
  onNavigate: (view: "edit" | "friends") => void;
}

export const Sidebar = ({
  servers,
  onServerSelect,
  onCreateNewServer,
  onNavigate,
}: SidebarProps): React.JSX.Element => {
  return (
    <div className="w-20 h-screen bg-blue-900 text-white flex flex-col items-center p-2 shadow-lg">
      {/* Edit Account */}
      <div
        className="w-14 h-14 flex items-center justify-center text-xl font-bold bg-gray-600 hover:bg-gray-500 transition rounded-full cursor-pointer my-3"
        title="Edit Account"
        onClick={() => onNavigate("edit")}
      >
        🛠️
      </div>

      {/* Friends */}
      <div
        className="w-14 h-14 flex items-center justify-center text-xl font-bold bg-gray-600 hover:bg-gray-500 transition rounded-full cursor-pointer my-3"
        title="Friends"
        onClick={() => onNavigate("friends")}
      >
        👫
      </div>

      {/* Separator */}
      <div className="w-12 h-1 bg-gray-500 my-3" />

      {/* Server Buttons */}
      {servers.map((server) => (
        <div
          key={server.id}
          className="w-14 h-14 flex items-center justify-center text-2xl font-bold bg-blue-500 hover:bg-blue-400 transition rounded-full cursor-pointer my-3 overflow-hidden"
          title={server.name}
          onClick={() => onServerSelect(server.id)}
        >
          {server.image ? (
            <img
              src={`http://localhost:8080/${server.image}`}
              alt={server.name}
              className="w-full h-full object-cover rounded-full"
            />
          ) : (
            <span className="text-lg">
              {server.name?.charAt(0).toUpperCase() || "?"}
            </span>
          )}
        </div>
      ))}

      {/* Create New Server button */}
      <div
        className="w-14 h-14 flex items-center justify-center text-2xl font-bold bg-green-500 hover:bg-green-400 transition rounded-full cursor-pointer my-3"
        title="Créer un serveur"
        onClick={onCreateNewServer}
      >
        ➕
      </div>
    </div>
  );
};