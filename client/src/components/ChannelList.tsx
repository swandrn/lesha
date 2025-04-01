import { useState } from "react";
import AddFriendToServer  from "./AddFriendToServer"; // Assurez-vous d'importer votre composant

interface ChannelListProps {
  serverId: number;
  onChannelSelect: (channelId: number) => void;
}

export function ChannelList({ serverId, onChannelSelect }: ChannelListProps): React.JSX.Element {
  const [isAddFriendVisible, setIsAddFriendVisible] = useState(false); // Ajouter l'état

  // Liste des canaux pour un serveur (exemple fictif)
  const channels = [
    { id: 101, name: "general" },
    { id: 102, name: "random" },
  ];

  const handleInviteClick = () => {
    setIsAddFriendVisible(true); // Afficher le composant AddFriendToServer lorsque le bouton est cliqué
  };

  return (
    <div className="w-40 md:w-56 h-screen bg-gray-800 text-white flex flex-col p-4 shadow-md">
      <h2 className="text-lg font-semibold mb-4">Channels</h2>
      
      {/* Liste des canaux */}
      {channels.map((channel) => (
        <div
          key={channel.id}
          className="p-2 text-sm rounded-md hover:bg-gray-700 cursor-pointer"
          onClick={() => onChannelSelect(channel.id)}
        >
          # {channel.name}
        </div>
      ))}

      {/* Bouton pour inviter un ami */}
      <button
        onClick={handleInviteClick}
        className="mt-4 bg-blue-600 text-white px-4 py-2 rounded-md text-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50"
      >
        Inviter un ami au serveur
      </button>

      {/* Affichage du composant AddFriendToServer lorsque isAddFriendVisible est true */}
      {isAddFriendVisible && <AddFriendToServer serverId={serverId} />}
    </div>
  );
}
