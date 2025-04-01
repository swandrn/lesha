import { useState, useEffect } from "react";
import AddFriendToServer from "./AddFriendToServer"; // Import the popup/modal component

interface Channel {
  id: number;
  name: string;
  serverId: number;
  createdAt: string;
  updatedAt: string;
}

interface ChannelListProps {
  serverId: number;
  serverOwnerId: number;
  currentUserId: number;
  onChannelSelect: (channelId: number) => void;
}

export function ChannelList({
  serverId,
  serverOwnerId,
  currentUserId,
  onChannelSelect,
}: ChannelListProps) {
  const [channels, setChannels] = useState<Channel[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [isInviteModalOpen, setIsInviteModalOpen] = useState(false); // State to track modal visibility

  const fetchChannels = () => {
    setLoading(true);
    fetch(`http://localhost:8080/servers/${serverId}/channels`, {
      credentials: "include",
    })
      .then((res) => {
        if (!res.ok) throw new Error("Failed to fetch channels");
        return res.json();
      })
      .then((data) => {
        setChannels(
          data.map((channel: any) => ({
            id: channel.ID,
            name: channel.Name,
            serverId: channel.ServerID,
            createdAt: channel.CreatedAt,
            updatedAt: channel.UpdatedAt,
          }))
        );
        setError("");
      })
      .catch((err) => {
        console.error(err);
        setError("Erreur lors du chargement des channels.");
      })
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    fetchChannels();
  }, [serverId]);

  const handleCreateChannel = async () => {
    const name = prompt("Nom du nouveau canal:");
    if (!name) return;

    const formData = new FormData();
    formData.append("name", name);
    formData.append("serverID", serverId.toString());

    try {
      const res = await fetch("http://localhost:8080/channels", {
        method: "POST",
        credentials: "include",
        body: formData,
      });

      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || "Erreur lors de la création du canal");
      }

      await fetchChannels();
    } catch (err) {
      console.error("Erreur lors de la création du canal:", err);
      alert("Erreur lors de la création du canal.");
    }
  };

  // Open the invite modal
  const handleInviteClick = () => {
    setIsInviteModalOpen(true);
  };

  // Close the invite modal
  const handleCloseModal = () => {
    setIsInviteModalOpen(false);
  };

  return (
    <div className="p-4 text-white bg-gray-800 w-64 flex flex-col">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl font-bold">Canaux</h2>

        {currentUserId === serverOwnerId && (
          <button
            onClick={handleCreateChannel}
            title="Créer un nouveau canal"
            className="text-green-400 hover:text-green-300 text-2xl"
          >
            ➕
          </button>
        )}
      </div>

      <button
        onClick={handleInviteClick}
        className="text-blue-400 hover:text-blue-300 mb-4"
      >
        Inviter un ami au serveur
      </button>

      {loading && <div>Chargement...</div>}
      {error && <div className="text-red-500">{error}</div>}

      <ul>
        {channels.map((channel) => (
          <li
            key={channel.id}
            onClick={() => onChannelSelect(channel.id)}
            className="cursor-pointer hover:bg-gray-700 px-3 py-2 rounded transition"
          >
            #{channel.name}
          </li>
        ))}
      </ul>

      {/* Modal for inviting friend */}
      {isInviteModalOpen && (
        <AddFriendToServer serverId={serverId} onClose={handleCloseModal} />
      )}
    </div>
  );
}
