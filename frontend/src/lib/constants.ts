export const QUEUES = {
    normal: {
        label: "Normal",
        id: 1090
    },
    ranked: {
        label: "Ranked",
        id: 1100
    },
    tutorial: {
        label: "Tutorial",
        id: 1110
    },
    hyperRoll: {
        label: "Hyper Roll",
        id: 1130
    },
    doubleUp: {
        label: "Double Up",
        id: 1160
    },
    tockersTrials: {
        label: "Tocker's Trials",
        id: 1220
    },
    revival: {
        label: "Revival",
        id: 6100
    }
} as const;

// maps queue ids to labels
export const QUEUE_LABELS = Object.fromEntries(
    Object.values(QUEUES).map(({ id, label }) => [id, label])
) as Record<number, string | undefined>;
