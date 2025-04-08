import { QUEUES } from "@/lib/constants";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../ui/select";

const QueueSelector: React.FC<{
  setQueueId: (queueId: number | null) => void;
}> = ({ setQueueId }) => {
  // I REALLY HATE THIS!!!!!
  const handleValueChange = (value: string) => {
    return setQueueId(
      (
        QUEUES as Record<
          string,
          (typeof QUEUES)[keyof typeof QUEUES] | undefined
        >
      )[value]?.id ?? null,
    );
  };

  return (
    <Select defaultValue="all" onValueChange={handleValueChange}>
      <SelectTrigger className="w-[200px]">
        <SelectValue />
      </SelectTrigger>
      <SelectContent>
        <SelectItem value="all">All Queues</SelectItem>
        <SelectItem value="ranked">Ranked</SelectItem>
        <SelectItem value="normal">Normal</SelectItem>
      </SelectContent>
    </Select>
  );
};
export default QueueSelector;
