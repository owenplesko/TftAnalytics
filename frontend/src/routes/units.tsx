import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import UnitIcon from "@/components/unitIcon";
import { GET_UNIT_STATS_KEY, getUnitStats } from "@/services/getUnitStats";
import { refreshUnitStats } from "@/services/refreshUnitStats";
import { UnitStat } from "@/services/types";
import { useMutation, useQueryClient, useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import {
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  getSortedRowModel,
  useReactTable,
} from "@tanstack/react-table";
import { ArrowUpDown } from "lucide-react";
import { unitData } from "../data/unitData";

export const Route = createFileRoute("/units")({
  component: UnitStatsPage,
  loader: async ({ context: { queryClient } }) => {
    await queryClient.ensureQueryData(getUnitStats());
  },
});

function UnitStatsPage() {
  return <UnitStatsTable />;
}

const columnHelper = createColumnHelper<UnitStat>();

const columns = [
  columnHelper.accessor("unit", {
    header: "Unit",
    cell: (props) => {
      const apiName = props.getValue().toLowerCase();
      return (
        <>
          <UnitIcon apiName={apiName} />
          {unitData[apiName]?.name ?? apiName}
        </>
      );
    },
    enableSorting: false,
  }),
  columnHelper.accessor("avgPlacement", {
    header: "Avg Place",
    cell: (props) => props.getValue().toPrecision(3),
  }),
  columnHelper.accessor("frequency", {
    header: "Frequency",
    cell: (props) => props.getValue(),
  }),
];

function UnitStatsTable() {
  const queryClient = useQueryClient()

  const query = useSuspenseQuery(getUnitStats());

  const refreshStatsMutation = useMutation({
    mutationKey: ["REFRESH_UNIT_STATS"], mutationFn: refreshUnitStats,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [GET_UNIT_STATS_KEY] })
    }
  })

  const table = useReactTable({
    columns,
    data: query.data,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
  });

  return (
    <div className="w-[1000px] space-y-2">
      <h1 className="col-span-2 text-4xl font-semibold">TFT Unit Stats</h1>
      <div className="flex items-end">
        <div className="grid w-64 grid-cols-[auto_auto] rounded-md bg-card p-4 text-xs">
          <span>Last Updated: </span>
          <span className="text-right">a few seconds ago</span>
          <span>Comps Analyzed: </span>
          <span className="text-right">69,420</span>
        </div>
        <Button className="ml-auto" onClick={() => refreshStatsMutation.mutate()}>Refresh</Button>
      </div>
      <div className="w-full rounded-md border bg-card">
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow>
                {headerGroup.headers.map((header) => (
                  <TableHead className="text-center first:text-left last:text-right">
                    {header.column.getCanSort() ? (
                      <Button
                        variant="ghost"
                        onClick={header.column.getToggleSortingHandler()}
                      >
                        {flexRender(
                          header.column.columnDef.header,
                          header.getContext(),
                        )}
                        <ArrowUpDown size={16} />
                      </Button>
                    ) : (
                      <div>
                        {flexRender(
                          header.column.columnDef.header,
                          header.getContext(),
                        )}
                      </div>
                    )}
                  </TableHead>
                ))}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows.map((row) => (
              <TableRow key={row.id}>
                {row.getVisibleCells().map((cell) => (
                  <TableCell
                    key={cell.id}
                    className="text-center first:text-left last:text-right"
                  >
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </TableCell>
                ))}
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}
