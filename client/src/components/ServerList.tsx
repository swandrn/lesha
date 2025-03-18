import React, { useEffect } from "react";
import axios from "axios";

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

export function ServerList(): React.JSX.Element {

  useEffect(() => {
    axios.get("http://localhost:8080/logout").then((res) => {
      console.log(res.data);
    });
  }, []);
  return (
    <div className="w-20 md:w-24 h-screen bg-blue-700 text-white flex flex-col items-center p-2 shadow-lg">
      {sampleServers.map((server) => (
        <div
          key={server.id}
          className="w-12 h-12 md:w-16 md:h-16 flex items-center justify-center text-2xl font-bold bg-blue-500 hover:bg-blue-400 transition rounded-full cursor-pointer my-2"
          title={server.name}
        >
          {server.icon}
        </div>
      ))}
    </div>
  );
}
