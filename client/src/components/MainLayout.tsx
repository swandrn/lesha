import { ServerList } from "./ServerList";
import { ChannelList } from "./ChannelList";
import { Chat } from "./Chat";
import EditAccount from "./EditAccount";
import { FriendList } from "./FriendList";
import CreateServer from "./CreateServer";

export interface Server {
    id: number;
    name: string;
    description: string;
    image: string;
    userId: number;
    createdAt: string;
    updatedAt: string;
  }

interface Props {
  servers: Server[];
  selectedServer: number | null;
  selectedChannel: number | null;
  isCreatingServer: boolean;
  isChannelListVisible: boolean;
  onServerSelect: (id: number) => void;
  onCreateNewServer: () => void;
  onChannelSelect: (id: number) => void;
  onServerCreated: () => void;
}

export const MainLayout = ({
  servers,
  selectedServer,
  selectedChannel,
  isCreatingServer,
  isChannelListVisible,
  onServerSelect,
  onCreateNewServer,
  onChannelSelect,
  onServerCreated,
}: Props) => {
  const selectedServerName = servers.find((s) => s.id === selectedServer)?.name;

  const specialViews: Record<string, React.JSX.Element> = {
    "Edit Account": <EditAccount />,
    "Friends": <FriendList />,
  };

  return (
    <div className="flex h-screen w-screen bg-gray-900">
      <ServerList
        servers={servers}
        onServerSelect={onServerSelect}
        onCreateNewServer={onCreateNewServer}
      />

      {selectedServerName && specialViews[selectedServerName]}

      {isCreatingServer && <CreateServer onServerCreated={onServerCreated} />}

      {selectedServer !== null &&
        !isCreatingServer &&
        !specialViews[selectedServerName!] && (
          <>
            {isChannelListVisible && (
              <ChannelList
                serverId={selectedServer}
                onChannelSelect={onChannelSelect}
              />
            )}
            {selectedChannel && (
              <Chat
                channelId={selectedChannel}
                onToggleChannels={() => onChannelSelect(1)}
              />
            )}
          </>
        )}
    </div>
  );
};