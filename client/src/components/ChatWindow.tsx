import { useState } from "react";
import { ChannelList } from "./ChannelList";
import { Chat } from "./Chat";
import { ServerList } from "./ServerList";
import EditAccount from "./EditAccount";
import { FriendList } from "./FriendList";

export function ChatWindow(): React.JSX.Element {
  const [selectedServer, setSelectedServer] = useState<number | null>(null);
  const [selectedChannel, setSelectedChannel] = useState<number | null>(null);
  const [isChannelListVisible, setIsChannelListVisible] = useState(true);
  const [servers, setServers] = useState([
    { id: 0, name: "Edit Account", icon: "🛠️" },
    { id: 1, name: "Friends", icon: "👫" },
    { id: 3, name: "Programming", icon: "💻" },
    { id: 4, name: "Music", icon: "🎵" },
    { id: 5, name: "Movies", icon: "🎬" },
  ]);

  const handleServerSelect = (serverId: number) => {
    setSelectedServer(serverId);
    setSelectedChannel(null); // Reset selected channel
    setIsChannelListVisible(true); // Ensure channel list is visible
  };

  const handleChannelSelect = (channelId: number) => {
    setSelectedChannel(channelId);
    setIsChannelListVisible(false); // Hide channel list when a channel is selected
  };

  // Fonction pour créer un nouveau serveur
  const handleCreateNewServer = () => {
    const newServer = {
      id: servers.length, // Assurez-vous que l'id est unique
      name: `New Server ${servers.length + 1}`,
      icon: "✨", // Vous pouvez ajouter un icône personnalisé
    };
    setServers((prevServers) => [...prevServers, newServer]);
  };

  return (
    <div className="flex h-screen w-screen bg-gray-900">
      <ServerList
        onServerSelect={handleServerSelect}
        onCreateNewServer={handleCreateNewServer}
      />

      {/* Afficher EditAccount si le serveur sélectionné est 0 */}
      {selectedServer === 0 && <EditAccount />}

      {/* Afficher FriendList si le serveur sélectionné est 1 */}
      {selectedServer === 1 && <FriendList />}

      {/* Show ChannelList only if a server is selected AND it's visible */}
      {selectedServer && selectedServer !== 0 && selectedServer !== 1 &&
        isChannelListVisible && (
        <ChannelList
          serverId={selectedServer}
          onChannelSelect={handleChannelSelect}
        />
      )}

      {/* Show Chat only if a channel is selected */}
      {selectedChannel && (
        <Chat
          channelId={selectedChannel}
          onToggleChannels={() => setIsChannelListVisible(true)}
        />
      )}
    </div>
  );
}
