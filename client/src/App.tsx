import { useState } from "react";
import { ServerList } from "./components/ServerList";
import { ChannelList } from "./components/ChannelList";
import { Chat } from "./components/Chat";
import EditAccount from "./components/EditAccount";
import { FriendList } from "./components/FriendList";
import CreateServer from "./components/CreateServer";

function App() {
  const [selectedServer, setSelectedServer] = useState<number | null>(null);
  const [selectedChannel, setSelectedChannel] = useState<number | null>(null);
  const [isChannelListVisible, setIsChannelListVisible] = useState(true);
  const [isCreatingServer, setIsCreatingServer] = useState(false); // État pour la création de serveur
  const [servers, setServers] = useState([
    { id: 0, name: "Edit Account", icon: "🛠️" },
    { id: 1, name: "Friends", icon: "👫" },
    { id: 3, name: "Programming", icon: "💻" },
    { id: 4, name: "Music", icon: "🎵" },
    { id: 5, name: "Movies", icon: "🎬" },
  ]);

  // Fonction pour afficher CreateServer lorsque le bouton dédié est cliqué
  const handleCreateNewServer = () => {
    setIsCreatingServer(true); // On active le mode création
    setSelectedServer(null); // On réinitialise la sélection du serveur
    setSelectedChannel(null); // On réinitialise la sélection du canal
  };

  // Fonction pour sélectionner un serveur et afficher ses channels
  const handleServerSelect = (serverId: number) => {
    setSelectedServer(serverId);
    setSelectedChannel(null); // Reset selected channel
    setIsChannelListVisible(true); // Ensure channel list is visible
    setIsCreatingServer(false); // On désactive le mode création quand un serveur est sélectionné
  };

  // Fonction pour sélectionner un canal
  const handleChannelSelect = (channelId: number) => {
    setSelectedChannel(channelId);
    setIsChannelListVisible(false); // Hide channel list when a channel is selected
  };

  return (
    <div className="flex h-screen w-screen bg-gray-900">
      <ServerList
        onServerSelect={handleServerSelect}
        onCreateNewServer={handleCreateNewServer}
      />

      {/* Affichage de EditAccount uniquement pour le serveur 0 */}
      {selectedServer === 0 && <EditAccount />}

      {/* Affichage de FriendList uniquement pour le serveur 1 */}
      {selectedServer === 1 && <FriendList />}

      {/* Affichage du composant CreateServer uniquement lorsque le mode création est activé */}
      {isCreatingServer && <CreateServer />}

      {/* Si un serveur est sélectionné et qu'il n'est pas en mode création, afficher les channels */}
      {selectedServer !== null && !isCreatingServer && selectedServer !== 0 && selectedServer !== 1 && isChannelListVisible && (
        <ChannelList
          serverId={selectedServer}
          onChannelSelect={handleChannelSelect}
        />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />

      {/* Si un canal est sélectionné, afficher le chat */}
      {selectedChannel && (
        <Chat
          channelId={selectedChannel}
          onToggleChannels={() => setIsChannelListVisible(true)}
        />
      )}
    </div>
  );
}

export default App;
