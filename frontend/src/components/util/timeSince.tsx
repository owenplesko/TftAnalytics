import { formatTimeSince } from "@/lib/utils";
import { useEffect, useState } from "react";

const TimeSince: React.FC<{ date: Date }> = ({ date }) => {
  const [timeSince, setTimeSince] = useState<string>();

  useEffect(() => {
    setTimeSince(formatTimeSince(date));

    const interval = setInterval(
      () => setTimeSince(formatTimeSince(date)),
      1000,
    );

    return () => clearInterval(interval);
  }, [date]);

  return timeSince;
};
export default TimeSince;
