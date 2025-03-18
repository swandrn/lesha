import React from "react";

interface Channel {
  id: number;
  name: string;
}

const sampleChannels: Record<number, Channel[]> = {
  1: [{ id: 101, name: "general" }, { id: 102, name: "random" }],
  2: [{ id: 201, name: "game-talk" }, { id: 202, name: "game-news" }],
  3: [{ id: 301, name: "coding" }, { id: 302, name: "debugging" }],
  4: [{ id: 401, name: "music-share" }, { id: 402, name: "playlists" }],
  5: [{ id: 501, name: "movie-reviews" }, { id: 502, name: "watch-party" }],
};

interface ChannelListProps {
  serverId: number;
  onChannelSelect: (channelId: number) => void;
}

export function ChannelList(
  { serverId, onChannelSelect }: ChannelListProps,
): React.JSX.Element {
  const channels = sampleChannels[serverId] || [];

  return (
    <div className="w-40 md:w-56 h-screen bg-gray-800 text-white flex flex-col p-4 shadow-md">
      <h2 className="text-lg font-semibold mb-4">Channels</h2>
      {channels.map((channel) => (
        <div
          key={channel.id}
          className="p-2 text-sm rounded-md hover:bg-gray-700 cursor-pointer"
          onClick={() => onChannelSelect(channel.id)}
        >
          # {channel.name}
        </div>
      ))}
    </div>
  );
}
