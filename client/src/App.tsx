import { useState } from "react";
import { ServerList } from "./components/ServerList";
import { ChannelList } from "./components/ChannelList";
import { Chat } from "./components/Chat";

function App() {
  const [selectedServer, setSelectedServer] = useState<number | null>(null);
  const [selectedChannel, setSelectedChannel] = useState<number | null>(null);
  const [isChannelListVisible, setIsChannelListVisible] = useState(true);

  const handleServerSelect = (serverId: number) => {
    setSelectedServer(serverId);
    setSelectedChannel(null); // Reset selected channel
    setIsChannelListVisible(true); // Ensure channel list is visible
  };

  const handleChannelSelect = (channelId: number) => {
    setSelectedChannel(channelId);
    setIsChannelListVisible(false); // Hide channel list when a channel is selected
  };

  return (
    <div className="flex h-screen w-screen bg-gray-900">
      <ServerList onServerSelect={handleServerSelect} />

      {/* Show ChannelList only if a server is selected AND it's visible */}
      {selectedServer && isChannelListVisible && (
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

export default App;
