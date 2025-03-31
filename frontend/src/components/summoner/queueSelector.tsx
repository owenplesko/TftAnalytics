import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../ui/select";

const queueIdValues = {
  null: null,
  ranked: 1100,
  normal: 1090,
};

const QueueSelector: React.FC<{
  setQueueId: (queueId: number | null) => void;
}> = ({ setQueueId }) => {
  return (
    <Select
      defaultValue="null"
      onValueChange={(value) => setQueueId(queueIdValues[value])}
    >
      <SelectTrigger className="w-[200px]">
        <SelectValue />
      </SelectTrigger>
      <SelectContent>
        <SelectItem value="null">All Queues</SelectItem>
        <SelectItem value="ranked">Ranked</SelectItem>
        <SelectItem value="normal">Normal</SelectItem>
      </SelectContent>
    </Select>
  );
};
export default QueueSelector;
