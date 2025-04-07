import { getSummoner } from "@/services/getSummoner";
import {
  useMutation,
  useQueryClient,
  useSuspenseQuery,
} from "@tanstack/react-query";
import { Badge } from "../ui/badge";
import { updateSummoner } from "@/services/updateSummoner";
import { Button } from "../ui/button";
import { IconLoader2 } from "@tabler/icons-react";
import TimeSince from "../timeSince";

const PlayerHeader: React.FC<{ name: string; tag: string }> = ({
  name,
  tag,
}) => {
  const queryClient = useQueryClient();

  const query = useSuspenseQuery(getSummoner(name, tag));

  const updateMutation = useMutation({
    mutationKey: ["UPDATE_SUMMONER", query.data.puuid],
    mutationFn: updateSummoner,
    onSuccess: () => {
      query.refetch();
      queryClient.invalidateQueries({
        predicate: ({ queryKey }) =>
          queryKey.includes(query.data.puuid) ||
          queryKey.includes(query.data.summonerId),
      });
    },
  });

  return (
    <div className="grid w-full grid-cols-[auto_1fr] items-center justify-items-start gap-2">
      <img
        className="row-span-3 h-full rounded border"
        width={124}
        src={`https://raw.communitydragon.org/latest/plugins/rcp-be-lol-game-data/global/default/v1/profile-icons/${query.data.profileIconId}.jpg`}
      />
      <div className="flex items-center gap-1">
        <h1 className="text-3xl font-bold">
          <span>{query.data.name}</span>
          <span className="text-muted-foreground">#{query.data.tag}</span>
        </h1>
        <Badge variant="secondary">{query.data.region}</Badge>
      </div>

      <Button
        onClick={() => updateMutation.mutate(query.data.puuid)}
        disabled={updateMutation.isPending}
      >
        {
          {
            idle: "Update",
            pending: (
              <>
                <IconLoader2 className="mr-1 animate-spin ease-in-out" />
                Updating...
              </>
            ),
            success: "Updated",
            error: "Error",
          }[updateMutation.status]
        }
      </Button>
      <span className="text-muted-foreground">
        {"Updated "}
        {query.data.statsUpdateTimestamp ? (
          <TimeSince date={new Date(query.data.statsUpdateTimestamp)} />
        ) : (
          "never"
        )}
      </span>
    </div>
  );
};
export default PlayerHeader;
